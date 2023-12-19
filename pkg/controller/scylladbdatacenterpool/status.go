package scylladbdatacenterpool

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	sointernalapi "github.com/scylladb/scylla-operator/pkg/internalapi"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
)

func (sdcpc *Controller) calculateStatus(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool, scyllaDBDatacenters []*scyllav1alpha1.ScyllaDBDatacenter) *pausingv1alpha1.ScyllaDBDatacenterPoolStatus {
	status := sdcp.Status.DeepCopy()
	status.ObservedGeneration = pointer.Int64(sdcp.Generation)

	instances := int32(0)
	readyInstances := int32(0)
	for _, sdc := range scyllaDBDatacenters {
		instances++

		isScyllaDBDatacenterPrewarmed := controllerhelpers.IsScyllaDBDatacenterPrewarmed(sdc)
		if isScyllaDBDatacenterPrewarmed {
			readyInstances++
		}
	}

	status.Instances = &instances
	status.ReadyInstances = &readyInstances

	switch {
	case instances < sdcp.Spec.Capacity:
		apimeta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcp.Generation,
			Reason:             "TooFewInstances",
			Message:            fmt.Sprintf("Only %d instances out of the desired capacity %d exist.", instances, sdcp.Spec.Capacity),
		})

	case readyInstances < sdcp.Spec.Capacity:
		apimeta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcp.Generation,
			Reason:             "TooFewReadyInstances",
			Message:            fmt.Sprintf("Only %d instances out of the desired capacity %d are prewarmed.", instances, sdcp.Spec.Capacity),
		})

	default:
		apimeta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: sdcp.Generation,
			Reason:             sointernalapi.AsExpectedReason,
			Message:            "",
		})

	}

	return status
}

func (sdcpc *Controller) updateStatus(ctx context.Context, currentSCP *pausingv1alpha1.ScyllaDBDatacenterPool, status *pausingv1alpha1.ScyllaDBDatacenterPoolStatus) error {
	if apiequality.Semantic.DeepEqual(&currentSCP.Status, status) {
		return nil
	}

	sdcp := currentSCP.DeepCopy()
	sdcp.Status = *status

	klog.V(2).InfoS("Updating status", "ScyllaDBDatacenterPool", klog.KObj(sdcp))

	_, err := sdcpc.pausingClient.ScyllaDBDatacenterPools(sdcp.Namespace).UpdateStatus(ctx, sdcp, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("can't update ScyllaDBDatacenterPool %q status: %w", sonaming.ObjRef(sdcp), err)
	}

	klog.V(2).InfoS("Status updated", "ScyllaDBDatacenterPool", klog.KObj(sdcp))

	return nil
}
