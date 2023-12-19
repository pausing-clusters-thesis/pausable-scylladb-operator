package pausablescylladbdatacenter

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/resourceapply"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	soresourceapply "github.com/scylladb/scylla-operator/pkg/resourceapply"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apierrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func makeScyllaDBDatacenterClaim(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) (*pausingv1alpha1.ScyllaDBDatacenterClaim, error) {
	if psdc.Spec.Paused != nil && *psdc.Spec.Paused {
		return nil, nil
	}

	sdcc := &pausingv1alpha1.ScyllaDBDatacenterClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      naming.GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc),
			Namespace: psdc.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(psdc, pausableScyllaDBDatacenterControllerGVK),
			},
			Labels: getSelectorLabels(psdc),
			// TODO: annotations
			// TODO: finalizers
		},
		Spec: pausingv1alpha1.ScyllaDBDatacenterClaimSpec{
			ScyllaDBDatacenterPoolName: psdc.Spec.ScyllaDBDatacenterPoolName,
			DNSDomains:                 []string{fmt.Sprintf(naming.PausableScyllaDBClusterDNSDomainFormat, psdc.Name, psdc.Namespace)},
			CertificateOptions: &pausingv1alpha1.CertificateOptions{
				ServingCASecretName: naming.GetPausableScyllaDBDatacenterLocalServingCAName(psdc.GetName()),
				ClientCASecretName:  naming.GetPausableScyllaDBDatacenterLocalClientCAName(psdc.GetName()),
			},
		},
	}

	if psdc.Spec.ExposeOptions != nil {
		sdcc.Spec.ExposeOptions = psdc.Spec.ExposeOptions.DeepCopy()
	}

	return sdcc, nil
}

func (psdcc *Controller) pruneScyllaDBDatacenterClaims(ctx context.Context, required *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) error {
	var errs []error

	for _, sdc := range scyllaDBDatacenterClaims {
		if sdc.DeletionTimestamp != nil {
			continue
		}

		if required != nil && sdc.Name == required.Name {
			continue
		}

		err := psdcc.pausingClient.ScyllaDBDatacenterClaims(sdc.Namespace).Delete(ctx, sdc.Name, metav1.DeleteOptions{
			Preconditions: &metav1.Preconditions{
				UID: &sdc.UID,
			},
			PropagationPolicy: ptr.To(metav1.DeletePropagationBackground),
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	return apierrors.NewAggregate(errs)
}

func (psdcc *Controller) syncScyllaDBDatacenterClaims(ctx context.Context, psdc *pausingv1alpha1.PausableScyllaDBDatacenter, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) ([]metav1.Condition, []metav1.Condition, error) {
	var availableConditions, progressingConditions []metav1.Condition
	var err error
	var sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim

	sdcc, err = makeScyllaDBDatacenterClaim(psdc)
	if err != nil {
		return availableConditions, progressingConditions, fmt.Errorf("can't make scyllaDBDatacenterClaim for pausablescylladbcluster: %w", err)
	}

	err = psdcc.pruneScyllaDBDatacenterClaims(ctx, sdcc, scyllaDBDatacenterClaims)
	if err != nil {
		return availableConditions, progressingConditions, fmt.Errorf("can't prune scyllaDBDatacenterClaims: %w", err)
	}

	if sdcc == nil {
		return availableConditions, progressingConditions, nil
	}

	klog.V(4).InfoS("Applying ScyllaDBDatacenterClaim", "PausableScyllaDBDatacenter", klog.KObj(psdc), "ScyllaDBDatacenterClaim", klog.KObj(sdcc))
	sdcc, changed, err := resourceapply.ApplyScyllaDBDatacenterClaim(ctx, psdcc.pausingClient, psdcc.scyllaDBDatacenterClaimLister, psdcc.eventRecorder, sdcc, soresourceapply.ApplyOptions{})
	if changed {
		socontrollerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, scyllaDBDatacenterClaimControllerProgressingCondition, sdcc, "apply", psdc.Generation)
	}
	if err != nil {
		return availableConditions, progressingConditions, fmt.Errorf("can't apply scyllaDBDatacenterClaim: %w", err)
	}

	isScyllaDBDatacenterClaimRolledOut, err := controllerhelpers.IsScyllaDBDatacenterClaimRolledOut(sdcc)
	if err != nil {
		return availableConditions, progressingConditions, fmt.Errorf("can't get ScyllaDBDatacenterClaim %q rollout status: %w", sonaming.ObjRef(sdcc), err)
	}
	if !isScyllaDBDatacenterClaimRolledOut {
		availableConditions = append(availableConditions, metav1.Condition{
			Type:               scyllaDBDatacenterClaimControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: psdc.Generation,
			Reason:             "ScyllaDBDatacenterClaimNotRolledOut",
			Message:            fmt.Sprintf("ScyllaDBDatacenterClaim %q is not rolled out.", naming.ObjRef(sdcc)),
		})
	}

	return availableConditions, progressingConditions, nil
}
