//  Copyright (C) 2024 ScyllaDB

package pausablescylladbdatacenter

const (
	persistentVolumeClaimControllerProgressingCondition   = "PersistentVolumeClaimControllerProgressing"
	persistentVolumeClaimControllerDegradedCondition      = "PersistentVolumeClaimControllerDegraded"
	scyllaDBDatacenterClaimControllerAvailableCondition   = "ScyllaDBDatacenterClaimControllerAvailable"
	scyllaDBDatacenterClaimControllerProgressingCondition = "ScyllaDBDatacenterClaimControllerProgressing"
	scyllaDBDatacenterClaimControllerDegradedCondition    = "ScyllaDBDatacenterClaimControllerDegraded"
	// TODO: delayed volume mount available condition?
	delayedVolumeMountControllerProgressingCondition = "DelayedVolumeMountControllerProgressing"
	delayedVolumeMountControllerDegradedCondition    = "DelayedVolumeMountControllerDegraded"
	certControllerProgressingCondition               = "CertControllerProgressing"
	certControllerDegradedCondition                  = "CertControllerDegraded"
)
