package controllerhelpers

import (
	"fmt"

	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
)

func RunSync(conditions *[]metav1.Condition, availableConditionType, progressingConditionType, degradedCondType string, observedGeneration int64, syncFn func() ([]metav1.Condition, []metav1.Condition, error)) error {
	availableConditions, progressingConditions, err := syncFn()
	socontrollerhelpers.SetStatusConditionFromError(conditions, err, degradedCondType, observedGeneration)
	if err != nil {
		return err
	}

	var errs []error

	availableCondition, err := socontrollerhelpers.AggregateStatusConditions(
		availableConditions,
		metav1.Condition{
			Type:               availableConditionType,
			Status:             metav1.ConditionTrue,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: observedGeneration,
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't aggregate available conditions %q: %w", availableConditionType, err))
	}

	progressingCondition, err := socontrollerhelpers.AggregateStatusConditions(
		progressingConditions,
		metav1.Condition{
			Type:               progressingConditionType,
			Status:             metav1.ConditionFalse,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: observedGeneration,
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't aggregate progressing conditions %q: %w", progressingConditionType, err))
	}

	err = apierrors.NewAggregate(errs)
	if err != nil {
		return err
	}

	apimeta.SetStatusCondition(conditions, availableCondition)
	apimeta.SetStatusCondition(conditions, progressingCondition)

	return nil
}
