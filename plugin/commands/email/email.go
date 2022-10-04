package email

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "email",
		Short: T("Classic infrastructure Email commands"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	return cobraCmd
}

func EmailNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "email",
		Description: T("Classic infrastructure Email commands"),
	}
}
