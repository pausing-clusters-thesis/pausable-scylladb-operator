// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	scheme "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ScyllaDBDatacenterPoolsGetter has a method to return a ScyllaDBDatacenterPoolInterface.
// A group's client should implement this interface.
type ScyllaDBDatacenterPoolsGetter interface {
	ScyllaDBDatacenterPools(namespace string) ScyllaDBDatacenterPoolInterface
}

// ScyllaDBDatacenterPoolInterface has methods to work with ScyllaDBDatacenterPool resources.
type ScyllaDBDatacenterPoolInterface interface {
	Create(ctx context.Context, scyllaDBDatacenterPool *pausingv1alpha1.ScyllaDBDatacenterPool, opts v1.CreateOptions) (*pausingv1alpha1.ScyllaDBDatacenterPool, error)
	Update(ctx context.Context, scyllaDBDatacenterPool *pausingv1alpha1.ScyllaDBDatacenterPool, opts v1.UpdateOptions) (*pausingv1alpha1.ScyllaDBDatacenterPool, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, scyllaDBDatacenterPool *pausingv1alpha1.ScyllaDBDatacenterPool, opts v1.UpdateOptions) (*pausingv1alpha1.ScyllaDBDatacenterPool, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*pausingv1alpha1.ScyllaDBDatacenterPool, error)
	List(ctx context.Context, opts v1.ListOptions) (*pausingv1alpha1.ScyllaDBDatacenterPoolList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *pausingv1alpha1.ScyllaDBDatacenterPool, err error)
	ScyllaDBDatacenterPoolExpansion
}

// scyllaDBDatacenterPools implements ScyllaDBDatacenterPoolInterface
type scyllaDBDatacenterPools struct {
	*gentype.ClientWithList[*pausingv1alpha1.ScyllaDBDatacenterPool, *pausingv1alpha1.ScyllaDBDatacenterPoolList]
}

// newScyllaDBDatacenterPools returns a ScyllaDBDatacenterPools
func newScyllaDBDatacenterPools(c *PausingV1alpha1Client, namespace string) *scyllaDBDatacenterPools {
	return &scyllaDBDatacenterPools{
		gentype.NewClientWithList[*pausingv1alpha1.ScyllaDBDatacenterPool, *pausingv1alpha1.ScyllaDBDatacenterPoolList](
			"scylladbdatacenterpools",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *pausingv1alpha1.ScyllaDBDatacenterPool { return &pausingv1alpha1.ScyllaDBDatacenterPool{} },
			func() *pausingv1alpha1.ScyllaDBDatacenterPoolList {
				return &pausingv1alpha1.ScyllaDBDatacenterPoolList{}
			},
		),
	}
}
