package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func druidCRList(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidIoWriter{
		out: streams.Out,
	}

	var namespace, CR string
	cmd := &cobra.Command{
		Use:          "list",
		Short:        "Lists druid CR's in all namespaces or a specific namespce",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return druidCmdList.druidCRListRun(namespace, CR, args)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&namespace, "namespace", "n", "", "namespace of druid CR")
	return cmd
}

func (sv *druidIoWriter) druidCRListRun(namespace, CR string, args []string) error {

	listDruids, err := di.listDruidCR(namespace)
	if err != nil {
		return err
	}

	// Format in tab-separated columns with a tab stop of 8.
	sv.w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintf(&sv.w, "NAME\tNAMESPACE\t\n")

	for namespace, crNames := range listDruids {
		for _, crName := range crNames {
			_, err := fmt.Fprintf(&sv.w, "%s\t%s\t\n", namespace, crName)
			if err != nil {
				return err
			}
		}
	}

	sv.w.Flush()

	return nil
}

func druidCRGet(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidIoWriter{
		out: streams.Out,
	}

	var namespace, cr string
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "get druid CR specfic info. Ex: kubectl druid get nodes --cr <cr> --namespace <namespace>",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := druidCmdList.validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			return druidCmdList.druidCRGetRun(namespace, cr)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&namespace, "namespace", "n", "", "namespace of druid CR")
	f.StringVar(&cr, "cr", "", "name of druid CR")

	return cmd
}

func (sv *druidIoWriter) druidCRGetRun(namespace, cr string) error {

	getDruidNodes, err := di.getDruidNodeNames(namespace, cr)
	if err != nil {
		return err
	}

	// Format in tab-separated columns with a tab stop of 8.
	sv.w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintf(&sv.w, "NAME\tNAMESPACE\t\n")

	for namespace, nodeNames := range getDruidNodes {
		for _, nodeName := range sort.StringSlice(nodeNames) {
			_, err := fmt.Fprintf(&sv.w, "%s\t%s\t\n", nodeName, namespace)
			if err != nil {
				return err
			}
		}
	}

	sv.w.Flush()

	return nil
}

func (sv *druidIoWriter) validate(args []string) error {
	if args == nil || args[0] != "nodes" {
		return errors.New("invalid arg, valid arg is [nodes] e.g. 'kubectl druid get nodes --cr <cr> --namespace <namespace>'")
	}
	return nil
}
