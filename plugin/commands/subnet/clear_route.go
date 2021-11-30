package subnet

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ClearRouteCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewClearRouteCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *ClearRouteCommand) {
	return &ClearRouteCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *ClearRouteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	subnetID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.NetworkManager.ClearRoute(subnetID)
	if err != nil {
		return cli.NewExitError(T("Failed to clear the route for the subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to clear the route is created, routes will be updated in one or two minutes."))
	return nil
}
