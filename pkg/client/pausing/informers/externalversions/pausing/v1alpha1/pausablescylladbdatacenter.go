// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"
	time "time"

	apipausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	versioned "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned"
	internalinterfaces "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions/internalinterfaces"
	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/listers/pausing/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PausableScyllaDBDatacenterInformer provides access to a shared informer and lister for
// PausableScyllaDBDatacenters.
type PausableScyllaDBDatacenterInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() pausingv1alpha1.PausableScyllaDBDatacenterLister
}

type pausableScyllaDBDatacenterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPausableScyllaDBDatacenterInformer constructs a new informer for PausableScyllaDBDatacenter type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPausableScyllaDBDatacenterInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPausableScyllaDBDatacenterInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPausableScyllaDBDatacenterInformer constructs a new informer for PausableScyllaDBDatacenter type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPausableScyllaDBDatacenterInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PausingV1alpha1().PausableScyllaDBDatacenters(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PausingV1alpha1().PausableScyllaDBDatacenters(namespace).Watch(context.TODO(), options)
			},
		},
		&apipausingv1alpha1.PausableScyllaDBDatacenter{},
		resyncPeriod,
		indexers,
	)
}

func (f *pausableScyllaDBDatacenterInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPausableScyllaDBDatacenterInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *pausableScyllaDBDatacenterInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apipausingv1alpha1.PausableScyllaDBDatacenter{}, f.defaultInformer)
}

func (f *pausableScyllaDBDatacenterInformer) Lister() pausingv1alpha1.PausableScyllaDBDatacenterLister {
	return pausingv1alpha1.NewPausableScyllaDBDatacenterLister(f.Informer().GetIndexer())
}
