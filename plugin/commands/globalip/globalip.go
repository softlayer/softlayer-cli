package globalip

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "globalip",
		Short: T("Classic infrastructure Global IP addresses"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAssignCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewUnassignCommand(sl).Command)
	return cobraCmd
}

func GlobalIpNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "globalip",
		Description: T("Classic infrastructure Global IP addresses"),
	}
}
