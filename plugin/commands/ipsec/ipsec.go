package ipsec

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "ipsec",
		Short: T("Classic infrastructure IPSEC VPN"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewConfigCommand(sl).Command)
	cobraCmd.AddCommand(NewCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewOrderCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewListCommand(sl).Command)
	cobraCmd.AddCommand(NewAddSubnetCommand(sl).Command)
	cobraCmd.AddCommand(NewRemoveSubnetCommand(sl).Command)
	cobraCmd.AddCommand(NewAddTranslationCommand(sl).Command)
	cobraCmd.AddCommand(NewRemoveTranslationCommand(sl).Command)
	cobraCmd.AddCommand(NewUpdateTranslationCommand(sl).Command)
	cobraCmd.AddCommand(NewUpdateCommand(sl).Command)
	return cobraCmd
}

func IpsecNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "ipsec",
		Description: T("Classic infrastructure IPSEC VPN"),
	}
}
