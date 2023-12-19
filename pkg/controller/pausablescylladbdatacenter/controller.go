package pausablescylladbdatacenter

import (
	"context"
	"fmt"
	"sync"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	pausingv1alpha1client "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned/typed/pausing/v1alpha1"
	pausingv1alpha1informers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions/pausing/v1alpha1"
	pausingv1alpha1listers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/listers/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/scheme"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	scyllav1alpha1informers "github.com/scylladb/scylla-operator/pkg/client/scylla/informers/externalversions/scylla/v1alpha1"
	scyllav1alpha1listers "github.com/scylladb/scylla-operator/pkg/client/scylla/listers/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	socrypto "github.com/scylladb/scylla-operator/pkg/crypto"
	"github.com/scylladb/scylla-operator/pkg/kubeinterfaces"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	ControllerName = "PausableScyllaDBDatacenterController"
)

var (
	keyFunc                                 = cache.DeletionHandlingMetaNamespaceKeyFunc
	pausableScyllaDBDatacenterControllerGVK = pausingv1alpha1.GroupVersion.WithKind("PausableScyllaDBDatacenter")
)

type Controller struct {
	kubeClient    kubernetes.Interface
	pausingClient pausingv1alpha1client.PausingV1alpha1Interface

	pausableScyllaDBDatacenterLister pausingv1alpha1listers.PausableScyllaDBDatacenterLister
	scyllaDBDatacenterClaimLister    pausingv1alpha1listers.ScyllaDBDatacenterClaimLister
	scyllaDBDatacenterPoolLister     pausingv1alpha1listers.ScyllaDBDatacenterPoolLister

	scyllaDBDatacenterLister scyllav1alpha1listers.ScyllaDBDatacenterLister

	persistentVolumeClaimLister corev1listers.PersistentVolumeClaimLister
	secretLister                corev1listers.SecretLister
	configMapLister             corev1listers.ConfigMapLister

	cachesToSync []cache.InformerSynced

	eventRecorder record.EventRecorder

	queue    workqueue.RateLimitingInterface
	handlers *socontrollerhelpers.Handlers[*pausingv1alpha1.PausableScyllaDBDatacenter]

	keyGetter socrypto.RSAKeyGetter
}

func NewController(
	kubeClient kubernetes.Interface,
	pausingClient pausingv1alpha1client.PausingV1alpha1Interface,
	pausableScyllaDBDatacenterInformer pausingv1alpha1informers.PausableScyllaDBDatacenterInformer,
	scyllaDBDatacenterClaimInformer pausingv1alpha1informers.ScyllaDBDatacenterClaimInformer,
	scyllaDBDatacenterPoolInformer pausingv1alpha1informers.ScyllaDBDatacenterPoolInformer,
	persistentVolumeClaimInformer corev1informers.PersistentVolumeClaimInformer,
	scyllaDBDatacenterInformer scyllav1alpha1informers.ScyllaDBDatacenterInformer,
	secretInformer corev1informers.SecretInformer,
	configMapInformer corev1informers.ConfigMapInformer,
	keyGetter socrypto.RSAKeyGetter,
) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&corev1client.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	pscc := &Controller{
		kubeClient:    kubeClient,
		pausingClient: pausingClient,

		pausableScyllaDBDatacenterLister: pausableScyllaDBDatacenterInformer.Lister(),
		scyllaDBDatacenterClaimLister:    scyllaDBDatacenterClaimInformer.Lister(),
		scyllaDBDatacenterPoolLister:     scyllaDBDatacenterPoolInformer.Lister(),
		persistentVolumeClaimLister:      persistentVolumeClaimInformer.Lister(),
		scyllaDBDatacenterLister:         scyllaDBDatacenterInformer.Lister(),
		secretLister:                     secretInformer.Lister(),
		configMapLister:                  configMapInformer.Lister(),

		cachesToSync: []cache.InformerSynced{
			pausableScyllaDBDatacenterInformer.Informer().HasSynced,
			scyllaDBDatacenterClaimInformer.Informer().HasSynced,
			persistentVolumeClaimInformer.Informer().HasSynced,
			scyllaDBDatacenterInformer.Informer().HasSynced,
			secretInformer.Informer().HasSynced,
			configMapInformer.Informer().HasSynced,
		},

		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "pausablescylladbdatacenter-controller"}),

		queue: workqueue.NewRateLimitingQueueWithConfig(
			workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{
				Name: "pausablescylladbdatacenter",
			},
		),

		keyGetter: keyGetter,
	}

	var err error
	pscc.handlers, err = socontrollerhelpers.NewHandlers[*pausingv1alpha1.PausableScyllaDBDatacenter](
		pscc.queue,
		keyFunc,
		scheme.Scheme,
		pausableScyllaDBDatacenterControllerGVK,
		kubeinterfaces.NamespacedGetList[*pausingv1alpha1.PausableScyllaDBDatacenter]{
			GetFunc: func(namespace, name string) (*pausingv1alpha1.PausableScyllaDBDatacenter, error) {
				return pscc.pausableScyllaDBDatacenterLister.PausableScyllaDBDatacenters(namespace).Get(name)
			},
			ListFunc: func(namespace string, selector labels.Selector) (ret []*pausingv1alpha1.PausableScyllaDBDatacenter, err error) {
				return pscc.pausableScyllaDBDatacenterLister.PausableScyllaDBDatacenters(namespace).List(selector)
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't create handlers: %w", err)
	}

	pausableScyllaDBDatacenterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pscc.addPausableScyllaDBDatacenter,
		UpdateFunc: pscc.updatePausableScyllaDBDatacenter,
		DeleteFunc: pscc.deletePausableScyllaDBDatacenter,
	})

	scyllaDBDatacenterClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pscc.addScyllaDBDatacenterClaim,
		UpdateFunc: pscc.updateScyllaDBDatacenterClaim,
		DeleteFunc: pscc.deleteScyllaDBDatacenterClaim,
	})

	persistentVolumeClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pscc.addPersistentVolumeClaim,
		UpdateFunc: pscc.updatePersistentVolumeClaim,
		DeleteFunc: pscc.deletePersistentVolumeClaim,
	})

	scyllaDBDatacenterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pscc.addScyllaDBDatacenter,
		UpdateFunc: pscc.updateScyllaDBDatacenter,
		DeleteFunc: pscc.deleteScyllaDBDatacenter,
	})

	secretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pscc.addSecret,
		UpdateFunc: pscc.updateSecret,
		DeleteFunc: pscc.deleteSecret,
	})

	return pscc, nil
}

func (psdcc *Controller) processNextItem(ctx context.Context) bool {
	key, quit := psdcc.queue.Get()
	if quit {
		return false
	}
	defer psdcc.queue.Done(key)

	err := psdcc.sync(ctx, key.(string))
	// TODO: Do smarter filtering then just Reduce to handle cases like 2 conflict errors.
	err = utilerrors.Reduce(err)
	switch {
	case err == nil:
		psdcc.queue.Forget(key)
		return true

	case apierrors.IsConflict(err):
		klog.V(2).InfoS("Hit conflict, will retry in a bit", "Key", key, "Error", err)

	case apierrors.IsAlreadyExists(err):
		klog.V(2).InfoS("Hit already exists, will retry in a bit", "Key", key, "Error", err)

	default:
		utilruntime.HandleError(fmt.Errorf("syncing key '%v' failed: %v", key, err))
	}

	psdcc.queue.AddRateLimited(key)

	return true
}

func (psdcc *Controller) runWorker(ctx context.Context) {
	for psdcc.processNextItem(ctx) {
	}
}

func (psdcc *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()

	klog.InfoS("Starting controller", "controller", ControllerName)

	var wg sync.WaitGroup
	defer func() {
		klog.InfoS("Shutting down controller", "controller", ControllerName)
		psdcc.queue.ShutDown()
		wg.Wait()
		klog.InfoS("Shut down controller", "controller", ControllerName)
	}()

	if !cache.WaitForNamedCacheSync(ControllerName, ctx.Done(), psdcc.cachesToSync...) {
		return
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wait.UntilWithContext(ctx, psdcc.runWorker, time.Second)
		}()
	}

	<-ctx.Done()
}

func (psdcc *Controller) addPausableScyllaDBDatacenter(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*pausingv1alpha1.PausableScyllaDBDatacenter),
		psdcc.handlers.Enqueue,
	)
}

func (psdcc *Controller) updatePausableScyllaDBDatacenter(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*pausingv1alpha1.PausableScyllaDBDatacenter),
		cur.(*pausingv1alpha1.PausableScyllaDBDatacenter),
		psdcc.handlers.Enqueue,
		psdcc.deletePausableScyllaDBDatacenter,
	)
}

func (psdcc *Controller) deletePausableScyllaDBDatacenter(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj,
		psdcc.handlers.Enqueue,
	)
}

func (psdcc *Controller) addScyllaDBDatacenterClaim(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) updateScyllaDBDatacenterClaim(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		cur.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		psdcc.handlers.EnqueueOwner,
		psdcc.deleteScyllaDBDatacenterClaim,
	)
}

func (psdcc *Controller) deleteScyllaDBDatacenterClaim(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) addPersistentVolumeClaim(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*corev1.PersistentVolumeClaim),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) updatePersistentVolumeClaim(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*corev1.PersistentVolumeClaim),
		cur.(*corev1.PersistentVolumeClaim),
		psdcc.handlers.EnqueueOwner,
		psdcc.deletePersistentVolumeClaim,
	)
}

func (psdcc *Controller) deletePersistentVolumeClaim(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj.(*corev1.PersistentVolumeClaim),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) addScyllaDBDatacenter(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*scyllav1alpha1.ScyllaDBDatacenter),
		psdcc.enqueueOwnerThroughScyllaDBDatacenterOwner,
	)
}

func (psdcc *Controller) updateScyllaDBDatacenter(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*scyllav1alpha1.ScyllaDBDatacenter),
		cur.(*scyllav1alpha1.ScyllaDBDatacenter),
		psdcc.enqueueOwnerThroughScyllaDBDatacenterOwner,
		psdcc.deleteScyllaDBDatacenter,
	)
}

func (psdcc *Controller) deleteScyllaDBDatacenter(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj.(*scyllav1alpha1.ScyllaDBDatacenter),
		psdcc.enqueueOwnerThroughScyllaDBDatacenterOwner,
	)
}

func (psdcc *Controller) addSecret(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*corev1.Secret),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) updateSecret(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*corev1.Secret),
		cur.(*corev1.Secret),
		psdcc.handlers.EnqueueOwner,
		psdcc.deleteSecret,
	)
}

func (psdcc *Controller) deleteSecret(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj.(*corev1.Secret),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) addConfigMap(obj interface{}) {
	psdcc.handlers.HandleAdd(
		obj.(*corev1.ConfigMap),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) updateConfigMap(old, cur interface{}) {
	psdcc.handlers.HandleUpdate(
		old.(*corev1.ConfigMap),
		cur.(*corev1.ConfigMap),
		psdcc.handlers.EnqueueOwner,
		psdcc.deleteConfigMap,
	)
}

func (psdcc *Controller) deleteConfigMap(obj interface{}) {
	psdcc.handlers.HandleDelete(
		obj.(*corev1.ConfigMap),
		psdcc.handlers.EnqueueOwner,
	)
}

func (psdcc *Controller) enqueueOwnerThroughScyllaDBDatacenterOwner(depth int, obj kubeinterfaces.ObjectInterface, op socontrollerhelpers.HandlerOperationType) {
	sdcc := psdcc.resolveScyllaDBDatacenterClaimController(obj)
	if sdcc == nil {
		return
	}

	psdcc.handlers.EnqueueOwner(depth+1, sdcc, op)
}

func (psdcc *Controller) resolveScyllaDBDatacenterClaimController(obj metav1.Object) *pausingv1alpha1.ScyllaDBDatacenterClaim {
	controllerRef := metav1.GetControllerOf(obj)
	if controllerRef == nil {
		return nil
	}

	if controllerRef.Kind != pausingv1alpha1.ScyllaDBDatacenterClaimControllerGVK.Kind {
		return nil
	}

	sdcc, err := psdcc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(obj.GetNamespace()).Get(controllerRef.Name)
	if err != nil {
		return nil
	}

	if sdcc.UID != controllerRef.UID {
		return nil
	}

	return sdcc
}
