package controllerhelpers

import (
	"fmt"

	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	sohelpers "github.com/scylladb/scylla-operator/pkg/helpers"
	soslices "github.com/scylladb/scylla-operator/pkg/helpers/slices"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
)

func GetRackStorageOptions(sdc *scyllav1alpha1.ScyllaDBDatacenter, rackName string) (*scyllav1alpha1.StorageOptions, error) {
	rackSpec, _, ok := soslices.Find(sdc.Spec.Racks, func(spec scyllav1alpha1.RackSpec) bool {
		return spec.Name == rackName
	})
	if !ok {
		return nil, fmt.Errorf("can't find rack %q in rack spec of ScyllaDBDatacenter %q", rackName, sonaming.ObjRef(sdc))
	}

	if rackSpec.ScyllaDB != nil && rackSpec.ScyllaDB.Storage != nil {
		return rackSpec.ScyllaDB.Storage, nil
	}

	if sdc.Spec.RackTemplate != nil && sdc.Spec.RackTemplate.ScyllaDB != nil && sdc.Spec.RackTemplate.ScyllaDB.Storage != nil {
		return sdc.Spec.RackTemplate.ScyllaDB.Storage, nil
	}

	return nil, fmt.Errorf("can't get storage options for rack %q of ScyllaDBDatacenter %q", rackName, sonaming.ObjRef(sdc))
}

func IsScyllaDBDatacenterPrewarmed(sdc *scyllav1alpha1.ScyllaDBDatacenter) bool {
	return sohelpers.IsStatusConditionPresentAndTrue(sdc.Status.Conditions, scyllav1alpha1.PrewarmedCondition, sdc.Generation)
}

func IsScyllaDBDatacenterAvailable(sdc *scyllav1alpha1.ScyllaDBDatacenter) bool {
	return sohelpers.IsStatusConditionPresentAndTrue(sdc.Status.Conditions, scyllav1alpha1.AvailableCondition, sdc.Generation)
}
