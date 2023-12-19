package pausablescylladbdatacenter

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	proxycsinaming "github.com/pausing-clusters-thesis/proxy-csi-driver/pkg/naming"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/internalapi"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

func (psdcc *Controller) syncDelayedVolumeMount(
	ctx context.Context,
	psdc *pausingv1alpha1.PausableScyllaDBDatacenter,
	scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim,
) ([]metav1.Condition, []metav1.Condition, error) {
	var availableConditions, progressingConditions []metav1.Condition

	if psdc.Spec.Paused == nil || *psdc.Spec.Paused {
		return availableConditions, progressingConditions, nil
	}

	sdcc, ok := scyllaDBDatacenterClaims[naming.GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc)]
	if !ok {
		return availableConditions, progressingConditions, nil
	}

	if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
		return availableConditions, progressingConditions, nil
	}

	sdc, err := psdcc.scyllaDBDatacenterLister.ScyllaDBDatacenters(psdc.Namespace).Get(*sdcc.Status.ScyllaDBDatacenterName)
	if err != nil {
		return availableConditions, progressingConditions, fmt.Errorf("can't get ScyllaDBDatacenter %q: %w", naming.ManualRef(psdc.Namespace, *sdcc.Status.ScyllaDBDatacenterName), err)
	}

	for _, r := range sdc.Spec.Racks {
		rackNodeCount, err := socontrollerhelpers.GetRackNodeCount(sdc, r.Name)
		if err != nil {
			return availableConditions, progressingConditions, fmt.Errorf("can't get rack node count for rack %q: %w", r.Name, err)
		}

		for i := range *rackNodeCount {
			proxyPVCName := sonaming.PVCNameForStatefulSet(sonaming.StatefulSetNameForRack(r, sdc), i)
			proxyPVC, err := psdcc.persistentVolumeClaimLister.PersistentVolumeClaims(psdc.Namespace).Get(proxyPVCName)
			if err != nil {
				return availableConditions, progressingConditions, fmt.Errorf("can't get persistent volume claim %q: %w", naming.ManualRef(psdc.Namespace, proxyPVCName), err)
			}

			backendPVCName := naming.GetBackendPersistentVolumeClaimNameForPausableScyllaDBDatacenterMember(psdc.Name, r.Name, i)
			if controllerhelpers.HasMatchingAnnotation(proxyPVC, proxycsinaming.DelayedStorageBackendPersistentVolumeClaimRefAnnotation, backendPVCName) {
				continue
			}

			err = psdcc.setAnnotation(ctx, proxyPVC, proxycsinaming.DelayedStorageBackendPersistentVolumeClaimRefAnnotation, backendPVCName)
			progressingConditions = append(progressingConditions, metav1.Condition{
				Type:               delayedVolumeMountControllerProgressingCondition,
				Status:             metav1.ConditionTrue,
				ObservedGeneration: psdc.Generation,
				Reason:             internalapi.ProgressingReason,
				Message:            fmt.Sprintf("Annotating proxy PersistentVolumeClaim %q with reference annotation to backend PersistentVolumeClaim %q.", naming.ObjRef(proxyPVC), backendPVCName),
			})
			if err != nil {
				return availableConditions, progressingConditions, fmt.Errorf("can't set annotation: %w", err)
			}
		}
	}

	return availableConditions, progressingConditions, nil
}

func (psdcc *Controller) setAnnotation(ctx context.Context, pvc *corev1.PersistentVolumeClaim, annotationKey, annotationValue string) error {
	patch, err := controllerhelpers.PrepareSetAnnotationPatch(pvc, annotationKey, ptr.To(annotationValue))
	if err != nil {
		return fmt.Errorf("can't prepare patch setting annotation: %w", err)
	}

	_, err = psdcc.kubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Patch(ctx, pvc.Name, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("can't patch PersistentVolumeClaim %q: %w", naming.ObjRef(pvc), err)
	}

	return nil
}
