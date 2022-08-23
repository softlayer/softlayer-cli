package dedicatedhost

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "dedicatedhost",
		Short: T("Classic infrastructure Dedicatedhost"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewListGuestsCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateOptionsCommand(sl).Command)
	return cobraCmd
}

func DedicatedhostNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "dedicatedhost",
		Description: T("Classic infrastructure Dedicatedhost"),
	}
}
