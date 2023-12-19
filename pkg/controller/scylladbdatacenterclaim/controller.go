package scylladbdatacenterclaim

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
	scyllav1alpha1client "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned/typed/scylla/v1alpha1"
	scyllav1alpha1informers "github.com/scylladb/scylla-operator/pkg/client/scylla/informers/externalversions/scylla/v1alpha1"
	scyllav1alpha1listers "github.com/scylladb/scylla-operator/pkg/client/scylla/listers/scylla/v1alpha1"
	"github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/kubeinterfaces"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	ControllerName = "ScyllaDBDatacenterClaimController"
)

var (
	keyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type Controller struct {
	kubeClient             kubernetes.Interface
	pausableScyllaDBClient pausingv1alpha1client.PausingV1alpha1Interface
	scyllaClient           scyllav1alpha1client.ScyllaV1alpha1Interface

	scyllaDBDatacenterClaimLister pausingv1alpha1listers.ScyllaDBDatacenterClaimLister
	scyllaDBDatacenterLister      scyllav1alpha1listers.ScyllaDBDatacenterLister

	cachesToSync []cache.InformerSynced

	eventRecorder record.EventRecorder

	queue    workqueue.RateLimitingInterface
	handlers *controllerhelpers.Handlers[*pausingv1alpha1.ScyllaDBDatacenterClaim]
}

func NewController(
	kubeClient kubernetes.Interface,
	pausableScyllaDBClient pausingv1alpha1client.PausingV1alpha1Interface,
	scyllaClient scyllav1alpha1client.ScyllaV1alpha1Interface,
	scyllaDBDatacenterClaimInformer pausingv1alpha1informers.ScyllaDBDatacenterClaimInformer,
	scyllaDBDatacenterInformer scyllav1alpha1informers.ScyllaDBDatacenterInformer,
) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&corev1client.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	sdccc := &Controller{
		kubeClient:             kubeClient,
		pausableScyllaDBClient: pausableScyllaDBClient,
		scyllaClient:           scyllaClient,

		scyllaDBDatacenterClaimLister: scyllaDBDatacenterClaimInformer.Lister(),
		scyllaDBDatacenterLister:      scyllaDBDatacenterInformer.Lister(),

		cachesToSync: []cache.InformerSynced{
			scyllaDBDatacenterClaimInformer.Informer().HasSynced,
			scyllaDBDatacenterInformer.Informer().HasSynced,
		},

		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "scylladbdatacenterclaim-controller"}),

		queue: workqueue.NewRateLimitingQueueWithConfig(
			workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{
				Name: "scylladbdatacenterclaim",
			},
		),
	}

	var err error
	sdccc.handlers, err = controllerhelpers.NewHandlers[*pausingv1alpha1.ScyllaDBDatacenterClaim](
		sdccc.queue,
		keyFunc,
		scheme.Scheme,
		pausingv1alpha1.ScyllaDBDatacenterClaimControllerGVK,
		kubeinterfaces.NamespacedGetList[*pausingv1alpha1.ScyllaDBDatacenterClaim]{
			GetFunc: func(namespace, name string) (*pausingv1alpha1.ScyllaDBDatacenterClaim, error) {
				return sdccc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(namespace).Get(name)
			},
			ListFunc: func(namespace string, selector labels.Selector) (ret []*pausingv1alpha1.ScyllaDBDatacenterClaim, err error) {
				return sdccc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(namespace).List(selector)
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't create handlers: %w", err)
	}

	scyllaDBDatacenterClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    sdccc.addScyllaDBDatacenterClaim,
		UpdateFunc: sdccc.updateScyllaDBDatacenterClaim,
		DeleteFunc: sdccc.deleteScyllaDBDatacenterClaim,
	})

	scyllaDBDatacenterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    sdccc.addScyllaDBDatacenter,
		UpdateFunc: sdccc.updateScyllaDBDatacenter,
		DeleteFunc: sdccc.deleteScyllaDBDatacenter,
	})

	return sdccc, nil
}

func (sdccc *Controller) processNextItem(ctx context.Context) bool {
	key, quit := sdccc.queue.Get()
	if quit {
		return false
	}
	defer sdccc.queue.Done(key)

	err := sdccc.sync(ctx, key.(string))
	// TODO: Do smarter filtering then just Reduce to handle cases like 2 conflict errors.
	err = utilerrors.Reduce(err)
	switch {
	case err == nil:
		sdccc.queue.Forget(key)
		return true

	case apierrors.IsConflict(err):
		klog.V(2).InfoS("Hit conflict, will retry in a bit", "Key", key, "Error", err)

	case apierrors.IsAlreadyExists(err):
		klog.V(2).InfoS("Hit already exists, will retry in a bit", "Key", key, "Error", err)

	default:
		utilruntime.HandleError(fmt.Errorf("syncing key '%v' failed: %v", key, err))
	}

	sdccc.queue.AddRateLimited(key)

	return true
}

func (sdccc *Controller) runWorker(ctx context.Context) {
	for sdccc.processNextItem(ctx) {
	}
}

func (sdccc *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()

	klog.InfoS("Starting controller", "controller", ControllerName)

	var wg sync.WaitGroup
	defer func() {
		klog.InfoS("Shutting down controller", "controller", ControllerName)
		sdccc.queue.ShutDown()
		wg.Wait()
		klog.InfoS("Shut down controller", "controller", ControllerName)
	}()

	if !cache.WaitForNamedCacheSync(ControllerName, ctx.Done(), sdccc.cachesToSync...) {
		return
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wait.UntilWithContext(ctx, sdccc.runWorker, time.Second)
		}()
	}

	<-ctx.Done()
}

func (sdccc *Controller) addScyllaDBDatacenterClaim(obj interface{}) {
	sdccc.handlers.HandleAdd(
		obj.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		sdccc.handlers.Enqueue,
	)
}

func (sdccc *Controller) updateScyllaDBDatacenterClaim(old, cur interface{}) {
	sdccc.handlers.HandleUpdate(
		old.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		cur.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		sdccc.handlers.Enqueue,
		sdccc.deleteScyllaDBDatacenterClaim,
	)
}

func (sdccc *Controller) deleteScyllaDBDatacenterClaim(obj interface{}) {
	sdccc.handlers.HandleDelete(
		obj,
		sdccc.handlers.Enqueue,
	)
}

func (sdccc *Controller) addScyllaDBDatacenter(obj interface{}) {
	sdccc.handlers.HandleAdd(
		obj.(*scyllav1alpha1.ScyllaDBDatacenter),
		sdccc.handlers.EnqueueOwner,
	)
}

func (sdccc *Controller) updateScyllaDBDatacenter(old, cur interface{}) {
	sdccc.handlers.HandleUpdate(
		old.(*scyllav1alpha1.ScyllaDBDatacenter),
		cur.(*scyllav1alpha1.ScyllaDBDatacenter),
		sdccc.handlers.EnqueueOwner,
		sdccc.deleteScyllaDBDatacenter,
	)
}

func (sdccc *Controller) deleteScyllaDBDatacenter(obj interface{}) {
	sdccc.handlers.HandleDelete(
		obj,
		sdccc.handlers.EnqueueOwner,
	)
}
