package scylladbdatacenterclaim

import (
	"context"
	"fmt"
	"sort"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	"github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (sdccc *Controller) syncScyllaDBDatacenter(ctx context.Context, sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter) ([]metav1.Condition, error) {
	var progressingConditions []metav1.Condition
	var err error

	if sdcc.Status.ScyllaDBDatacenterName == nil || len(*sdcc.Status.ScyllaDBDatacenterName) == 0 {
		// No ScyllaDBDatacenter bound, nothing to do.
		progressingConditions = append(progressingConditions, metav1.Condition{
			Type:               scyllaDBDatacenterControllerProgressingCondition,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: sdcc.Generation,
			Reason:             "AwaitingScyllaDBDatacenterBinding",
			Message:            "Waiting for a ScyllaDBDatacenter to be bound.",
		})
		return progressingConditions, nil
	}

	sdcName := *sdcc.Status.ScyllaDBDatacenterName
	sdc, ok := scyllaDBDatacenters[sdcName]
	if !ok {
		return progressingConditions, fmt.Errorf("can't get bound ScyllaDBDatacenter %q", naming.ManualRef(sdcc.Namespace, sdcName))
	}

	sdcCopy := sdc.DeepCopy()
	dnsDomains := sets.New(sdcCopy.Spec.DNSDomains...).Insert(sdcc.Spec.DNSDomains...).UnsortedList()
	// Sort to make the manifest consistent across iterations.
	sort.Strings(dnsDomains)

	sdcCopy.Spec.DNSDomains = dnsDomains

	if sdcCopy.Spec.ExposeOptions == nil {
		sdcCopy.Spec.ExposeOptions = &scyllav1alpha1.ExposeOptions{}
	}
	sdcCopy.Spec.ExposeOptions.CQL = &scyllav1alpha1.CQLExposeOptions{
		Ingress: &scyllav1alpha1.CQLExposeIngressOptions{
			IngressClassName: naming.IngressClassName,
		},
	}

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

	// TODO: only check for changed fields?
	if !equality.Semantic.DeepEqual(sdc, sdcCopy) {
		_, err = sdccc.scyllaClient.ScyllaDBDatacenters(sdc.Namespace).Update(ctx, sdcCopy, metav1.UpdateOptions{})
		controllerhelpers.AddGenericProgressingStatusCondition(&progressingConditions, scyllaDBDatacenterControllerProgressingCondition, sdcCopy, "update", sdcc.Generation)
		if err != nil {
			return progressingConditions, fmt.Errorf("can't update ScyllaDBDatacenter: %w", err)
		}
	}

	return progressingConditions, nil
}
