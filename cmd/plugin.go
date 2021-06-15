package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var rootCmd = &cobra.Command{
	Use:          "druid",
	Long:         "kubectl druid plugin",
	SilenceUsage: true,
}

func NewCmdDruidPlugin(streams genericclioptions.IOStreams) *cobra.Command {
	rootCmd.AddCommand(druidCRList(streams))
	rootCmd.AddCommand(druidCRGet(streams))
	rootCmd.AddCommand(druidCRWriterNodeSpecReplicas(streams))
	rootCmd.AddCommand(druidCRWriterUpdates(streams))
	rootCmd.AddCommand(druidCRPatcher(streams))
	return rootCmd
}
