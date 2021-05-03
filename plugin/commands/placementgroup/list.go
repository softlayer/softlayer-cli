package placementgroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type PlacementGroupListCommand struct {
	UI                terminal.UI
	PlaceGroupManager managers.PlaceGroupManager
}

func NewPlacementGroupListCommand(ui terminal.UI, placeGroupManager managers.PlaceGroupManager) (cmd *PlacementGroupListCommand) {
	return &PlacementGroupListCommand{
		UI:                ui,
		PlaceGroupManager: placeGroupManager,
	}
}

func (cmd *PlacementGroupListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	placementGroups, err := cmd.PlaceGroupManager.List("")
	if err != nil {
		return cli.NewExitError(T("Failed to list placement groups\n")+err.Error(), 2)
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
