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

type ReloadCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewReloadCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *ReloadCommand) {
	return &ReloadCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *ReloadCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will reload operating system for hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.HardwareManager.Reload(hardwareId, c.String("i"), c.IntSlice("k"), c.IsSet("b"), c.IsSet("w"))
	if err != nil {
		return cli.NewExitError(T("Failed to reload operating system for hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Started to reload operating system for hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	return nil
}

func HardwareReloadMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "reload",
		Description: T("Reload operating system on a server"),
		Usage:       "${COMMAND_NAME} sl hardware reload IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "i,postinstall",
				Usage: T("Post-install script to download, only HTTPS executes, HTTP leaves file in /root"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("IDs of SSH key to add to the root user, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "b,upgrade-bios",
				Usage: T("Upgrade BIOS"),
			},
			cli.BoolFlag{
				Name:  "w,upgrade-firmware",
				Usage: T("Upgrade all hard drives' firmware"),
			},
			metadata.ForceFlag(),
		},
	}
}
