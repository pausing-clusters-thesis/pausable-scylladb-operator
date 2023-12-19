package pausablescylladbdatacenter

const (
	persistentVolumeClaimControllerAvailableCondition     = "PersistentVolumeClaimControllerAvailable"
	persistentVolumeClaimControllerProgressingCondition   = "PersistentVolumeClaimControllerProgressing"
	persistentVolumeClaimControllerDegradedCondition      = "PersistentVolumeClaimControllerDegraded"
	scyllaDBDatacenterClaimControllerAvailableCondition   = "ScyllaDBDatacenterClaimControllerAvailable"
	scyllaDBDatacenterClaimControllerProgressingCondition = "ScyllaDBDatacenterClaimControllerProgressing"
	scyllaDBDatacenterClaimControllerDegradedCondition    = "ScyllaDBDatacenterClaimControllerDegraded"
	delayedVolumeMountControllerAvailableCondition        = "DelayedVolumeMountControllerAvailable"
	delayedVolumeMountControllerProgressingCondition      = "DelayedVolumeMountControllerProgressing"
	delayedVolumeMountControllerDegradedCondition         = "DelayedVolumeMountControllerDegraded"
	certControllerAvailableCondition                      = "CertControllerAvailable"
	certControllerProgressingCondition                    = "CertControllerProgressing"
	certControllerDegradedCondition                       = "CertControllerDegraded"
)
