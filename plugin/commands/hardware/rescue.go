package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RescueCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewRescueCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *RescueCommand) {
	return &RescueCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *RescueCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument"))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will reboot hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.Rescure(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to rescue hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} was rebooted to a rescue image.", map[string]interface{}{"ID": hardwareId}))
	return nil
}

func HardwareRescueMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "rescue",
		Description: T("Reboot server into a rescue image"),
		Usage:       "${COMMAND_NAME} sl hardware rescue IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
