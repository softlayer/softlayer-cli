package placementgroup

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateOptionsCommand struct {
	UI                terminal.UI
	PlaceGroupManager managers.PlaceGroupManager
}

func NewPlacementGroupCreateOptionsCommand(ui terminal.UI, placeGroupManager managers.PlaceGroupManager) (cmd *PlacementGroupCreateOptionsCommand) {
	return &PlacementGroupCreateOptionsCommand{
		UI:                ui,
		PlaceGroupManager: placeGroupManager,
	}
}

func (cmd *PlacementGroupCreateOptionsCommand) Run(c *cli.Context) error {

	routers, err := cmd.PlaceGroupManager.GetRouter(nil, "")
	if err != nil {
		return cli.NewExitError(T("Failed to list routers\n")+err.Error(), 2)
	}

	cmd.UI.Print("Available Router:")
	if len(routers) == 0 {
		cmd.UI.Print(T("No available router was found."))
	} else {
		sort.Sort(utils.RouterByHostname(routers))
		table := cmd.UI.Table([]string{T("Data Center"), T("Host Name"), T("Backend Router Id")})
		for _, router := range routers {
			center := "-"
			if router.TopLevelLocation != nil {
				center = utils.FormatStringPointer(router.TopLevelLocation.LongName)
			}
			table.Add(center,
				utils.FormatStringPointer(router.Hostname),
				utils.FormatIntPointer(router.Id))
		}
		table.Print()
	}

	rules, err := cmd.PlaceGroupManager.GetRules()
	if err != nil {
		return cli.NewExitError(T("Failed to list rules\n")+err.Error(), 2)
	}
	cmd.UI.Print("Rules:")
	if len(rules) == 0 {
		cmd.UI.Print(T("No rules was found."))
	} else {
		ruleTable := cmd.UI.Table([]string{T("ID"), T("Name")})
		for _, rule := range rules {
			ruleTable.Add(utils.FormatIntPointer(rule.Id),
				utils.FormatStringPointer(rule.KeyName),
			)
		}
		ruleTable.Print()
	}

	return nil
}

func PlacementGroupCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "placement-group",
		Name:        "create-options",
		Description: T("List options for creating a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group create-options",
		Flags:       []cli.Flag{},
	}
}
