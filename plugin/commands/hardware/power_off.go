package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type PowerOffCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewPowerOffCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *PowerOffCommand) {
	return &PowerOffCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *PowerOffCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will power off hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.PowerOff(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to power off hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} was power off.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
