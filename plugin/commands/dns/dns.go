package dns

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "dns",
		Short: T("Classic infrastructure Domain Name System"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewImportCommand(sl).Command)
	cobraCmd.AddCommand(NewRecordAddCommand(sl).Command)
	cobraCmd.AddCommand(NewRecordEditCommand(sl).Command)
	cobraCmd.AddCommand(NewRecordListCommand(sl).Command)
	cobraCmd.AddCommand(NewRecordRemoveCommand(sl).Command)
	cobraCmd.AddCommand(NewZoneCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewZoneDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewZoneListCommand(sl).Command)
	cobraCmd.AddCommand(NewZonePrintCommand(sl).Command)
	return cobraCmd
}

func DnsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "dns",
		Description: T("Classic infrastructure Domain Name System"),
	}
}
