//  Copyright (C) 2024 ScyllaDB

package pausablescylladbdatacenter

import (
	"context"
	"fmt"
	"time"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
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

	psc, err := psdcc.pausableScyllaDBDatacenterLister.PausableScyllaDBDatacenters(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("PausableScyllaDBDatacenters has been deleted", "PausableScyllaDBDatacenters", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return fmt.Errorf("can't get PausableScyllaDBDatacenter %q: %w", naming.ManualRef(namespace, name), err)
	}

	selector := getSelector(psc)
	type CT = *pausingv1alpha1.PausableScyllaDBDatacenter
	var objectErrs []error

	persistentVolumeClaimMap, err := socontrollerhelpers.GetObjects[CT, *corev1.PersistentVolumeClaim](
		ctx,
		psc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *corev1.PersistentVolumeClaim]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psc.Namespace).Get,
			ListObjectsFunc:           psdcc.persistentVolumeClaimLister.PersistentVolumeClaims(psc.Namespace).List,
			PatchObjectFunc:           psdcc.kubeClient.CoreV1().PersistentVolumeClaims(psc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get persistentvolumeclaims: %w", err))
	}

	scyllaDBDatacenterClaimMap, err := socontrollerhelpers.GetObjects[CT, *pausingv1alpha1.ScyllaDBDatacenterClaim](
		ctx,
		psc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *pausingv1alpha1.ScyllaDBDatacenterClaim]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psc.Namespace).Get,
			ListObjectsFunc:           psdcc.scyllaDBDatacenterClaimLister.ScyllaDBDatacenterClaims(psc.Namespace).List,
			PatchObjectFunc:           psdcc.pausingClient.ScyllaDBDatacenterClaims(psc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get ScyllaDBDatacenterClaims: %w", err))
	}

	secretMap, err := socontrollerhelpers.GetObjects[CT, *corev1.Secret](
		ctx,
		psc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *corev1.Secret]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psc.Namespace).Get,
			ListObjectsFunc:           psdcc.secretLister.Secrets(psc.Namespace).List,
			PatchObjectFunc:           psdcc.kubeClient.CoreV1().Secrets(psc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get secrets: %w", err))
	}

	configMapMap, err := socontrollerhelpers.GetObjects[CT, *corev1.ConfigMap](
		ctx,
		psc,
		pausableScyllaDBDatacenterControllerGVK,
		selector,
		socontrollerhelpers.ControlleeManagerGetObjectsFuncs[CT, *corev1.ConfigMap]{
			GetControllerUncachedFunc: psdcc.pausingClient.PausableScyllaDBDatacenters(psc.Namespace).Get,
			ListObjectsFunc:           psdcc.configMapLister.ConfigMaps(psc.Namespace).List,
			PatchObjectFunc:           psdcc.kubeClient.CoreV1().ConfigMaps(psc.Namespace).Patch,
		},
	)
	if err != nil {
		objectErrs = append(objectErrs, fmt.Errorf("can't get configmaps: %w", err))
	}

	objectErr := utilerrors.NewAggregate(objectErrs)
	if objectErr != nil {
		return objectErr
	}

	status := psdcc.calculateStatus(psc, scyllaDBDatacenterClaimMap)

	if psc.DeletionTimestamp != nil {
		return psdcc.updateStatus(ctx, psc, status)
	}

	var errs []error

	err = socontrollerhelpers.RunSync(
		&status.Conditions,
		persistentVolumeClaimControllerProgressingCondition,
		persistentVolumeClaimControllerDegradedCondition,
		psc.Generation,
		func() ([]metav1.Condition, error) {
			return psdcc.syncPersistentVolumeClaims(ctx, psc, persistentVolumeClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync PersistentVolumeClaims: %w", err))
	}

	// TODO: move claim logic here and remove scc controller?
	err = socontrollerhelpers.RunSync(
		&status.Conditions,
		scyllaDBDatacenterClaimControllerProgressingCondition,
		scyllaDBDatacenterClaimControllerDegradedCondition,
		psc.Generation,
		func() ([]metav1.Condition, error) {
			return psdcc.syncScyllaDBDatacenterClaims(ctx, psc, scyllaDBDatacenterClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync ScyllaDBDatacenterClaims: %w", err))
	}

	err = socontrollerhelpers.RunSync(
		&status.Conditions,
		delayedVolumeMountControllerProgressingCondition,
		delayedVolumeMountControllerDegradedCondition,
		psc.Generation,
		func() ([]metav1.Condition, error) {
			return psdcc.syncDelayedVolumeMount(ctx, psc, scyllaDBDatacenterClaimMap)
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("can't sync delayed volume mounting: %w", err))
	}

	err = socontrollerhelpers.RunSync(
		&status.Conditions,
		certControllerProgressingCondition,
		certControllerDegradedCondition,
		psc.Generation,
		func() ([]metav1.Condition, error) {
			return psdcc.syncCerts(ctx, psc, secretMap, configMapMap)
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
			ObservedGeneration: psc.Generation,
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
			ObservedGeneration: psc.Generation,
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
			ObservedGeneration: psc.Generation,
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

	err = psdcc.updateStatus(ctx, psc, status)
	if err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}
