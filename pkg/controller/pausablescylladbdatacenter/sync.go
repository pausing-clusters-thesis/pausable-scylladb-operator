package pausablescylladbdatacenter

import (
	"context"
	"fmt"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func getSelectorLabels(psc *pausingv1alpha1.PausableScyllaDBDatacenter) labels.Set {
	return labels.Set{
		naming.PausableScyllaDBDatacenterNameLabel: psc.Name,
	}
}

func getSelector(psc *pausingv1alpha1.PausableScyllaDBDatacenter) labels.Selector {
	return labels.SelectorFromSet(getSelectorLabels(psc))
}

func (psdcc *Controller) sync(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return fmt.Errorf("can't split meta namespace key: %w", err)
	}

	now := time.Now()
	klog.V(4).InfoS("Started syncing PausableScyllaDBDatacenter", "PausableScyllaDBDatacenter", klog.KRef(namespace, name), "startTime", now)
	defer func() {
		klog.V(4).InfoS("Finished syncing PausableScyllaDBDatacenter", "PausableScyllaDBDatacenter", klog.KRef(namespace, name), "duration", time.Since(now))
	}()

	psdc, err := psdcc.pausableScyllaDBDatacenterLister.PausableScyllaDBDatacenters(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("PausableScyllaDBDatacenters has been deleted", "PausableScyllaDBDatacenters", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return fmt.Errorf("can't get PausableScyllaDBDatacenter %q: %w", naming.ManualRef(namespace, name), err)
	}

	selector := getSelector(psdc)
	type CT = *pausingv1alpha1.PausableScyllaDBDatacenter
	var objectErrs []error

	persistentVolumeClaimMap, err := socontrollerhelpers.GetObjects[CT, *corev1.PersistentVolumeClaim](
		ctx,
		psdc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *corev1.PersistentVolumeClaim]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psdc.Namespace).Get,
			ListObjectsFunc:           psdcc.persistentVolumeClaimLister.PersistentVolumeClaims(psdc.Namespace).List,
			PatchObjectFunc:           psdcc.kubeClient.CoreV1().PersistentVolumeClaims(psdc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get persistentvolumeclaims: %w", err))
	}

	scyllaDBDatacenterClaimMap, err := socontrollerhelpers.GetObjects[CT, *pausingv1alpha1.ScyllaDBDatacenterClaim](
		ctx,
		psdc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *pausingv1alpha1.ScyllaDBDatacenterClaim]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psdc.Namespace).Get,
			ListObjectsFunc:           psdcc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(psdc.Namespace).List,
			PatchObjectFunc:           psdcc.pausingClient.ScyllaDBDatacenterClaims(psdc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get ScyllaDBDatacenterClaims: %w", err))
	}

	secretMap, err := socontrollerhelpers.GetObjects[CT, *corev1.Secret](
		ctx,
		psdc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *corev1.Secret]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psdc.Namespace).Get,
			ListObjectsFunc:           psdcc.secretLister.Secrets(psdc.Namespace).List,
			PatchObjectFunc:           psdcc.kubeClient.CoreV1().Secrets(psdc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get secrets: %w", err))
	}

	objectErr := utilerrors.NewAggregate(objectErrs)
	if objectErr != nil {
		return objectErr
	}

	status := psdcc.calculateStatus(psdc, persistentVolumeClaimMap, scyllaDBDatacenterClaimMap)

	if psdc.DeletionTimestamp != nil {
		return psdcc.updateStatus(ctx, psdc, status)
	}

	var errs []error

	err = controllerhelpers.RunSync(
		&status.Conditions,
		persistentVolumeClaimControllerAvailableCondition,
		persistentVolumeClaimControllerProgressingCondition,
		persistentVolumeClaimControllerDegradedCondition,
		psdc.Generation,
		func() ([]metav1.Condition, []metav1.Condition, error) {
			return psdcc.syncPersistentVolumeClaims(ctx, psdc, persistentVolumeClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync PersistentVolumeClaims: %w", err))
	}

	err = controllerhelpers.RunSync(
		&status.Conditions,
		scyllaDBDatacenterClaimControllerAvailableCondition,
		scyllaDBDatacenterClaimControllerProgressingCondition,
		scyllaDBDatacenterClaimControllerDegradedCondition,
		psdc.Generation,
		func() ([]metav1.Condition, []metav1.Condition, error) {
			return psdcc.syncScyllaDBDatacenterClaims(ctx, psdc, scyllaDBDatacenterClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync ScyllaDBDatacenterClaims: %w", err))
	}

	err = controllerhelpers.RunSync(
		&status.Conditions,
		delayedVolumeMountControllerAvailableCondition,
		delayedVolumeMountControllerProgressingCondition,
		delayedVolumeMountControllerDegradedCondition,
		psdc.Generation,
		func() ([]metav1.Condition, []metav1.Condition, error) {
			return psdcc.syncDelayedVolumeMount(ctx, psdc, scyllaDBDatacenterClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync delayed volume mounting: %w", err))
	}

	err = controllerhelpers.RunSync(
		&status.Conditions,
		certControllerAvailableCondition,
		certControllerProgressingCondition,
		certControllerDegradedCondition,
		psdc.Generation,
		func() ([]metav1.Condition, []metav1.Condition, error) {
			return psdcc.syncCerts(ctx, psdc, secretMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync certs: %w", err))
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
			ObservedGeneration: psdc.Generation,
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
			ObservedGeneration: psdc.Generation,
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
			ObservedGeneration: psdc.Generation,
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

	err = psdcc.updateStatus(ctx, psdc, status)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}
