package globalip

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UnassignCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewUnassignCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *UnassignCommand) {
	return &UnassignCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *UnassignCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	globalIPID, err := utils.ResolveGloablIPId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Globalip ID")
	}

	resp, err := cmd.NetworkManager.UnassignGlobalIP(globalIPID)
	if err != nil {
		return cli.NewExitError(T("Failed to unassign global IP {{.ID}}.\n", map[string]interface{}{"ID": globalIPID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to unroute a global IP address is created, routes will be updated in one or two minutes."))
	return nil
}
