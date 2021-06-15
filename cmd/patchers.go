package cmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type druidPatcherCmd struct {
	out io.Writer
}

func druidCRPatcher(streams genericclioptions.IOStreams) *cobra.Command {
	druidCmdList := &druidListCmd{
		out: streams.Out,
	}

	var deleteOrphanPvc, rollingDeploy, namespace, cr string
	cmd := &cobra.Command{
		Use:          "patch",
		Short:        "patches druid CR",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			return druidCmdList.druidCRPatcherRun(namespace, cr, deleteOrphanPvc, rollingDeploy)
		},
	}

	f := cmd.Flags()
	f.StringVar(&namespace, "namespace", "", "namespace of druid CR")
	f.StringVar(&cr, "cr", "", "name of the druid CR")
	f.StringVar(&deleteOrphanPvc, "deleteOrphanPvc", "", "deleteOrphanPvc is a bool, enabling to true will lead to delete of orphan pvc")
	f.StringVar(&rollingDeploy, "rollingDeploy", "", "rollingDeploy is a bool, enabling to true will lead to sequential rolling upgrades")

	return cmd
}

func (sv *druidListCmd) druidCRPatcherRun(namespace, CR, deleteOrphanPvc, rollingDeploy string) error {

	if deleteOrphanPvc != "" {
		b, _ := strconv.ParseBool(deleteOrphanPvc)
		patcherResult, err := di.patcherDruidDeleteOrphanPvc(namespace, CR, b)
		if err != nil {
			return err
		}

		if patcherResult {
			_, err := fmt.Fprintf(sv.out, "Druid CR [%s],successfully patched DeleteOrphanPvc to [%v] in Namespace [%s]\n", CR, deleteOrphanPvc, namespace)
			if err != nil {
				return err
			}
		}
	} else if rollingDeploy != "" {
		b, _ := strconv.ParseBool(rollingDeploy)
		patcherResult, err := di.patcherDruidRollingDeploy(namespace, CR, b)
		if err != nil {
			return err
		}

		if patcherResult {
			_, err := fmt.Fprintf(sv.out, "Druid CR [%s],successfully patched rollingDeploy to [%v] in Namespace [%s]\n", CR, rollingDeploy, namespace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
