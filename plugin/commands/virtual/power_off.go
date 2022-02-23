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

type PowerOffCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPowerOffCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PowerOffCommand) {
	return &PowerOffCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PowerOffCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if c.IsSet("hard") && c.IsSet("soft") {
		return bmxErr.NewExclusiveFlagsError("--hard", "--soft")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will power off virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.PowerOffInstance(vsID, c.IsSet("soft"), c.IsSet("hard"))
	if err != nil {
		return cli.NewExitError(T("Failed to power off virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was power off.", map[string]interface{}{"VsId": vsID}))
	return nil
}

func VSPowerOffMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "power-off",
		Description: T("Power off an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs power-off IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs power-off 12345678 --soft
   This command performs a soft power off for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "hard",
				Usage: T("Perform a hard shutdown"),
			},
			cli.BoolFlag{
				Name:  "soft",
				Usage: T("Perform a soft shutdown"),
			},
			metadata.ForceFlag(),
		},
	}
}