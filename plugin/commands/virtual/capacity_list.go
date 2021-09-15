package virtual

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityListCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCapacityListCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CapacityListCommand) {
	return &CapacityListCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CapacityListCommand) Run(c *cli.Context) error {
	capacities, err := cmd.VirtualServerManager.CapacityList("")
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Failed to get virtual Reserved capacity groups on your account.\n")
	}
	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Capacity"),
		T("Flavor"), T("Location"), T("Created")})
	for _, capacity := range capacities {
		var occupied string
		var available string
		if utils.FormatUIntPointer(capacity.OccupiedInstanceCount) == "0" {
			occupied = "#"
		}
		if  utils.FormatUIntPointer(capacity.AvailableInstanceCount) == "0" {
			available = "-"
		}
		table.Add(utils.FormatIntPointer(capacity.Id),
			utils.FormatStringPointer(capacity.Name),
			fmt.Sprintf("%s%s",occupied,available),
			utils.FormatStringPointer(capacity.Instances[0].BillingItem.Description),
			utils.FormatStringPointer(capacity.BackendRouter.Hostname),
			utils.FormatSLTimePointer(capacity.CreateDate))
	}
	table.Print()
	return nil
}
