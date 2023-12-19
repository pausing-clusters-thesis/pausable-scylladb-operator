package scylladbdatacenterclaim

import (
	"context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func (sdccc *Controller) calculateStatus(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim, scyllaDBDatacenters map[string]*scyllav1alpha1.ScyllaDBDatacenter) *pausingv1alpha1.ScyllaDBDatacenterClaimStatus {
	status := sdcc.Status.DeepCopy()
	status.ObservedGeneration = ptr.To(sdcc.Generation)

	return status
}

func (sdccc *Controller) updateStatus(ctx context.Context, currentSDCC *pausingv1alpha1.ScyllaDBDatacenterClaim, status *pausingv1alpha1.ScyllaDBDatacenterClaimStatus) error {
	if apiequality.Semantic.DeepEqual(&currentSDCC.Status, status) {
		return nil
	}

	sdcc := currentSDCC.DeepCopy()
	sdcc.Status = *status

	klog.V(2).InfoS("Updating status", "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

	_, err := sdccc.pausableScyllaDBClient.ScyllaDBDatacenterClaims(sdcc.Namespace).UpdateStatus(ctx, sdcc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	klog.V(2).InfoS("Status updated", "ScyllaDBDatacenterClaim", klog.KObj(sdcc))

	return nil
}
