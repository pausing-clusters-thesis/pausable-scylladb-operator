// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// PausableScyllaDBDatacenters returns a PausableScyllaDBDatacenterInformer.
	PausableScyllaDBDatacenters() PausableScyllaDBDatacenterInformer
	// ScyllaDBDatacenterClaims returns a ScyllaDBDatacenterClaimInformer.
	ScyllaDBDatacenterClaims() ScyllaDBDatacenterClaimInformer
	// ScyllaDBDatacenterPools returns a ScyllaDBDatacenterPoolInformer.
	ScyllaDBDatacenterPools() ScyllaDBDatacenterPoolInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// PausableScyllaDBDatacenters returns a PausableScyllaDBDatacenterInformer.
func (v *version) PausableScyllaDBDatacenters() PausableScyllaDBDatacenterInformer {
	return &pausableScyllaDBDatacenterInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ScyllaDBDatacenterClaims returns a ScyllaDBDatacenterClaimInformer.
func (v *version) ScyllaDBDatacenterClaims() ScyllaDBDatacenterClaimInformer {
	return &scyllaDBDatacenterClaimInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ScyllaDBDatacenterPools returns a ScyllaDBDatacenterPoolInformer.
func (v *version) ScyllaDBDatacenterPools() ScyllaDBDatacenterPoolInformer {
	return &scyllaDBDatacenterPoolInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
