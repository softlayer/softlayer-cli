package subnet

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditIpCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewEditIpCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *EditIpCommand) {
	return &EditIpCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *EditIpCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	if !c.IsSet("note") {
		return errors.NewMissingInputError(T("--note"))
	}

	subnetIpAddressID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		ipAddress := c.Args()[0]
		subnetIpAddress, err := cmd.NetworkManager.GetIpByAddress(ipAddress)
		if err != nil {
			return cli.NewExitError(T("Failed to get Subnet IP by address\n")+err.Error(), 2)
		}
		subnetIpAddressID = *subnetIpAddress.Id
	}

	note := c.String("note")
	subnetIpAddressTemplate := datatypes.Network_Subnet_IpAddress{
		Note: sl.String(note),
	}
	response, err := cmd.NetworkManager.EditSubnetIpAddress(subnetIpAddressID, subnetIpAddressTemplate)
	if err != nil {
		return cli.NewExitError(T("Failed to set note: {{.note}}.\n", map[string]interface{}{"note": note})+err.Error(), 2)
	}
	if response {
		cmd.UI.Ok()
		cmd.UI.Print(T("Set note successfully"))
	}
	return nil
}

func SubnetEditIpMetaData() cli.Command {
	return cli.Command{
		Category:    "subnet",
		Name:        "edit-ip",
		Description: T("Set the note of the ipAddress."),
		Usage: T(`${COMMAND_NAME} sl subnet edit IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet edit-ip 11.22.33.44 --note myNote
   ${COMMAND_NAME} sl subnet edit-ip 12345678 --note myNote`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "note",
				Usage: T("The note "),
			},
		},
	}
}
