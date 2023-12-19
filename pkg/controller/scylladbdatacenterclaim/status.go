//  Copyright (C) 2024 ScyllaDB

package scylladbdatacenterclaim

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func (sdccc *Controller) calculateStatus(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter) *pausingv1alpha1.ScyllaDBDatacenterClaimStatus {
	status := sdcc.Status.DeepCopy()
	status.ObservedGeneration = ptr.To(sdcc.Generation)

	scyllaDBDatacenterControllerAvailableCondition := getScyllaDBDatacenterControllerAvailableCondition(sdcc, scyllaDBDatacenters)
	apimeta.SetStatusCondition(&status.Conditions, scyllaDBDatacenterControllerAvailableCondition)

	return status
}

func getScyllaDBDatacenterControllerAvailableCondition(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter) metav1.Condition {
	if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
		return metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcc.Generation,
			Reason:             "AwaitingScyllaDBDatacenterBinding",
			Message:            "Waiting for a ScyllaDBDatacenter to be bound.",
		}
	}

	sdcName := *sdcc.Status.ScyllaDBDatacenterName
	sdc, ok := scyllaDBDatacenters[sdcName]
	if !ok {
		return metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcc.Generation,
			Reason:             "AwaitingScyllaDBDatacenterCreation",
			Message:            fmt.Sprintf("Waiting for ScyllaDBDatacenter %q to be created.", sonaming.ManualRef(sdcc.Namespace, sdcName)),
		}
	}

	rolledOut, err := socontrollerhelpers.IsScyllaDBDatacenterRolledOut(sdc)
	if err != nil {
		klog.ErrorS(err, "can't verify if ScyllaDBDatacenter is rolled out", "ScyllaDBDatacenterClaim", naming.ObjRef(sdcc), "ScyllaDBDatacenter", naming.ObjRef(sdc))
		return metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionUnknown,
			ObservedGeneration: sdcc.Generation,
			Reason:             "ScyllaDBDatacenterRolloutStatusUnknown",
			Message:            fmt.Sprintf("Can't determine ScyllaDBDatacenter %q rollout status.", naming.ObjRef(sdc)),
		}
	}
	if !rolledOut {
		return metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcc.Generation,
			Reason:             "AwaitingScyllaDBDatacenterRollout",
			Message:            fmt.Sprintf("Waiting for ScyllaDBDatacenter %q to be rolled out.", naming.ObjRef(sdc)),
		}
	}

	return metav1.Condition{
		Type:               scyllaDBDatacenterControllerAvailableCondition,
		Status:             metav1.ConditionTrue,
		ObservedGeneration: sdcc.Generation,
		Reason:             internalapi.AsExpectedReason,
		Message:            "",
	}
}

func (sdccc *Controller) updateStatus(ctx context.Context, currentSDCC *pausingv1alpha1.ScyllaDBDatacenterClaim, status *pausingv1alpha1.ScyllaDBDatacenterClaimStatus) error {
	if apiequality.Semantic.DeepEqual(&currentSDCC.Status, status) {
		return nil
	}

	sdcc := currentSDCC.DeepCopy()
	sdcc.Status = *status

	klog.V(2).InfoS("Updating status", "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

	_, err := sdccc.pausableScyllaDBClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdcc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	klog.V(2).InfoS("Status updated", "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

	return nil
}
