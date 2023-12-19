package controllerhelpers

import (
	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	sohelpers "github.com/scylladb/scylla-operator/pkg/helpers"
)

func IsScyllaDBDatacenterClaimRolledOut(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim) (bool, error) {
	if !sohelpers.IsStatusConditionPresentAndTrue(sdcc.Status.Conditions, pausingv1alpha1.AvailableCondition, sdcc.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(sdcc.Status.Conditions, pausingv1alpha1.ProgressingCondition, sdcc.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(sdcc.Status.Conditions, pausingv1alpha1.DegradedCondition, sdcc.Generation) {
		return false, nil
	}

	return true, nil
}

func IsPausableScyllaDBDatacenterRolledOut(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) (bool, error) {
	if !sohelpers.IsStatusConditionPresentAndTrue(psdc.Status.Conditions, pausingv1alpha1.AvailableCondition, psdc.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(psdc.Status.Conditions, pausingv1alpha1.ProgressingCondition, psdc.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(psdc.Status.Conditions, pausingv1alpha1.DegradedCondition, psdc.Generation) {
		return false, nil
	}

	return true, nil
}

func IsPausableScyllaDBDatacenterAvailable(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) (bool, error) {
	if !sohelpers.IsStatusConditionPresentAndTrue(psdc.Status.Conditions, pausingv1alpha1.AvailableCondition, psdc.Generation) {
		return false, nil
	}

	return true, nil
}

func IsScyllaDBDatacenterPoolRolledOut(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) (bool, error) {
	if !sohelpers.IsStatusConditionPresentAndTrue(sdcp.Status.Conditions, pausingv1alpha1.AvailableCondition, sdcp.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(sdcp.Status.Conditions, pausingv1alpha1.ProgressingCondition, sdcp.Generation) {
		return false, nil
	}

	if !sohelpers.IsStatusConditionPresentAndFalse(sdcp.Status.Conditions, pausingv1alpha1.DegradedCondition, sdcp.Generation) {
		return false, nil
	}

	return true, nil
}
