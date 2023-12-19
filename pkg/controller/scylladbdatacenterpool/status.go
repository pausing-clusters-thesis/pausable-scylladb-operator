//  Copyright (C) 2024 ScyllaDB

package scylladbdatacenterpool

import (
	"context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/utils/pointer"
)

func (sdcpc *Controller) calculateStatus(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) *pausingv1alpha1.ScyllaDBDatacenterPoolStatus {
	status := sdcp.Status.DeepCopy()
	status.ObservedGeneration = pointer.Int64(sdcp.Generation)

	return status
}

func (sdcpc *Controller) updateStatus(ctx context.Context, currentSCP *pausingv1alpha1.ScyllaDBDatacenterPool, status *pausingv1alpha1.ScyllaDBDatacenterPoolStatus) error {
	if apiequality.Semantic.DeepEqual(&currentSCP.Status, status) {
		return nil
	}

	return nil
}
