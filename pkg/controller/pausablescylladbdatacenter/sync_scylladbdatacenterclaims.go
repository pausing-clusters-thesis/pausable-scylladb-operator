//  Copyright (C) 2024 ScyllaDB

package pausablescylladbdatacenter

import (
	"context"
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/resourceapply"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
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

	return &pausingv1alpha1.ScyllaDBDatacenterClaim{
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
	}, nil
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

func (psdcc *Controller) syncScyllaDBDatacenterClaims(ctx context.Context, psc *pausingv1alpha1.PausableScyllaDBDatacenter, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) ([]metav1.Condition, error) {
	var progressingConditions []metav1.Condition
	var err error
	var required *pausingv1alpha1.ScyllaDBDatacenterClaim

	required, err = makeScyllaDBDatacenterClaim(psc)
	if err != nil {
		return progressingConditions, fmt.Errorf("can't make scyllaDBDatacenterClaim for pausablescylladbcluster: %w", err)
	}

	err = psdcc.pruneScyllaDBDatacenterClaims(ctx, required, scyllaDBDatacenterClaims)
	if err != nil {
		return progressingConditions, fmt.Errorf("can't prune scyllaDBDatacenterClaims: %w", err)
	}

	if required == nil {
		return progressingConditions, nil
	}

	klog.V(4).InfoS("Applying ScyllaDBDatacenterClaim", "PausableScyllaDBDatacenter", klog.KObj(psc), "ScyllaDBDatacenterClaim", klog.KObj(required))
	_, changed, err := resourceapply.ApplyScyllaDBDatacenterClaim(ctx, psdcc.pausingClient, psdcc.scyllaDBDatacenterClaimLister, psdcc.eventRecorder, required, soresourceapply.ApplyOptions{})
	if changed {
		socontrollerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, scyllaDBDatacenterClaimControllerProgressingCondition, required, "apply", psc.Generation)
	}
	if err != nil {
		return progressingConditions, fmt.Errorf("can't apply scyllaDBDatacenterClaim: %w", err)
	}

	return progressingConditions, nil
}
