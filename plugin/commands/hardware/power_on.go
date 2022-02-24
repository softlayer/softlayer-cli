package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type PowerOnCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewPowerOnCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *PowerOnCommand) {
	return &PowerOnCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *PowerOnCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	err = cmd.HardwareManager.PowerOn(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to power on hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} is power on.", map[string]interface{}{"ID": hardwareId}))
	return nil
}

func HardwarePowerOnMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "power-on",
		Description: T("Power on a server"),
		Usage:       "${COMMAND_NAME} sl hardware power-on IDENTIFIER",
	}
}
