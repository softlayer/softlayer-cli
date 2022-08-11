package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ReflashFirmwareCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewReflashFirmwareCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *ReflashFirmwareCommand) {
	return &ReflashFirmwareCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *ReflashFirmwareCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		hardwareMapValue := map[string]interface{}{"hardwareID": hardwareId}
		confirm, err := cmd.UI.Confirm(T("This will power off the server with id {{.hardwareID}} and reflash device firmware. Continue?", hardwareMapValue))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	response, err := cmd.HardwareManager.CreateFirmwareReflashTransaction(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to reflash firmware.")+"\n"+err.Error(), 2)
	}
	if response {
		cmd.UI.Print(T("Successfully device firmware reflashed"))
	}

	return nil
}

func HardwareReflashFirmwareMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "reflash-firmware",
		Description: T("Reflash server firmware."),
		Usage: T(`${COMMAND_NAME} sl hardware reflash-firmware IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl hardware reflash-firmware 123456`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
