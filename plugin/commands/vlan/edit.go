package vlan

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewEditCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	vlanID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}
	if !c.IsSet("n") {
		return errors.NewMissingInputError("-n|--name")

	}

	err = cmd.NetworkManager.EditVlan(vlanID, c.String("n"))
	if err != nil {
		return cli.NewExitError(T("Failed to edit VLAN: {{.VlanID}}.\n", map[string]interface{}{"VlanID": vlanID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("VLAN {{.VlanID}} was updated.", map[string]interface{}{"VlanID": vlanID}))
	return nil
}

func VlanEditMetaData() cli.Command {
	return cli.Command{
		Category:    "vlan",
		Name:        "edit",
		Description: T("Edit the details about a VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan edit 12345678 -n myvlan-rename
   This command updates vlan with ID 12345678 and gives it a new name "myvlan-rename".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the VLAN"),
			},
		},
	}
}
