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

type PauseCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPauseCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PauseCommand) {
	return &PauseCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PauseCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will pause virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.PauseInstance(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to pause virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was paused.", map[string]interface{}{"VsId": vsID}))

	return nil
}

func VSPauseMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "pause",
		Description: T("Pause an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs pause IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs pause 12345678 -f
   This command pauses virtual server instance with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
