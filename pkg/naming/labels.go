package naming

import (
	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetPausableScyllaDBLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       AppName,
		"app.kubernetes.io/managed-by": OperatorAppName,
	}
}

func ScyllaDBDatacenterClaimSelectorLabels(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim) map[string]string {
	return map[string]string{
		ScyllaDBDatacenterClaimNameLabel: sdcc.Name,
	}
}

func ScyllaDBDatacenterClaimSelector(sdcc *pausingv1alpha1.ScyllaDBDatacenterClaim) labels.Selector {
	return labels.SelectorFromSet(ScyllaDBDatacenterClaimSelectorLabels(sdcc))
}

func ScyllaDBDatacenterPoolSelectorLabels(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) map[string]string {
	return map[string]string{
		ScyllaDBDatacenterPoolNameLabel: sdcp.Name,
	}
}

func ScyllaDBDatacenterPoolSelectorLabelKeys() []string {
	return []string{
		ScyllaDBDatacenterPoolNameLabel,
	}
}

func ScyllaDBDatacenterPoolSelector(sdcp *pausingv1alpha1.ScyllaDBDatacenterPool) labels.Selector {
	return labels.SelectorFromSet(ScyllaDBDatacenterPoolSelectorLabels(sdcp))
}
