package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateOptionsCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPlacementGruopCreateOptionsCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PlacementGroupCreateOptionsCommand) {
	return &PlacementGroupCreateOptionsCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PlacementGroupCreateOptionsCommand) Run(c *cli.Context) error {
	datacenters, err := cmd.VirtualServerManager.GetDatacenters()
	if err != nil {
		return slErrors.NewInvalidUsageError("Internal error.")
	}
	tableRegion := cmd.UI.Table([]string{T("Datacenter"), T("Hostname"), T("BackendRouterId")})
	for _, datacenter := range datacenters{
		routers, err := cmd.VirtualServerManager.GetAvailablePlacementRouters(utils.IntPointertoInt(datacenter.Id))
		if err != nil {
			return slErrors.NewInvalidUsageError("Internal error.")
		}
		for _, routerAvalaible := range routers{
			tableRegion.Add(utils.FormatStringPointer(datacenter.LongName),utils.FormatStringPointer(routerAvalaible.Hostname),utils.FormatIntPointer(routerAvalaible.Id))
		}
	}

	rules, err := cmd.VirtualServerManager.GetRules()
	if err != nil {
		return slErrors.NewInvalidUsageError("Internal error.")
	}
	tableRules := cmd.UI.Table([]string{T("Id"), T("Rule")})
	for _, rule := range rules{
			tableRules.Add(utils.FormatIntPointer(rule.Id),utils.FormatStringPointer(rule.KeyName))

	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return slErrors.NewInvalidUsageError("Internal error.")
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSONList(cmd.UI, datacenters)
	}

	tableRegion.Print()
	tableRules.Print()
	return nil
}

func VSPlacementGroupCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "placementgroup-create-options",
		Description: T("Get List options for creating a placement group.."),
		Usage: T(`${COMMAND_NAME} sl vs placementgroup-create-options
EXAMPLE:
   ${COMMAND_NAME} sl vs placementgroup-create-options
    Get List options for creating a placement group.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		}}
}