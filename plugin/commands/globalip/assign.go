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

type AssignCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewAssignCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *AssignCommand) {
	return &AssignCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *AssignCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	globalIPID, err := utils.ResolveGloablIPId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Globalip ID")
	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	targetIPAddress := c.Args()[1]
	resp, err := cmd.NetworkManager.AssignGlobalIP(globalIPID, targetIPAddress)
	if err != nil {
		return cli.NewExitError(T("Failed to assign global IP {{.IpID}} to target {{.Target}}.\n",
			map[string]interface{}{"IpID": globalIPID, "Target": targetIPAddress})+err.Error(), 2)

	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to modify a global IP route is created, routes will be updated in one or two minutes."))
	return nil
}

func GlobalIpAssignMetaData() cli.Command {
	return cli.Command{
		Category:    "globalip",
		Name:        "assign",
		Description: T("Assign a global IP to a target router or device"),
		Usage: T(`${COMMAND_NAME} sl globalip assign IDENTIFIER TARGET [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip assign 12345678 9.111.123.456
	This command assigns IP address with ID 12345678 to a target device whose IP address is 9.111.123.456.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
