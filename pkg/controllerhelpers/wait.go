package controllerhelpers

import (
	"context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	pausingv1alpha1client "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned/typed/pausing/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
)

func WaitForScyllaDBDatacenterPoolState(ctx context.Context, client pausingv1alpha1client.ScyllaDBDatacenterPoolInterface, name string, options socontrollerhelpers.WaitForStateOptions, condition func(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) (bool, error), additionalConditions ...func(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) (bool, error)) (*pausingv1alpha1.ScyllaDBDatacenterPool, error) {
	return socontrollerhelpers.WaitForObjectState[*pausingv1alpha1.ScyllaDBDatacenterPool, *pausingv1alpha1.ScyllaDBDatacenterPoolList](ctx, client, name, options, condition, additionalConditions...)
}

func WaitForPausableScyllaDBDatacenterState(ctx context.Context, client pausingv1alpha1client.PausableScyllaDBDatacenterInterface, name string, options socontrollerhelpers.WaitForStateOptions, condition func(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) (bool, error), additionalConditions ...func(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) (bool, error)) (*pausingv1alpha1.PausableScyllaDBDatacenter, error) {
	return socontrollerhelpers.WaitForObjectState[*pausingv1alpha1.PausableScyllaDBDatacenter, *pausingv1alpha1.PausableScyllaDBDatacenterList](ctx, client, name, options, condition, additionalConditions...)
}
