package securitygroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "securitygroup",
		Short: T("Classic infrastructure network security groups"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewInterfaceAddCommand(sl).Command)
	cobraCmd.AddCommand(NewInterfaceListCommand(sl).Command)
	//cobraCmd.AddCommand(NewInterfaceRemoveCommand(sl).Command)
	//cobraCmd.AddCommand(NewListCommand(sl).Command)
	//cobraCmd.AddCommand(NewRuleAddCommand(sl).Command)
	//cobraCmd.AddCommand(NewRuleEditCommand(sl).Command)
	//cobraCmd.AddCommand(NewRuleListCommand(sl).Command)
	//cobraCmd.AddCommand(NewRuleRemoveCommand(sl).Command)
	return cobraCmd
}

func SecurityGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "securitygroup",
		Description: T("Classic infrastructure network security groups"),
	}
}
