package main

import (
	"flag"
	"os"

	"github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/cmd/pausablescylladboperator"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/component-base/cli"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(flag.CommandLine)
	command := pausablescylladboperator.NewCommand(genericiooptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})
	code := cli.Run(command)
	os.Exit(code)
}
