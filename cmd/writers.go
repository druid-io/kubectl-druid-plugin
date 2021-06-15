package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type druidWritererCmd struct {
	out io.Writer
}

func druidCRWriterNodeSpecReplicas(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidListCmd{
		out: streams.Out,
	}

	var node, namespace, cr string
	var replicas int64
	cmd := &cobra.Command{
		Use:          "scale",
		Short:        "scale druid node replica counts",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return druidCmdList.druidCRWriterNodeSpecReplicasRun(node, namespace, cr, replicas)
		},
	}

	f := cmd.Flags()
	f.StringVar(&namespace, "namespace", "", "namespace of druid CR")
	f.StringVar(&node, "node", "", "name of druid node created by the druid operator, can be a statefulset or deployment")
	f.StringVar(&cr, "cr", "", "name of the druid CR")
	f.Int64Var(&replicas, "replicas", replicas, "number of replicas to scale")

	return cmd
}

func (sv *druidListCmd) druidCRWriterNodeSpecReplicasRun(nodeName, namespace, CR string, replica int64) error {

	writerResult, err := di.writerDruidNodeSpecReplicas(nodeName, namespace, CR, replica)
	if err != nil {
		return err
	}

	if writerResult {
		_, err := fmt.Fprintf(sv.out, "Druid CR [%s],NodeName [%s] successfully updated in Namespace [%s] with Replica Count [%d]\n", CR, nodeName, namespace, replica)
		if err != nil {
			return err
		}
	}

	return nil
}

func druidCRWriterUpdates(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidListCmd{
		out: streams.Out,
	}

	var nodeName, namespace, CR, image string
	cmd := &cobra.Command{
		Use:          "update",
		Short:        "update druid node images",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return druidCmdList.druidCRWriterUpdatesRun(nodeName, namespace, CR, image)
		},
	}

	f := cmd.Flags()
	f.StringVar(&namespace, "namespace", "", "namespace of druid CR")

	f.StringVar(&nodeName, "node", "", "name of druid node created by the druid operator, can be a statefulset or deployment")

	f.StringVar(&CR, "cr", "", "name of the druid CR")
	f.StringVar(&image, "image", "", "image of the druid node")

	return cmd
}

func (sv *druidListCmd) druidCRWriterUpdatesRun(nodeName, namespace, CR, image string) error {

	writerResult, err := di.writerDruidNodeImages(nodeName, namespace, CR, image)
	if err != nil {
		return err
	}

	if writerResult {
		_, err := fmt.Fprintf(sv.out, "Druid CR [%s],NodeName [%s] successfully update with image [%s] in Namespace [%s]\n", CR, nodeName, image, namespace)
		if err != nil {
			return err
		}
	}

	return nil
}
