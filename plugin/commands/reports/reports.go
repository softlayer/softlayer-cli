package reports

import (
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "report",
		Short: T("Classic Infrastructure Reports"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewDCClosuresCommand(sl).Command)
	cobraCmd.AddCommand(NewBandwidthCommand(sl).Command)
	return cobraCmd
}

func ReportsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "report",
		Description: T("Classic Infrastructure Reports"),
	}
}
