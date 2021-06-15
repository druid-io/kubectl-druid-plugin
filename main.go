package main

import (
	"os"

	"github.com/druid-io/druid-kubectl-plugin/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
)

func main() {

	if err := cmd.NewCmdDruidPlugin(genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}).Execute(); err != nil {
		klog.Flush()
		os.Exit(1)
	}

}
