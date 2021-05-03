package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type UpdateFirmwareCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewUpdateFirmwareCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *UpdateFirmwareCommand) {
	return &UpdateFirmwareCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *UpdateFirmwareCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will power off hardware server: {{.ID}} and update device firmware. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.UpdateFirmware(hardwareId, true, true, true, true)
	if err != nil {
		return cli.NewExitError(T("Failed to update firmware for hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Started to update firmware for hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
