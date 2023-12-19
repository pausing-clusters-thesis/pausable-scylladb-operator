package scylladbdatacenterpool

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
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
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
	ControllerName = "ScyllaDBDatacenterPoolController"
)

var (
	keyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type Controller struct {
	kubeClient    kubernetes.Interface
	pausingClient pausingv1alpha1client.PausingV1alpha1Interface
	scyllaClient  scyllav1alpha1client.ScyllaV1alpha1Interface

	scyllaDBDatacenterPoolLister  pausingv1alpha1listers.ScyllaDBDatacenterPoolLister
	scyllaDBDatacenterClaimLister pausingv1alpha1listers.ScyllaDBDatacenterClaimLister
	scyllaDBDatacenterLister      scyllav1alpha1listers.ScyllaDBDatacenterLister

	cachesToSync []cache.InformerSynced

	eventRecorder record.EventRecorder

	queue    workqueue.RateLimitingInterface
	handlers *socontrollerhelpers.Handlers[*pausingv1alpha1.ScyllaDBDatacenterPool]
}

func NewController(
	kubeClient kubernetes.Interface,
	pausableScyllaDBClient pausingv1alpha1client.PausingV1alpha1Interface,
	scyllaClient scyllav1alpha1client.ScyllaV1alpha1Interface,
	scyllaDBDatacenterPoolInformer pausingv1alpha1informers.ScyllaDBDatacenterPoolInformer,
	scyllaDBDatacenterClaimInformer pausingv1alpha1informers.ScyllaDBDatacenterClaimInformer,
	scyllaDBDatacenterInformer scyllav1alpha1informers.ScyllaDBDatacenterInformer,
) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&corev1client.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	scpc := &Controller{
		kubeClient:    kubeClient,
		pausingClient: pausableScyllaDBClient,
		scyllaClient:  scyllaClient,

		scyllaDBDatacenterPoolLister:  scyllaDBDatacenterPoolInformer.Lister(),
		scyllaDBDatacenterClaimLister: scyllaDBDatacenterClaimInformer.Lister(),
		scyllaDBDatacenterLister:      scyllaDBDatacenterInformer.Lister(),

		cachesToSync: []cache.InformerSynced{
			scyllaDBDatacenterPoolInformer.Informer().HasSynced,
			scyllaDBDatacenterClaimInformer.Informer().HasSynced,
			scyllaDBDatacenterInformer.Informer().HasSynced,
		},

		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "scylladbdatacenterpool-controller"}),

		queue: workqueue.NewRateLimitingQueueWithConfig(
			workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{
				Name: "scylladbdatacenterpool",
			},
		),
	}

	var err error
	scpc.handlers, err = socontrollerhelpers.NewHandlers[*pausingv1alpha1.ScyllaDBDatacenterPool](
		scpc.queue,
		keyFunc,
		scheme.Scheme,
		pausingv1alpha1.ScyllaDBDatacenterPoolControllerGVK,
		kubeinterfaces.NamespacedGetList[*pausingv1alpha1.ScyllaDBDatacenterPool]{
			GetFunc: func(namespace, name string) (*pausingv1alpha1.ScyllaDBDatacenterPool, error) {
				return scpc.scyllaDBDatacenterPoolLister.ScyllaDBDatacenterPools(namespace).Get(name)
			},
			ListFunc: func(namespace string, selector labels.Selector) (ret []*pausingv1alpha1.ScyllaDBDatacenterPool, err error) {
				return scpc.scyllaDBDatacenterPoolLister.ScyllaDBDatacenterPools(namespace).List(selector)
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("can't create handlers: %w", err)
	}

	scyllaDBDatacenterPoolInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    scpc.addScyllaDBDatacenterPool,
		UpdateFunc: scpc.updateScyllaDBDatacenterPool,
		DeleteFunc: scpc.deleteScyllaDBDatacenterPool,
	})

	scyllaDBDatacenterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    scpc.addScyllaDBDatacenter,
		UpdateFunc: scpc.updateScyllaDBDatacenter,
		DeleteFunc: scpc.deleteScyllaDBDatacenter,
	})

	scyllaDBDatacenterClaimInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    scpc.addScyllaDBDatacenterClaim,
		UpdateFunc: scpc.updateScyllaDBDatacenterClaim,
		DeleteFunc: scpc.deleteScyllaDBDatacenterClaim,
	})

	return scpc, nil
}

func (sdcpc *Controller) processNextItem(ctx context.Context) bool {
	key, quit := sdcpc.queue.Get()
	if quit {
		return false
	}
	defer sdcpc.queue.Done(key)

	err := sdcpc.sync(ctx, key.(string))
	// TODO: Do smarter filtering then just Reduce to handle cases like 2 conflict errors.
	err = utilerrors.Reduce(err)
	switch {
	case err == nil:
		sdcpc.queue.Forget(key)
		return true

	case apierrors.IsConflict(err):
		klog.V(2).InfoS("Hit conflict, will retry in a bit", "Key", key, "Error", err)

	case apierrors.IsAlreadyExists(err):
		klog.V(2).InfoS("Hit already exists, will retry in a bit", "Key", key, "Error", err)

	default:
		utilruntime.HandleError(fmt.Errorf("syncing key '%v' failed: %v", key, err))
	}

	sdcpc.queue.AddRateLimited(key)

	return true
}

func (sdcpc *Controller) runWorker(ctx context.Context) {
	for sdcpc.processNextItem(ctx) {
	}
}

func (sdcpc *Controller) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()

	klog.InfoS("Starting controller", "controller", ControllerName)

	var wg sync.WaitGroup
	defer func() {
		klog.InfoS("Shutting down controller", "controller", ControllerName)
		sdcpc.queue.ShutDown()
		wg.Wait()
		klog.InfoS("Shut down controller", "controller", ControllerName)
	}()

	if !cache.WaitForNamedCacheSync(ControllerName, ctx.Done(), sdcpc.cachesToSync...) {
		return
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wait.UntilWithContext(ctx, sdcpc.runWorker, time.Second)
		}()
	}

	<-ctx.Done()
}

func (sdcpc *Controller) addScyllaDBDatacenterPool(obj interface{}) {
	sdcpc.handlers.HandleAdd(
		obj.(*pausingv1alpha1.ScyllaDBDatacenterPool),
		sdcpc.handlers.Enqueue,
	)
}

func (sdcpc *Controller) updateScyllaDBDatacenterPool(old, cur interface{}) {
	sdcpc.handlers.HandleUpdate(
		old.(*pausingv1alpha1.ScyllaDBDatacenterPool),
		cur.(*pausingv1alpha1.ScyllaDBDatacenterPool),
		sdcpc.handlers.Enqueue,
		sdcpc.deleteScyllaDBDatacenterPool,
	)
}

func (sdcpc *Controller) deleteScyllaDBDatacenterPool(obj interface{}) {
	sdcpc.handlers.HandleDelete(
		obj,
		sdcpc.handlers.Enqueue,
	)
}

func (sdcpc *Controller) addScyllaDBDatacenter(obj interface{}) {
	sdcpc.handlers.HandleAdd(
		obj.(*scyllav1alpha1.ScyllaDBDatacenter),
		sdcpc.handlers.EnqueueOwner,
	)
}

func (sdcpc *Controller) updateScyllaDBDatacenter(old, cur interface{}) {
	sdcpc.handlers.HandleUpdate(
		old.(*scyllav1alpha1.ScyllaDBDatacenter),
		cur.(*scyllav1alpha1.ScyllaDBDatacenter),
		sdcpc.handlers.EnqueueOwner,
		sdcpc.deleteScyllaDBDatacenter,
	)
}

func (sdcpc *Controller) deleteScyllaDBDatacenter(obj interface{}) {
	sdcpc.handlers.HandleDelete(
		obj,
		sdcpc.handlers.EnqueueOwner,
	)
}

func (sdcpc *Controller) addScyllaDBDatacenterClaim(obj interface{}) {
	sdcpc.handlers.HandleAdd(
		obj.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		sdcpc.enqueueScyllaDBDatacenterPoolUsingScyllaDBDatacenterClaim,
	)
}

func (sdcpc *Controller) updateScyllaDBDatacenterClaim(old, cur interface{}) {
	sdcpc.handlers.HandleUpdate(
		old.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		cur.(*pausingv1alpha1.ScyllaDBDatacenterClaim),
		sdcpc.enqueueScyllaDBDatacenterPoolUsingScyllaDBDatacenterClaim,
		sdcpc.deleteScyllaDBDatacenterClaim,
	)
}

func (sdcpc *Controller) deleteScyllaDBDatacenterClaim(obj interface{}) {
	sdcpc.handlers.HandleDelete(
		obj,
		sdcpc.enqueueScyllaDBDatacenterPoolUsingScyllaDBDatacenterClaim,
	)
}

func (sdcpc *Controller) enqueueScyllaDBDatacenterPoolUsingScyllaDBDatacenterClaim(depth int, obj kubeinterfaces.ObjectInterface, op socontrollerhelpers.HandlerOperationType) {
	sdcc := obj.(*pausingv1alpha1.ScyllaDBDatacenterClaim)

	isPoolReferencedByClaim := func(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) bool {
		return sdcp.Namespace == sdcc.Namespace &&
			sdcp.Name == sdcc.Spec.ScyllaDBDatacenterPoolName
	}

	sdcpc.handlers.EnqueueAllFunc(sdcpc.handlers.EnqueueWithFilterFunc(isPoolReferencedByClaim))(depth+1, obj, op)
}
