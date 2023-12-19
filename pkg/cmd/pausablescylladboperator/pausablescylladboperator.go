package pausablescylladboperator

import (
	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/cmd/pausablescylladboperator/run"
	"github.com/scylladb/scylla-operator/pkg/cmdutil"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func NewCommand(streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdutil.ReadFlagsFromEnv("PAUSABLE_SCYLLADB_OPERATOR_", cmd)
		},
	}

	cmd.AddCommand(run.NewCommand(streams))

	cmdutil.InstallKlog(cmd)

	return cmd
}
