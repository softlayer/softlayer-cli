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

type RebootCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewRebootCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *RebootCommand) {
	return &RebootCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *RebootCommand) Run(c *cli.Context) error {
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
		confirm, err := cmd.UI.Confirm(T("This will reboot virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.RebootInstance(vsID, c.IsSet("soft"), c.IsSet("hard"))
	if err != nil {
		return cli.NewExitError(T("Failed to reboot virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was rebooted.", map[string]interface{}{"VsId": vsID}))
	return nil
}

func VSRebootMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "reboot",
		Description: T("Reboot an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs reboot IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs reboot 12345678 --hard
   This command performs a hard reboot for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "hard",
				Usage: T("Perform a hard reboot"),
			},
			cli.BoolFlag{
				Name:  "soft",
				Usage: T("Perform a soft reboot"),
			},
			metadata.ForceFlag(),
		},
	}
}