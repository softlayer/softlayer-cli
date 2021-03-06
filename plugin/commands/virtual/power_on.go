package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PowerOnCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPowerOnCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PowerOnCommand) {
	return &PowerOnCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PowerOnCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will power on virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.PowerOnInstance(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to power on virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was power on.", map[string]interface{}{"VsId": vsID}))
	return nil
}

func VSPowerOnMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "power-on",
		Description: T("Power on a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs power-on IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs power-on 12345678
   This command performs a power on for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}