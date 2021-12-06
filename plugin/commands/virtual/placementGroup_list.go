package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupListCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPlacementGroupListCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PlacementGroupListCommand) {
	return &PlacementGroupListCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PlacementGroupListCommand) Run(c *cli.Context) error {
	placements, err := cmd.VirtualServerManager.PlacementsGroupList(utils.EMPTY_STRING)
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Failed to get virtual Placement groups on your account.\n")
	}
	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Backend Router"),
		T("Rule"), T("Guests"), T("Created")})
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
	table.Print()
	return nil
}
