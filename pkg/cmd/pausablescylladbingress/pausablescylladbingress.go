package pausablescylladbingress

import (
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/cmd/pausablescylladbingress/run"
	"github.com/scylladb/scylla-operator/pkg/cmdutil"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func NewCommand(streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{}

	cmd.AddCommand(run.NewCommand(streams))

	cmdutil.InstallKlog(cmd)

	return cmd
}
