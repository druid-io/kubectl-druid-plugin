package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type druidListCmd struct {
	out io.Writer
}

func druidCRList(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidListCmd{
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
	f.StringVar(&namespace, "namespace", "", "namespace of druid CR")
	return cmd
}

func (sv *druidListCmd) druidCRListRun(namespace, CR string, args []string) error {

	listDruids, err := di.listDruidCR(namespace)
	if err != nil {
		return err
	}

	for _, l := range listDruids {
		_, err := fmt.Fprintf(sv.out, "%s\n", l)
		if err != nil {
			return err
		}
	}

	return nil
}

func druidCRGet(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidListCmd{
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
	f.StringVar(&namespace, "namespace", "", "namespace of druid CR")
	f.StringVar(&cr, "cr", "", "name of druid CR")

	return cmd
}

func (sv *druidListCmd) druidCRGetRun(namespace, cr string) error {

	getDruidNodes, err := di.getDruidNodeNames(namespace, cr)
	if err != nil {
		return err
	}

	for _, l := range getDruidNodes {
		_, err := fmt.Fprintf(sv.out, "%s\n", l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sv *druidListCmd) validate(args []string) error {
	if args == nil || args[0] != "nodes" {
		return errors.New("invalid arg, valid arg is [nodes] e.g. 'kubectl druid get nodes --cr <cr> --namespace <namespace>'")
	}
	return nil
}
