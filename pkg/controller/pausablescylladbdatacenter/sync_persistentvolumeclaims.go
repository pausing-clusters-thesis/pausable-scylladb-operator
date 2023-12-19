//  Copyright (C) 2024 ScyllaDB

package pausablescylladbdatacenter

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	soresourceapply "github.com/scylladb/scylla-operator/pkg/resourceapply"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func (psdcc *Controller) makePersistentVolumeClaims(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) ([]*corev1.PersistentVolumeClaim, error) {
	var pvcs []*corev1.PersistentVolumeClaim
	var errs []error

	sdcp, err := psdcc.scyllaDBDatacenterPoolLister.ScyllaDBDatacenterPools(psdc.Namespace).Get(psdc.Spec.ScyllaDBDatacenterPoolName)
	if err != nil {
		return nil, fmt.Errorf("can't get ScyllaDBDatacenterPool %q: %w", naming.ManualRef(psdc.Namespace, psdc.Spec.ScyllaDBDatacenterPoolName), err)
	}

	for _, r := range sdcp.Spec.Template.Spec.Racks {
		// Use a dummy SDC to use SO helpers.
		sdc := &scyllav1alpha1.ScyllaDBDatacenter{
			Spec: sdcp.Spec.Template.Spec,
		}

		storageOptions, err := controllerhelpers.GetRackStorageOptions(sdc, r.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't get storage options for rack %q: %w", r.Name, err))
			continue
		}

		storageCapacity, err := resource.ParseQuantity(storageOptions.Capacity)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't parse storage capacity: %w", err))
			continue
		}

		rackNodeCount, err := socontrollerhelpers.GetRackNodeCount(sdc, r.Name)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't get rack node count for rack %q: %w", r.Name, err))
			continue
		}

		for i := range *rackNodeCount {
			pvc := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      naming.GetBackendPersistentVolumeClaimNameForPausableScyllaDBDatacenterMember(psdc.Name, r.Name, i),
					Namespace: psdc.Namespace,
					// TODO: better labels
					Labels: getSelectorLabels(psdc),
					// TODO: should this set an owner ref? Kind of dangerous?
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(psdc, pausableScyllaDBDatacenterControllerGVK),
					},
					Finalizers:    nil,
					ManagedFields: nil,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: storageCapacity,
						},
					},
					// TODO: verify storage class is immediate? We could log a warning.
					StorageClassName: storageOptions.StorageClassName,
					// TODO: data source?
					// TODO: data source ref?
					// TODO: volume attributes class name?
				},
			}

			pvcs = append(pvcs, pvc)
		}
	}

	err = utilerrors.NewAggregate(errs)
	if err != nil {
		return nil, err
	}

	return pvcs, nil
}

func (psdcc *Controller) syncPersistentVolumeClaims(ctx context.Context, psc *pausingv1alpha1.PausableScyllaDBDatacenter, persistentVolumeClaims map[string]*corev1.PersistentVolumeClaim) ([]metav1.Condition, error) {
	var progressingConditions []metav1.Condition

	requiredPersistentVolumeClaims, err := psdcc.makePersistentVolumeClaims(psc)
	if err != nil {
		return progressingConditions, fmt.Errorf("can't make PersistentVolumeClaim(s): %w", err)
	}

	// TODO: Pruning?

	var errs []error
	for _, pvc := range requiredPersistentVolumeClaims {
		_, changed, err := soresourceapply.ApplyPersistentVolumeClaim(ctx, psdcc.kubeClient.CoreV1(), psdcc.persistentVolumeClaimLister, psdcc.eventRecorder, pvc, soresourceapply.ApplyOptions{})
		if changed {
			socontrollerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, persistentVolumeClaimControllerProgressingCondition, pvc, "apply", psc.Generation)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("can't apply PersistentVolumeClaim: %w", err))
			continue
		}
	}
	err = utilerrors.NewAggregate(errs)
	if err != nil {
		return progressingConditions, err
	}

	return progressingConditions, nil
}
