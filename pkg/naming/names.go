package naming

import (
	"fmt"

	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ManualRef(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func ObjRef(obj metav1.Object) string {
	namespace := obj.GetNamespace()
	if len(namespace) == 0 {
		return obj.GetName()
	}

	return ManualRef(obj.GetNamespace(), obj.GetName())
}

func GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc *pausingv1alpha1.PausableScyllaDBDatacenter) string {
	return fmt.Sprintf("claim-%s", psdc.Name)
}

func GetBackendPersistentVolumeClaimNameForPausableScyllaDBDatacenterMember(pausableScyllaDBDatacenterName string, rackName string, memberIdx int32) string {
	return fmt.Sprintf("%s-%s-%s-%d", BackendPVCTemplateName, pausableScyllaDBDatacenterName, rackName, memberIdx)
}

func GetPausableScyllaDBDatacenterLocalClientCAName(pscName string) string {
	return fmt.Sprintf("%s-pausable-local-client-ca", pscName)
}

func GetPausableScyllaDBDatacenterLocalServingCAName(pscName string) string {
	return fmt.Sprintf("%s-pausable-local-serving-ca", pscName)
}
