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

// ToggleIPMICommand is the implementation of Softlayer `hardware toggle-ipmi` command
type ToggleIPMICommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

// NewToggleIPMICommand will create an instance of ToggleIPMICommand
func NewToggleIPMICommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *ToggleIPMICommand) {
	return &ToggleIPMICommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

// Run will execute ToggleIPMICommand
func (cmd *ToggleIPMICommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument"))
	}
	hardwareID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if c.IsSet("enable") && c.IsSet("disable") {
		return errors.NewExclusiveFlagsErrorWithDetails([]string{"--enable", "--disable"}, "")
	}

	if !c.IsSet("enable") && !c.IsSet("disable") {
		return errors.NewInvalidUsageError(T("Either '--enable' or '--disable' is required."))
	}

	enabled := !c.IsSet("disable")
	if err = cmd.HardwareManager.ToggleIPMI(hardwareID, enabled); err != nil {
		return cli.NewExitError(T("Failed to toggle IPMI interface of hardware server '{{.ID}}'.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully send request to toggle IPMI interface of hardware server '{{.ID}}'.", map[string]interface{}{"ID": hardwareID}))
	return nil
}

func HardwareToggleIPMIMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "toggle-ipmi",
		Description: T("Toggle the IPMI interface on and off. This command is asynchronous."),
		Usage:       "${COMMAND_NAME} sl hardware toggle-ipmi IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "enable",
				Usage: T("Enable the IPMI interface."),
			},
			cli.BoolFlag{
				Name:  "disable",
				Usage: T("Disable the IPMI interface."),
			},
			metadata.QuietFlag(),
		},
	}
}
