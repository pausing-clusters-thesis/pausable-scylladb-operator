//  Copyright (C) 2024 ScyllaDB

package pausablescylladbdatacenter

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	sohelpers "github.com/scylladb/scylla-operator/pkg/helpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func (psdcc *Controller) calculateStatus(psdc *pausingv1alpha1.PausableScyllaDBDatacenter, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) *pausingv1alpha1.PausableScyllaDBDatacenterStatus {
	status := psdc.Status.DeepCopy()
	status.ObservedGeneration = ptr.To(psdc.Generation)

	scyllaDBDatacenterClaimControllerAvailableCondition := getScyllaDBDatacenterClaimControllerAvailableCondition(
		naming.GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc),
		psdc.Namespace,
		psdc.Generation,
		scyllaDBDatacenterClaims,
	)
	apimeta.SetStatusCondition(&status.Conditions, scyllaDBDatacenterClaimControllerAvailableCondition)

	return status
}

func getScyllaDBDatacenterClaimControllerAvailableCondition(name string, namespace string, generation int64, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) metav1.Condition {
	sdcc, ok := scyllaDBDatacenterClaims[name]
	if !ok {
		return metav1.Condition{
			Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: generation,
			Reason:             "AwaitingCreation",
			Message:            fmt.Sprintf("Awaiting ScyllaDBDatacenterClaim %q to be created.", sonaming.ManualRef(namespace, name)),
		}
	}

	if !isScyllaDBDatacenterClaimRolledOut(sdcc) {
		return metav1.Condition{
			Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcc.Generation,
			Reason:             "AwaitingScyllaDBDatacenterClaimRollout",
			Message:            fmt.Sprintf("Waiting for ScyllaDBDatacenterClaim %q to be rolled out.", naming.ObjRef(sdcc)),
		}
	}

	//availableCondition, _, availableConditionPresent := soslices.Find(sdcc.Status.Conditions, func(c metav1.Condition) bool {
	//	return c.Type == scyllav1.AvailableCondition
	//})
	//if !availableConditionPresent {
	//	return metav1.Condition{
	//		Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
	//		Status:             metav1.ConditionFalse,
	//		ObservedGeneration: generation,
	//		Reason:             "AwaitingCondition",
	//		Message:            fmt.Sprintf("Awaiting ScyllaDBDatacenterClaim %q Available condition to be present.", sonaming.ObjRef(sdcc)),
	//	}
	//}
	//
	//if availableCondition.Status == metav1.ConditionFalse {
	//	return metav1.Condition{
	//		Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
	//		Status:             metav1.ConditionFalse,
	//		ObservedGeneration: generation,
	//		Reason:             availableCondition.Reason,
	//		Message:            availableCondition.Message,
	//	}
	//}

	return metav1.Condition{
		Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
		Status:             metav1.ConditionTrue,
		ObservedGeneration: generation,
		Reason:             internalapi.AsExpectedReason,
		Message:            "",
	}
}

func isScyllaDBDatacenterClaimRolledOut(scc *pausingv1alpha1.ScyllaDBDatacenterClaim) bool {
	if !sohelpers.IsStatusConditionPresentAndTrue(scc.Status.Conditions, pausingv1alpha1.AvailableCondition, scc.Generation) {
		return false
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(scc.Status.Conditions, pausingv1alpha1.ProgressingCondition, scc.Generation) {
		return false
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(scc.Status.Conditions, pausingv1alpha1.DegradedCondition, scc.Generation) {
		return false
	}

	return true
}

func (psdcc *Controller) updateStatus(ctx context.Context, curPSDC *pausingv1alpha1.PausableScyllaDBDatacenter, status *pausingv1alpha1.PausableScyllaDBDatacenterStatus) error {
	if apiequality.Semantic.DeepEqual(&curPSDC.Status, status) {
		return nil
	}

	scc := curPSDC.DeepCopy()
	scc.Status = *status

	klog.V(2).InfoS("Updating status", "PausableScyllaDBDatacenter", klog.KObj(scc))

	_, err := psdcc.pausingClient.PausableScyllaDBDatacenters(scc.Namespace).UpdateStatus(ctx, scc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	klog.V(2).InfoS("Status updated", "PausableScyllaDBDatacenter", klog.KObj(scc))

	return nil
}
