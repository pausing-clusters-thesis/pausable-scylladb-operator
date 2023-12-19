package scylladbdatacenterpool

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	soslices "github.com/scylladb/scylla-operator/pkg/helpers/slices"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func (sdcpc *Controller) sync(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return fmt.Errorf("can't split meta namespace key: %w", err)
	}

	now := time.Now()
	klog.V(4).InfoS("Started syncing ScyllaDBDatacenterPool", "ScyllaDBDatacenterPool", klog.KRef(namespace, name), "startTime", now)
	defer func() {
		klog.V(4).InfoS("Finished syncing ScyllaDBDatacenterPool", "ScyllaDBDatacenterPool", klog.KRef(namespace, name), "duration", time.Since(now))
	}()

	sdcp, err := sdcpc.scyllaDBDatacenterPoolLister.ScyllaDBDatacenterPools(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("ScyllaDBDatacenterPool has been deleted", "ScyllaDBDatacenterPool", klog.KObj(sdcp))
		return nil
	}
	if err != nil {
		return fmt.Errorf("can't get ScyllaDBDatacenterPool %q: %w", naming.ManualRef(namespace, name), err)
	}

	scyllaDBDatacenterClaims, err := sdcpc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(sdcp.Namespace).List(labels.Everything())
	if err != nil {
		return fmt.Errorf("can't list ScyllaDBDatacenterClaims: %w", err)
	}

	scyllaDBDatacenterClaims = soslices.Filter(scyllaDBDatacenterClaims, func(scc *pausingv1alpha1.ScyllaDBDatacenterClaim) bool {
		return scc.Spec.ScyllaDBDatacenterPoolName == sdcp.Name
	})

	selector := naming.ScyllaDBDatacenterPoolSelector(sdcp)

	var objectErrs []error
	scyllaDBDatacenterMap, err := socontrollerhelpers.GetObjects[*pausingv1alpha1.ScyllaDBDatacenterPool, *scyllav1alpha1.ScyllaDBDatacenter](
		ctx,
		sdcp,
		pausingv1alpha1.ScyllaDBDatacenterPoolControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[*pausingv1alpha1.ScyllaDBDatacenterPool, *scyllav1alpha1.ScyllaDBDatacenter]{
			GetControllerUncachedFunc: sdcpc.pausingClient.ScyllaDBDatacenterPools(sdcp.Namespace).Get,
			ListObjectsFunc:           sdcpc.scyllaDBDatacenterLister.ScyllaDBDatacenters(sdcp.Namespace).List,
			PatchObjectFunc:           sdcpc.scyllaClient.ScyllaDBDatacenters(sdcp.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, err)
	}

	objectErr := utilerrors.NewAggregate(objectErrs)
	if objectErr != nil {
		return objectErr
	}

	status := sdcpc.calculateStatus(sdcp, slices.Collect(maps.Values(scyllaDBDatacenterMap)))

	if sdcp.DeletionTimestamp != nil {
		return sdcpc.updateStatus(ctx, sdcp, status)
	}

	var errs []error

	err = socontrollerhelpers.RunSync(
		&status.Conditions,
		scyllaDBDatacenterControllerProgressingCondition,
		scyllaDBDatacenterControllerDegradedCondition,
		sdcp.Generation,
		func() ([]metav1.Condition, error) {
			return sdcpc.syncScyllaDBDatacenters(ctx, key, sdcp, scyllaDBDatacenterClaims, scyllaDBDatacenterMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync ScyllaDBDatacenters: %w", err))
	}

	// Aggregate conditions.
	var aggregationErrs []error
	availableConditions, err := socontrollerhelpers.AggregateStatusConditions(
		socontrollerhelpers.FindStatusConditionsWithSuffix(status.Conditions, pausingv1alpha1.AvailableCondition),
		metav1.Condition{
			Type:               pausingv1alpha1.AvailableCondition,
			Status:             metav1.ConditionTrue,
			Reason:             internalapi.AsExpectedReason,
			Message:            "",
			ObservedGeneration: sdcp.Generation,
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
			ObservedGeneration: sdcp.Generation,
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
			ObservedGeneration: sdcp.Generation,
		},
	)
	if err != nil {
		aggregationErrs = append(aggregationErrs, fmt.Errorf("can't aggregate degraded status conditions: %w", err))
	}

	if len(aggregationErrs) > 0 {
		errs = append(errs, aggregationErrs...)
		return utilerrors.NewAggregate(errs)
	}

	apimeta.SetStatusCondition(&status.Conditions, availableConditions)
	apimeta.SetStatusCondition(&status.Conditions, progressingCondition)
	apimeta.SetStatusCondition(&status.Conditions, degradedCondition)

	err = sdcpc.updateStatus(ctx, sdcp, status)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}
