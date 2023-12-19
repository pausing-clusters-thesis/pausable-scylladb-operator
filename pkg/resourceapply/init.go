package resourceapply

import (
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/scheme"
	"github.com/scylladb/scylla-operator/pkg/resource"
)

func init() {
	resource.Scheme = scheme.Scheme
}
