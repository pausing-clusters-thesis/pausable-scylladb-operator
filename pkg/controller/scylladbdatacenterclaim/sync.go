package scylladbdatacenterclaim

import (
	"context"
	"fmt"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func (sdccc *Controller) sync(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return fmt.Errorf("can't split meta namespace key: %w", err)
	}

	now := time.Now()
	klog.V(4).InfoS("Started syncing ScyllaDBDatacenterClaim", "ScyllaDBDatacenterClaim", klog.KRef(namespace, name), "startTime", now)
	defer func() {
		klog.V(4).InfoS("Finished syncing ScyllaDBDatacenterClaim", "ScyllaDBDatacenterClaim", klog.KRef(namespace, name), "duration", time.Since(now))
	}()

	sdcc, err := sdccc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			return fmt.Errorf("can't get ScyllaDBDatacenterClaim %q: %w", naming.ManualRef(namespace, name), err)

		}

		klog.V(2).InfoS("ScyllaDBDatacenterClaim has been deleted", "ScyllaDBDatacenterClaim", klog.KObj(sdcc))
		return nil
	}

	selector := naming.ScyllaDBDatacenterClaimSelector(sdcc)

	var objectErrs []error
	scyllaDBDatacenters, err := socontrollerhelpers.GetObjects[*pausingv1alpha1.ScyllaDBDatacenterClaim, *scyllav1alpha1.ScyllaDBDatacenter](
		ctx,
		sdcc,
		pausingv1alpha1.ScyllaDBDatacenterClaimControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[*pausingv1alpha1.ScyllaDBDatacenterClaim, *scyllav1alpha1.ScyllaDBDatacenter]{
			GetControllerUncachedFunc: sdccc.pausableScyllaDBClient.ScyllaDBDatacenterClaims(sdcc.Namespace).Get,
			ListObjectsFunc:           sdccc.scyllaDBDatacenterLister.ScyllaDBDatacenters(sdcc.Namespace).List,
			PatchObjectFunc:           sdccc.scyllaClient.ScyllaDBDatacenters(sdcc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get scyllaDBDatacenters: %w", err))
	}

	status := sdccc.calculateStatus(sdcc, scyllaDBDatacenters)

	if sdcc.DeletionTimestamp != nil {
		return sdccc.updateStatus(ctx, sdcc, status)
	}

	objectErr := utilerrors.NewAggregate(objectErrs)
	if objectErr != nil {
		return objectErr
	}

	var errs []error

	err = controllerhelpers.RunSync(
		&status.Conditions,
		scyllaDBDatacenterControllerAvailableCondition,
		scyllaDBDatacenterControllerProgressingCondition,
		scyllaDBDatacenterControllerDegradedCondition,
		sdcc.Generation,
		func() ([]metav1.Condition, []metav1.Condition, error) {
			return sdccc.syncScyllaDBDatacenter(ctx, sdcc, scyllaDBDatacenters)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync ScyllaDBDatacenters: %w", err))
	}

	// Aggregate conditions.
	var aggregationErrs []error
	availableCondition, err := socontrollerhelpers.AggregateStatusConditions(
		socontrollerhelpers.FindStatusConditionsWithSuffix(status.Conditions, scyllav1.AvailableCondition),
		metav1.Condition{
			Type:               scyllav1.AvailableCondition,
			Status:             metav1.ConditionTrue,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: sdcc.Generation,
		},
	)
	if err != nil {
		aggregationErrs = append(aggregationErrs, fmt.Errorf("can't aggregate available status conditions: %w", err))
	}

	progressingCondition, err := socontrollerhelpers.AggregateStatusConditions(
		socontrollerhelpers.FindStatusConditionsWithSuffix(status.Conditions, pausingv1alpha1.ProgressingCondition),
		metav1.Condition{
			Type:               pausingv1alpha1.ProgressingCondition,
			Status:             metav1.ConditionFalse,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: sdcc.Generation,
		},
	)
	if err != nil {
		aggregationErrs = append(aggregationErrs, fmt.Errorf("can't aggregate progressing status conditions: %w", err))
	}

	degradedCondition, err := socontrollerhelpers.AggregateStatusConditions(
		socontrollerhelpers.FindStatusConditionsWithSuffix(status.Conditions, pausingv1alpha1.DegradedCondition),
		metav1.Condition{
			Type:               pausingv1alpha1.DegradedCondition,
			Status:             metav1.ConditionFalse,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: sdcc.Generation,
		},
	)
	if err != nil {
		aggregationErrs = append(aggregationErrs, fmt.Errorf("can't aggregate degraded status conditions: %w", err))
	}

	if len(aggregationErrs) > 0 {
		errs = append(errs, aggregationErrs...)
		return utilerrors.NewAggregate(errs)
	}

	apimeta.SetStatusCondition(&status.Conditions, availableCondition)
	apimeta.SetStatusCondition(&status.Conditions, progressingCondition)
	apimeta.SetStatusCondition(&status.Conditions, degradedCondition)

	err = sdccc.updateStatus(ctx, sdcc, status)
	if err != nil {
		return fmt.Errorf("can't update status: %w", err)
	}

	return nil
}
