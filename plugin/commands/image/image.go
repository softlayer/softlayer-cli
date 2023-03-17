package image

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "image",
		Short: T("Classic infrastructure Compute images"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewExportCommand(sl).Command)
	cobraCmd.AddCommand(NewImportCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewDatacenterCommand(sl).Command)
	cobraCmd.AddCommand(NewShareCommand(sl).Command)
	return cobraCmd
}

func ImageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "image",
		Description: T("Classic infrastructure Compute images"),
	}
}
