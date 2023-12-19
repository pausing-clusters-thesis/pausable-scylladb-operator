package scylladbdatacenterclaim

import (
	"context"
	"fmt"
	"sort"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (sdccc *Controller) syncScyllaDBDatacenter(ctx context.Context, sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter) ([]metav1.Condition, []metav1.Condition, error) {
	var availableConditions, progressingConditions []metav1.Condition
	var err error

	if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
		// No ScyllaDBDatacenter bound, nothing to do.
		progressingConditions = append(availableConditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: sdcc.Generation,
			Reason:             "ScyllaDBDatacenterNotBound",
			Message:            "Waiting for a ScyllaDBDatacenter to be bound.",
		})
		return availableConditions, progressingConditions, nil
	}

	sdcName := *sdcc.Status.ScyllaDBDatacenterName
	sdc, ok := scyllaDBDatacenters[sdcName]
	if !ok {
		return availableConditions, progressingConditions, fmt.Errorf("can't get bound ScyllaDBDatacenter %q", sonaming.ManualRef(sdcc.Namespace, sdcName))
	}

	sdcAvailable := controllerhelpers.IsScyllaDBDatacenterAvailable(sdc)
	if !sdcAvailable {
		availableConditions = append(availableConditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerAvailableCondition,
			Status:             metav1.ConditionFalse,
			ObservedGeneration: sdcc.Generation,
			Reason:             "ScyllaDBDatacenterNotAvailable",
			Message:            fmt.Sprintf("Waiting for ScyllaDBDatacenter %q to be available.", sonaming.ObjRef(sdc)),
		})
	}

	sdcCopy := sdc.DeepCopy()
	dnsDomains := sets.New(sdcCopy.Spec.DNSDomains...).Insert(sdcc.Spec.DNSDomains...).UnsortedList()
	// Sort to make the manifest consistent across iterations.
	sort.Strings(dnsDomains)

	sdcCopy.Spec.DNSDomains = dnsDomains

	if sdcCopy.Spec.ExposeOptions == nil {
		sdcCopy.Spec.ExposeOptions = &scyllav1alpha1.ExposeOptions{}
	}
	sdcCopy.Spec.ExposeOptions.CQL = sdcc.Spec.ExposeOptions.CQL.DeepCopy()

	sdcCopy.Spec.CertificateOptions = &scyllav1alpha1.CertificateOptions{
		ServingCA: &scyllav1alpha1.TLSCertificateAuthority{
			Type: scyllav1alpha1.TLSCertificateAuthorityTypeUserManaged,
			UserManagedOptions: &scyllav1alpha1.UserManagedTLSCertificateAuthorityOptions{
				SecretName: sdcc.Spec.CertificateOptions.ServingCASecretName,
			},
		},
		ClientCA: &scyllav1alpha1.TLSCertificateAuthority{
			Type: scyllav1alpha1.TLSCertificateAuthorityTypeUserManaged,
			UserManagedOptions: &scyllav1alpha1.UserManagedTLSCertificateAuthorityOptions{
				SecretName: sdcc.Spec.CertificateOptions.ClientCASecretName,
			},
		},
	}

	if !equality.Semantic.DeepEqual(sdc, sdcCopy) {
		_, err = sdccc.scyllaClient.ScyllaDBDatacenters(sdc.Namespace).Update(ctx, sdcCopy, metav1.UpdateOptions{})
		socontrollerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, scyllaDBDatacenterControllerProgressingCondition, sdcCopy, "update", sdcc.Generation)
		if err != nil {
			return availableConditions, progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenter %q: %w", sonaming.ObjRef(sdcCopy), err)
		}
	}

	return availableConditions, progressingConditions, nil
}
