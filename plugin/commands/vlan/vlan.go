package vlan

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "vlan",
		Short: T("Classic infrastructure Network VLANs"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewOptionsCommand(sl).Command)
	return cobraCmd
}

func VlanNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "vlan",
		Description: T("Classic infrastructure Network VLANs"),
	}
}
