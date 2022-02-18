package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type PowerCycleCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewPowerCycleCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *PowerCycleCommand) {
	return &PowerCycleCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *PowerCycleCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will power cycle hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.PowerCycle(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to power cycle hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} was power cycle.", map[string]interface{}{"ID": hardwareId}))
	return nil
}

func HardwarePowerCycleMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "power-cycle",
		Description: T("Power cycle a server"),
		Usage:       "${COMMAND_NAME} sl hardware power-cycle IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
