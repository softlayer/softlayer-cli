package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupListCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewPlacementGroupListCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupListCommand) {
	thisCmd := &PlacementGroupListCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "placementgroup-list",
		Short: T("List placement groups."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupListCommand) Run(args []string) error {
	placements, err := cmd.VirtualServerManager.PlacementsGroupList(utils.EMPTY_STRING)
	if err != nil {
		return slErrors.NewAPIError("Failed to get virtual Placement groups on your account.\n", err.Error(), 2)
	}
	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Backend Router"), T("Rule"), T("Guests"), T("Created")})
	for _, placement := range placements {
		var backendName string
		if placement.BackendRouter == nil {
			backendName = "-"
		} else {
			backendName = utils.FormatStringPointer(placement.BackendRouter.Hostname)
		}
		table.Add(utils.FormatIntPointer(placement.Id),
			utils.FormatStringPointer(placement.Name),
			backendName,
			utils.FormatStringPointer(placement.Rule.Name),
			utils.FormatUIntPointer(placement.GuestCount),
			utils.FormatSLTimePointer(placement.CreateDate))
	}
	outputFormat := cmd.GetOutputFlag()
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
