package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "loadbal",
		Short: T("Classic infrastructure Load Balancers"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewHealthChecksCommand(sl).Command)
	cobraCmd.AddCommand(NewL7MembersAddCommand(sl).Command)
	cobraCmd.AddCommand(NewL7MembersDelCommand(sl).Command)
	cobraCmd.AddCommand(NewL7PolicyAddCommand(sl).Command)
	cobraCmd.AddCommand(NewL7PolicyDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewL7PolicyEditCommand(sl).Command)
	// cobraCmd.AddCommand(NewMembersAddCommand(sl).Command)
	// cobraCmd.AddCommand(NewMembersDelCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7PolicyListCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7PoolAddCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7PoolDelCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7PoolDetailCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7PoolEditCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7RuleAddCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7RuleDelCommand(sl).Command)
	// cobraCmd.AddCommand(NewL7RuleListCommand(sl).Command)
	// cobraCmd.AddCommand(NewListCommand(sl).Command)
	// cobraCmd.AddCommand(NewMembersAddCommand(sl).Command)
	// cobraCmd.AddCommand(NewMembersDelCommand(sl).Command)
	// cobraCmd.AddCommand(NewProtocolAddCommand(sl).Command)
	// cobraCmd.AddCommand(NewProtocolDeleteCommand(sl).Command)
	// cobraCmd.AddCommand(NewProtocolEditCommand(sl).Command)
	// cobraCmd.AddCommand(NewNetscalerDetailCommand(sl).Command)
	// cobraCmd.AddCommand(NewNetscalerListCommand(sl).Command)
	return cobraCmd
}

func LoadbalNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "loadbal",
		Description: T("Classic infrastructure Load Balancers"),
	}
}
