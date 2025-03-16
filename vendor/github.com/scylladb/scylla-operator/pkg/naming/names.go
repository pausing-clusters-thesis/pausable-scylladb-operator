package naming

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containers/image/v5/docker/reference"
	"github.com/pkg/errors"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ManualRef(namespace, name string) string {
	if len(namespace) == 0 {
		return name
	}
	return fmt.Sprintf("%s/%s", namespace, name)
}

func ObjRef(obj metav1.Object) string {
	return ManualRef(obj.GetNamespace(), obj.GetName())
}

func ObjRefWithUID(obj metav1.Object) string {
	return fmt.Sprintf("%s(UID=%s)", ObjRef(obj), obj.GetUID())
}

func StatefulSetNameForRack(r scyllav1alpha1.RackSpec, sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return fmt.Sprintf("%s-%s-%s", sdc.Name, GetScyllaDBDatacenterGossipDatacenterName(sdc), r.Name)
}

func StatefulSetNameForRackForScyllaCluster(r scyllav1.RackSpec, sc *scyllav1.ScyllaCluster) string {
	return fmt.Sprintf("%s-%s-%s", sc.Name, sc.Spec.Datacenter.Name, r.Name)
}

func PodNameFromService(svc *corev1.Service) string {
	// Pod and its corresponding Service have the same name
	return svc.Name
}

func AgentAuthTokenSecretName(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return fmt.Sprintf("%s-auth-token", sdc.Name)
}

func AgentAuthTokenSecretNameForScyllaCluster(sc *scyllav1.ScyllaCluster) string {
	return AgentAuthTokenSecretName(&scyllav1alpha1.ScyllaDBDatacenter{
		ObjectMeta: metav1.ObjectMeta{
			Name: sc.Name,
		},
	})
}

func MemberServiceName(r scyllav1alpha1.RackSpec, sdc *scyllav1alpha1.ScyllaDBDatacenter, idx int) string {
	return fmt.Sprintf("%s-%d", StatefulSetNameForRack(r, sdc), idx)
}

func MemberServiceNameForScyllaCluster(r scyllav1.RackSpec, sc *scyllav1.ScyllaCluster, idx int) string {
	return fmt.Sprintf("%s-%d", StatefulSetNameForRackForScyllaCluster(r, sc), idx)
}

func PodNameForScyllaCluster(r scyllav1.RackSpec, sc *scyllav1.ScyllaCluster, idx int) string {
	return MemberServiceNameForScyllaCluster(r, sc, idx)
}

func IdentityServiceName(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return fmt.Sprintf("%s-client", sdc.Name)
}

func IdentityServiceNameForScyllaCluster(sc *scyllav1.ScyllaCluster) string {
	return fmt.Sprintf("%s-client", sc.Name)
}

func PodDisruptionBudgetName(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return sdc.Name
}

func PodDisruptionBudgetNameForScyllaCluster(sc *scyllav1.ScyllaCluster) string {
	return sc.Name
}

func CrossNamespaceServiceNameForCluster(sc *scyllav1.ScyllaCluster) string {
	return fmt.Sprintf("%s.%s.svc", IdentityServiceNameForScyllaCluster(sc), sc.Namespace)
}

func ManagerClusterName(sc *scyllav1.ScyllaCluster) string {
	return sc.Namespace + "/" + sc.Name
}

func PVCNameForPod(podName string) string {
	return fmt.Sprintf("%s-%s", PVCTemplateName, podName)
}

func PVCNameForService(svcName string) string {
	return PVCNameForPod(svcName)
}

func PVCNamePrefix(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return fmt.Sprintf("%s-%s-", PVCTemplateName, sdc.Name)
}

func PVCNamePrefixForScyllaCluster(sc *scyllav1.ScyllaCluster) string {
	return fmt.Sprintf("%s-%s-", PVCTemplateName, sc.Name)
}

func PVCNameForStatefulSet(stsName string, ordinal int32) string {
	return fmt.Sprintf("%s-%s-%d", PVCTemplateName, stsName, ordinal)
}

// IndexFromName attempts to get the index from a name using the
// naming convention <name>-<index>.
func IndexFromName(n string) (int32, error) {

	// index := svc.Name[strings.LastIndex(svc.Name, "-") + 1 : len(svc.Name)]
	delimIndex := strings.LastIndex(n, "-")
	if delimIndex == -1 {
		return -1, errors.New(fmt.Sprintf("didn't find '-' delimiter in string %s", n))
	}

	index, err := strconv.Atoi(n[delimIndex+1:])
	if err != nil {
		return -1, errors.New(fmt.Sprintf("couldn't convert '%s' to a number", n[delimIndex+1:]))
	}

	return int32(index), nil
}

// ImageToVersion strips version part from container image.
func ImageToVersion(image string) (string, error) {
	_, version, err := ImageToRepositoryVersion(image)
	if err != nil {
		return "", err
	}

	return version, nil
}

func ImageToRepositoryVersion(image string) (string, string, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", "", fmt.Errorf("can't parse image: %w", err)
	}
	tagged, ok := named.(reference.NamedTagged)
	if !ok {
		return "", "", fmt.Errorf("invalid, non-tagged image reference of type %T: %s", named, image)
	}

	return tagged.Name(), tagged.Tag(), nil
}

// FindScyllaContainer returns Scylla container from given list.
func FindScyllaContainer(containers []corev1.Container) (int, error) {
	return FindContainerWithName(containers, ScyllaContainerName)
}

// FindSidecarInjectorContainer returns sidecar injector container from given list.
func FindSidecarInjectorContainer(containers []corev1.Container) (int, error) {
	return FindContainerWithName(containers, SidecarInjectorContainerName)
}

// FindContainerWithName returns container having
func FindContainerWithName(containers []corev1.Container, name string) (int, error) {
	for idx := range containers {
		if containers[idx].Name == name {
			return idx, nil
		}
	}
	return 0, errors.Errorf(" '%s' container not found", name)
}

// ScyllaVersion returns version of Scylla container.
func ScyllaVersion(containers []corev1.Container) (string, error) {
	idx, err := FindScyllaContainer(containers)
	if err != nil {
		return "", errors.Wrap(err, "find scylla container")
	}

	version, err := ImageToVersion(containers[idx].Image)
	if err != nil {
		return "", errors.Wrap(err, "parse scylla container version")
	}
	return version, nil
}

func GetTuningConfigMapNameForPod(pod *corev1.Pod) string {
	return fmt.Sprintf("nodeconfig-podinfo-%s", pod.UID)
}

func MemberServiceAccountNameForScyllaDBDatacenter(sdcName string) string {
	return fmt.Sprintf("%s-member", sdcName)
}

func GetScyllaClusterLocalClientCAName(scName string) string {
	return fmt.Sprintf("%s-local-client-ca", scName)
}

func GetScyllaClusterLocalUserAdminCertName(scName string) string {
	return fmt.Sprintf("%s-local-user-admin", scName)
}

func GetScyllaClusterLocalServingCAName(scName string) string {
	return fmt.Sprintf("%s-local-serving-ca", scName)
}

func GetScyllaClusterLocalServingCertName(scName string) string {
	return fmt.Sprintf("%s-local-serving-certs", scName)
}

func GetScyllaClusterLocalAdminCQLConnectionConfigsName(scName string) string {
	return fmt.Sprintf("%s-local-cql-connection-configs-admin", scName)
}

func GetScyllaClusterAlternatorLocalServingCAName(scName string) string {
	return fmt.Sprintf("%s-alternator-local-serving-ca", scName)
}

func GetScyllaClusterAlternatorLocalServingCertName(scName string) string {
	return fmt.Sprintf("%s-alternator-local-serving-certs", scName)
}

func GetProtocolSubDomain(protocol, domain string) string {
	return fmt.Sprintf("%s.%s", protocol, domain)
}

func GetCQLProtocolSubDomain(domain string) string {
	return GetProtocolSubDomain(CQLProtocolDNSLabel, domain)
}

func GetCQLHostIDSubDomain(hostID, domain string) string {
	return fmt.Sprintf("%s.%s", hostID, GetCQLProtocolSubDomain(domain))
}

func CleanupJobForService(svcName string) string {
	return fmt.Sprintf("cleanup-%s", svcName)
}

func GetScyllaDBManagedConfigCMName(clusterName string) string {
	return fmt.Sprintf("%s-managed-config", clusterName)
}

func GetScyllaDBDatacenterGossipDatacenterName(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	if sdc.Spec.DatacenterName != nil && len(*sdc.Spec.DatacenterName) != 0 {
		return *sdc.Spec.DatacenterName
	}
	return sdc.Name
}

func UpgradeContextConfigMapName(sdc *scyllav1alpha1.ScyllaDBDatacenter) string {
	return fmt.Sprintf("%s-upgrade-context", sdc.Name)
}

func DCNameFromSeedServiceAddress(sc *scyllav1alpha1.ScyllaDBCluster, seedServiceAddress, namespace string) string {
	dcName := strings.TrimPrefix(seedServiceAddress, fmt.Sprintf("%s-", sc.Name))
	dcName = strings.TrimSuffix(dcName, fmt.Sprintf("-seed.%s.svc", namespace))
	return dcName
}

func SeedService(sc *scyllav1alpha1.ScyllaDBCluster, dc *scyllav1alpha1.ScyllaDBClusterDatacenter) string {
	return fmt.Sprintf("%s-%s-seed", sc.Name, dc.Name)
}

func ScyllaDBDatacenterName(sc *scyllav1alpha1.ScyllaDBCluster, dc *scyllav1alpha1.ScyllaDBClusterDatacenter) string {
	return fmt.Sprintf("%s-%s", sc.Name, dc.Name)
}
