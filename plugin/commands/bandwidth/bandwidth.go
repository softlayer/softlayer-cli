package bandwidth

import (
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "bandwidth",
		Short: T("Classic infrastructure Bandwidth commands"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewPoolsCommand(sl).Command)
	cobraCmd.AddCommand(NewPoolsDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewSummaryCommand(sl).Command)
	cobraCmd.AddCommand(NewPoolsCreateCommand(sl).Command)
	return cobraCmd
}

func BandwidthNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "bandwidth",
		Description: T("Classic Infrastructure Bandwidth commands"),
	}
}
