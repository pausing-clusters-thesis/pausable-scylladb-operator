package ingress

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/routing"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/scheme"
	"github.com/scylladb/scylla-operator/pkg/helpers/slices"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1informers "k8s.io/client-go/informers/core/v1"
	networkingv1informers "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	networkingv1listers "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	ControllerName  = "IngressController"
	maxSyncDuration = 30 * time.Second
)

var (
	controllerKey        = "key"
	keyFunc              = cache.DeletionHandlingMetaNamespaceKeyFunc
	ingressControllerGVK = networkingv1.SchemeGroupVersion.WithKind("Ingress")
)

type Controller struct {
	kubeClient kubernetes.Interface

	ingressLister      networkingv1listers.IngressLister
	ingressClassLister networkingv1listers.IngressClassLister
	serviceLister      corev1listers.ServiceLister
	podLister          corev1listers.PodLister

	cachesToSync []cache.InformerSynced

	eventRecorder record.EventRecorder

	queue workqueue.RateLimitingInterface

	routingTable *routing.Table
}

func NewController(
	kubeClient kubernetes.Interface,
	ingressInformer networkingv1informers.IngressInformer,
	ingressClassInformer networkingv1informers.IngressClassInformer,
	serviceInformer corev1informers.ServiceInformer,
	podInformer corev1informers.PodInformer,
	routingTable *routing.Table,
) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&corev1client.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	ic := &Controller{
		kubeClient: kubeClient,

		ingressLister:      ingressInformer.Lister(),
		ingressClassLister: ingressClassInformer.Lister(),
		serviceLister:      serviceInformer.Lister(),
		podLister:          podInformer.Lister(),

		cachesToSync: []cache.InformerSynced{
			ingressInformer.Informer().HasSynced,
			ingressClassInformer.Informer().HasSynced,
			serviceInformer.Informer().HasSynced,
			podInformer.Informer().HasSynced,
		},

		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "ingress-controller"}),

		queue: workqueue.NewRateLimitingQueueWithConfig(
			workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{
				Name: "ingress",
			},
		),

		routingTable: routingTable,
	}

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ic.addIngress,
		UpdateFunc: ic.updateIngress,
		DeleteFunc: ic.deleteIngress,
	})

	ingressClassInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ic.addIngressClass,
		UpdateFunc: ic.updateIngressClass,
		DeleteFunc: ic.deleteIngressClass,
	})

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ic.addService,
		UpdateFunc: ic.updateService,
		DeleteFunc: ic.deleteService,
	})

	// Enqueue preemptively to initialize the routing table.
	ic.enqueue()

	return ic, nil
}

func (ic *Controller) processNextItem(ctx context.Context) bool {
	key, quit := ic.queue.Get()
	if quit {
		return false
	}
	defer ic.queue.Done(key)

	ctx, cancel := context.WithTimeoutCause(ctx, maxSyncDuration, fmt.Errorf("exceeded max sync duration (%v)", maxSyncDuration))
	defer cancel()
	err := ic.sync(ctx)
	// TODO: Do smarter filtering then just Reduce to handle cases like 2 conflict errors.
	err = utilerrors.Reduce(err)
	switch {
	case err == nil:
		ic.queue.Forget(key)
		return true

	case apierrors.IsConflict(err):
		klog.V(2).InfoS("Hit conflict, will retry in a bit", "Key", key, "Error", err)

	case apierrors.IsAlreadyExists(err):
		klog.V(2).InfoS("Hit already exists, will retry in a bit", "Key", key, "Error", err)

	default:
		utilruntime.HandleError(fmt.Errorf("syncing key '%v' failed: %v", key, err))
	}

	ic.queue.AddRateLimited(key)

	return true
}

func (ic *Controller) runWorker(ctx context.Context) {
	for ic.processNextItem(ctx) {
	}
}

func (ic *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()

	klog.InfoS("Starting controller", "controller", ControllerName)

	var wg sync.WaitGroup
	defer func() {
		klog.InfoS("Shutting down controller", "controller", ControllerName)
		ic.queue.ShutDown()
		wg.Wait()
		klog.InfoS("Shut down controller", "controller", ControllerName)
	}()

	if !cache.WaitForNamedCacheSync(ControllerName, ctx.Done(), ic.cachesToSync...) {
		return
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wait.UntilWithContext(ctx, ic.runWorker, time.Second)
		}()
	}

	<-ctx.Done()
}

func (ic *Controller) enqueue() {
	ic.queue.Add(controllerKey)
}

func (ic *Controller) addIngress(obj interface{}) {
	ingress := obj.(*networkingv1.Ingress)

	provisionedByUs, err := ic.isIngressProvisionedByUs(ingress)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueueing Ingress referencing IngressClass not provisioned by us", "Ingress", klog.KObj(ingress), "RV", ingress.ResourceVersion)
	}

	klog.V(4).InfoS("Observed addition of Ingress", "Ingress", klog.KObj(ingress))
	ic.enqueue()
}

func (ic *Controller) updateIngress(old, cur interface{}) {
	oldIngress := old.(*networkingv1.Ingress)
	curIngress := cur.(*networkingv1.Ingress)

	if curIngress.UID != oldIngress.UID {
		key, err := keyFunc(oldIngress)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", oldIngress, err))
			return
		}
		ic.deleteIngressClass(cache.DeletedFinalStateUnknown{
			Key: key,
			Obj: oldIngress,
		})
	}

	provisionedByUs, err := ic.isIngressProvisionedByUs(curIngress)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueuing Ingress referencing IngressClass not provisioned by us", "Ingress", klog.KObj(curIngress), "RV", curIngress.ResourceVersion)
		return
	}

	klog.V(4).InfoS(
		"Observed update of Ingress",
		"Ingress", klog.KObj(curIngress),
		"RV", fmt.Sprintf("%s-%s", oldIngress.ResourceVersion, curIngress.ResourceVersion),
		"UID", fmt.Sprintf("%s-%s", oldIngress.UID, curIngress.UID),
	)
	ic.enqueue()
}

func (ic *Controller) deleteIngress(obj interface{}) {
	var ingress *networkingv1.Ingress
	var ok bool

	ingress, ok = obj.(*networkingv1.Ingress)
	if !ok {
		var tombstone cache.DeletedFinalStateUnknown
		tombstone, ok = obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}

		ingress, ok = tombstone.Obj.(*networkingv1.Ingress)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an Ingress %#v", obj))
			return
		}
	}

	provisionedByUs, err := ic.isIngressProvisionedByUs(ingress)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueueing Ingress referencing IngressClass not provisioned by us", "Ingress", klog.KObj(ingress), "RV", ingress.ResourceVersion)
	}

	klog.V(4).InfoS("Observed deletion of Ingress", "Ingress", klog.KObj(ingress), "RV", ingress.ResourceVersion)
	ic.enqueue()
}

func (ic *Controller) addIngressClass(obj interface{}) {
	ingressClass := obj.(*networkingv1.IngressClass)

	provisionedByUs := isIngressClassProvisionedByUs(ingressClass)
	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueuing IngressClass not provisioned by us", "IngressClass", klog.KObj(ingressClass), "RV", ingressClass.ResourceVersion)
		return
	}

	klog.V(4).InfoS("Observed addition of IngressClass", "IngressClass", klog.KObj(ingressClass))
	ic.enqueue()
}

func (ic *Controller) updateIngressClass(old, cur interface{}) {
	oldIngressClass := old.(*networkingv1.IngressClass)
	curIngressClass := cur.(*networkingv1.IngressClass)

	if curIngressClass.UID != oldIngressClass.UID {
		key, err := keyFunc(oldIngressClass)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", oldIngressClass, err))
			return
		}
		ic.deleteIngressClass(cache.DeletedFinalStateUnknown{
			Key: key,
			Obj: oldIngressClass,
		})
	}

	provisionedByUs := isIngressClassProvisionedByUs(curIngressClass)
	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueuing IngressClass not provisioned by us", "IngressClass", klog.KObj(curIngressClass), "RV", curIngressClass.ResourceVersion)
		return
	}

	klog.V(4).InfoS(
		"Observed update of IngressClass",
		"IngressClass", klog.KObj(curIngressClass),
		"RV", fmt.Sprintf("%s-%s", oldIngressClass.ResourceVersion, curIngressClass.ResourceVersion),
		"UID", fmt.Sprintf("%s-%s", oldIngressClass.UID, curIngressClass.UID),
	)
	ic.enqueue()
}

func (ic *Controller) deleteIngressClass(obj interface{}) {
	var ingressClass *networkingv1.IngressClass
	var ok bool

	ingressClass, ok = obj.(*networkingv1.IngressClass)
	if !ok {
		var tombstone cache.DeletedFinalStateUnknown
		tombstone, ok = obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}

		ingressClass, ok = tombstone.Obj.(*networkingv1.IngressClass)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not an IngressClass %#v", obj))
			return
		}
	}

	provisionedByUs := isIngressClassProvisionedByUs(ingressClass)
	if !provisionedByUs {
		klog.V(5).InfoS("Not enqueuing IngressClass not provisioned by us", "IngressClass", klog.KObj(ingressClass), "RV", ingressClass.ResourceVersion)
		return
	}

	klog.V(4).InfoS("Observed deletion of IngressClass", "IngressClass", klog.KObj(ingressClass), "RV", ingressClass.ResourceVersion)
	ic.enqueue()
}

func (ic *Controller) addService(obj interface{}) {
	service := obj.(*corev1.Service)

	referencedByIngressProvisionedByUs, err := ic.isServiceReferencedByIngressProvisionedByUs(service)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !referencedByIngressProvisionedByUs {
		klog.V(5).InfoS("Not enqueuing Service not referenced by any Ingress provisioned by us", "Service", klog.KObj(service), "RV", service.ResourceVersion)
		return
	}

	klog.V(4).InfoS("Observed addition of Service", "Service", klog.KObj(service))
	ic.enqueue()
}

func (ic *Controller) updateService(old, cur interface{}) {
	oldService := old.(*corev1.Service)
	curService := cur.(*corev1.Service)

	if curService.UID != oldService.UID {
		key, err := keyFunc(oldService)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", oldService, err))
			return
		}
		ic.deleteIngressClass(cache.DeletedFinalStateUnknown{
			Key: key,
			Obj: oldService,
		})
	}

	referencedByIngressProvisionedByUs, err := ic.isServiceReferencedByIngressProvisionedByUs(curService)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !referencedByIngressProvisionedByUs {
		klog.V(5).InfoS("Not enqueuing Service not referenced by any Ingress provisioned by us", "Service", klog.KObj(curService), "RV", curService.ResourceVersion)
		return
	}

	klog.V(4).InfoS(
		"Observed update of Service",
		"Service", klog.KObj(curService),
		"RV", fmt.Sprintf("%s-%s", oldService.ResourceVersion, curService.ResourceVersion),
		"UID", fmt.Sprintf("%s-%s", oldService.UID, curService.UID),
	)
	ic.enqueue()
}

func (ic *Controller) deleteService(obj interface{}) {
	var service *corev1.Service
	var ok bool

	service, ok = obj.(*corev1.Service)
	if !ok {
		var tombstone cache.DeletedFinalStateUnknown
		tombstone, ok = obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}

		service, ok = tombstone.Obj.(*corev1.Service)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Service %#v", obj))
			return
		}
	}

	referencedByIngressProvisionedByUs, err := ic.isServiceReferencedByIngressProvisionedByUs(service)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	if !referencedByIngressProvisionedByUs {
		klog.V(5).InfoS("Not enqueuing Service not referenced by any Ingress provisioned by us", "Service", klog.KObj(service), "RV", service.ResourceVersion)
		return
	}

	klog.V(4).InfoS("Observed deletion of Service", "Service", klog.KObj(service), "RV", service.ResourceVersion)
	ic.enqueue()
}

func (ic *Controller) isServiceReferencedByIngressProvisionedByUs(svc *corev1.Service) (bool, error) {
	ingresses, err := ic.ingressLister.Ingresses(svc.Namespace).List(labels.Everything())
	if err != nil {
		return false, fmt.Errorf("can't list Ingresses: %w", err)
	}

	var ingressesProvisionedByUs []*networkingv1.Ingress
	var errs []error
	for _, ingress := range ingresses {
		var provisionedByUs bool
		provisionedByUs, err = ic.isIngressProvisionedByUs(ingress)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't verify if Ingress %q is provisioned by us: %w", naming.ObjRef(ingress), err))
			continue
		}

		if provisionedByUs {
			ingressesProvisionedByUs = append(ingressesProvisionedByUs, ingress)
		}
	}

	isIngressBackendReferencingService := func(ingressBackend *networkingv1.IngressBackend) bool {
		return ingressBackend.Service != nil && ingressBackend.Service.Name == svc.Name
	}

	isHTTPIngressPathReferencingService := func(httpIngressPath networkingv1.HTTPIngressPath) bool {
		return isIngressBackendReferencingService(&httpIngressPath.Backend)
	}

	isIngressRuleReferencingService := func(ingressRule networkingv1.IngressRule) bool {
		return ingressRule.HTTP != nil && slices.Contains(ingressRule.HTTP.Paths, isHTTPIngressPathReferencingService)
	}

	isIngressReferencingService := func(ingress *networkingv1.Ingress) bool {
		if ingress.Spec.DefaultBackend != nil && isIngressBackendReferencingService(ingress.Spec.DefaultBackend) {
			return true
		}

		return slices.Contains(ingress.Spec.Rules, isIngressRuleReferencingService)
	}

	return slices.Contains(ingressesProvisionedByUs, isIngressReferencingService), nil
}

func (ic *Controller) isIngressProvisionedByUs(ingress *networkingv1.Ingress) (bool, error) {
	if ingress.Spec.IngressClassName == nil {
		return false, nil
	}

	ingressClass, err := ic.ingressClassLister.Get(*ingress.Spec.IngressClassName)
	if err != nil {
		return false, fmt.Errorf("can't get IngressClass %q: %w", *ingress.Spec.IngressClassName, err)
	}

	return isIngressClassProvisionedByUs(ingressClass), nil
}

func isIngressClassProvisionedByUs(ingressClass *networkingv1.IngressClass) bool {
	return ingressClass.Spec.Controller == naming.IngressControllerName
}
