package placementgroup

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupListCommand struct {
	*metadata.SoftlayerCommand
	PlaceGroupManager managers.PlaceGroupManager
	Command           *cobra.Command
}

func NewPlacementGroupListCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupListCommand) {
	thisCmd := &PlacementGroupListCommand{
		SoftlayerCommand:  sl,
		PlaceGroupManager: managers.NewPlaceGroupManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List placement groups"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	placementGroups, err := cmd.PlaceGroupManager.List("")
	if err != nil {
		return errors.NewAPIError(T("Failed to list placement groups"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, placementGroups)
	}

	cmd.UI.Ok()
	if len(placementGroups) == 0 {
		cmd.UI.Print(T("No placement group was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Backend Router"), T("Rule"), T("Guests"), T("Created")})
		for _, placementGroup := range placementGroups {
			backendRouter := "-"
			rule := "-"
			if placementGroup.BackendRouter != nil {
				backendRouter = utils.FormatStringPointer(placementGroup.BackendRouter.Hostname)
			}
			if placementGroup.Rule != nil {
				rule = utils.FormatStringPointer(placementGroup.Rule.Name)
			}

			table.Add(utils.FormatIntPointer(placementGroup.Id),
				utils.FormatStringPointer(placementGroup.Name),
				backendRouter,
				rule,
				utils.FormatUIntPointer(placementGroup.GuestCount),
				utils.FormatSLTimePointer(placementGroup.CreateDate),
			)
		}
		table.Print()
	}

	return nil
}
