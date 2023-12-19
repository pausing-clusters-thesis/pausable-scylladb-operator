package pausablescylladbdatacenter

import (
	"context"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func (psdcc *Controller) calculateStatus(psdc *pausingv1alpha1.PausableScyllaDBDatacenter, persistentVolumeClaims map[string]*corev1.PersistentVolumeClaim, scyllaDBDatacenterClaims map[string]*pausingv1alpha1.ScyllaDBDatacenterClaim) *pausingv1alpha1.PausableScyllaDBDatacenterStatus {
	status := psdc.Status.DeepCopy()
	status.ObservedGeneration = ptr.To(psdc.Generation)

	return status
}

func (psdcc *Controller) updateStatus(ctx context.Context, curPSDC *pausingv1alpha1.PausableScyllaDBDatacenter, status *pausingv1alpha1.PausableScyllaDBDatacenterStatus) error {
	if apiequality.Semantic.DeepEqual(&curPSDC.Status, status) {
		return nil
	}

	scc := curPSDC.DeepCopy()
	scc.Status = *status

	klog.V(2).InfoS("Updating status", "PausableScyllaDBDatacenter", klog.KObj(scc))

	_, err := psdcc.pausingClient.PausableScyllaDBDatacenters(scc.Namespace).UpdateStatus(ctx, scc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	klog.V(2).InfoS("Status updated", "PausableScyllaDBDatacenter", klog.KObj(scc))

	return nil
}
