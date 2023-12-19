package naming

import (
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
)

const (
	AppName         = "pausable-scylladb-operator.scylladb.com/scylladb"
	OperatorAppName = "pausable-scylladb-operator.scylladb.com/pausable-scylladb-operator"
)

const (
	PausableScyllaDBDatacenterNameLabel = "pausable-scylladb-operator.scylladb.com/pausable-scylladb-datacenter-name"
)

const (
	PausableScyllaDBDatacenterDNSDomainSuffix = "pausablescylladbdatacenter.pausingv1alpha1.scylladb.com"
	PausableScyllaDBClusterDNSDomainFormat    = "%s.%s." + PausableScyllaDBDatacenterDNSDomainSuffix
)

const (
	IngressControllerName = "pausable-scylladb-operator.scylladb.com/ingress-controller"
	IngressClassName      = "pausable-scylladb-ingress"
)

const (
	BackendPVCTemplateName = "backend"
)

const (
	ScyllaDBDatacenterClaimNameLabel = sonaming.IgnoreInternalKeyPrefix + "/scylladb-datacenter-claim"
)

const (
	ScyllaDBDatacenterClaimBindCompletedAnnotation = sonaming.IgnoreInternalKeyPrefix + "/scylladb-datacenter-bind-completed"
)

const (
	ScyllaDBDatacenterPoolNameLabel = sonaming.IgnoreInternalKeyPrefix + "/scylladb-datacenter-pool-name"
)

const (
	IngressControllerScyllaDBMemberPodConditionType = "pausable-scylladb-operator.scylladb.com/member-ingress-ready"
)
