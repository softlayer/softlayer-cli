package firewall

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "firewall",
		Short: T("Classic infrastructure Firewalls"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAddCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewEditCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	return cobraCmd
}

func FirewallNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "firewall",
		Description: T("Classic infrastructure Firewalls"),
	}
}
