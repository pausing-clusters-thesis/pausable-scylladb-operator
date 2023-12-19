package resourceapply

import (
	"context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	pausingv1alpha1client "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned/typed/pausing/v1alpha1"
	pausablescylladbv1alpha1lister "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/listers/pausing/v1alpha1"
	"github.com/scylladb/scylla-operator/pkg/resourceapply"
	"k8s.io/client-go/tools/record"
)

func ApplyScyllaDBDatacenterClaimWithControl(
	ctx context.Context,
	control resourceapply.ApplyControlInterface[*pausingv1alpha1.ScyllaDBDatacenterClaim],
	recorder record.EventRecorder,
	required *pausingv1alpha1.ScyllaDBDatacenterClaim,
	options resourceapply.ApplyOptions,
) (*pausingv1alpha1.ScyllaDBDatacenterClaim, bool, error) {
	return resourceapply.ApplyGeneric[*pausingv1alpha1.ScyllaDBDatacenterClaim](ctx, control, recorder, required, options)
}

func ApplyScyllaDBDatacenterClaim(
	ctx context.Context,
	client pausingv1alpha1client.ScyllaDBDatacenterClaimsGetter,
	lister pausablescylladbv1alpha1lister.ScyllaDBDatacenterClaimLister,
	recorder record.EventRecorder,
	required *pausingv1alpha1.ScyllaDBDatacenterClaim,
	options resourceapply.ApplyOptions,
) (*pausingv1alpha1.ScyllaDBDatacenterClaim, bool, error) {
	return ApplyScyllaDBDatacenterClaimWithControl(
		ctx,
		resourceapply.ApplyControlFuncs[*pausingv1alpha1.ScyllaDBDatacenterClaim]{
			GetCachedFunc: lister.ScyllaDBDatacenterClaims(required.Namespace).Get,
			CreateFunc:    client.ScyllaDBDatacenterClaims(required.Namespace).Create,
			UpdateFunc:    client.ScyllaDBDatacenterClaims(required.Namespace).Update,
			DeleteFunc:    client.ScyllaDBDatacenterClaims(required.Namespace).Delete,
		},
		recorder,
		required,
		options,
	)
}
