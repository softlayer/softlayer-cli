package placementgroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "placement-group",
		Short: T("Classic infrastructure Placement Group"),
		RunE:  nil,
	}

	cobraCmd.AddCommand(NewPlacementGroupCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewPlacementGroupListCommand(sl).Command)
	cobraCmd.AddCommand(NewPlacementGroupDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewPlacementGroupCreateOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewPlacementGroupDetailCommand(sl).Command)
	return cobraCmd
}

func PlacementGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "placement-group",
		Description: T("Classic infrastructure Placement Group"),
	}
}
