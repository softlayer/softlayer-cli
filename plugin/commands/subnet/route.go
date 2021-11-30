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

type RouteCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewRouteCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *RouteCommand) {
	return &RouteCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *RouteCommand) Run(c *cli.Context) error {
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

	if !c.IsSet("t") {
		return errors.NewInvalidUsageError(T("[-t/--type] is required."))
	}

	if !c.IsSet("i") {
		return errors.NewInvalidUsageError(T("[-i/--type-id] is required."))
	}

	resp, err := cmd.NetworkManager.Route(subnetID, c.String("t"), c.String("i"))
	if err != nil {
		return cli.NewExitError(T("Failed to route using the type: {{.TYPE}} and identifier: {{.IDENTIFIER}}.\n",
			map[string]interface{}{"TYPE": c.String("t"), "IDENTIFIER": c.String("i")})+err.Error(), 2)

	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to route is created, routes will be updated in one or two minutes."))
	return nil
}
