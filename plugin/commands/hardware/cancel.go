package hardware

import (
	"strconv"
	"strings"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CancelCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewCancelCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the hardware server: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": hardwareID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.HardwareManager.CancelHardware(hardwareID, c.String("r"), c.String("c"), c.IsSet("i"))
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find hardware server with ID: {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server {{.ID}} was cancelled.", map[string]interface{}{"ID": hardwareID}))
	return nil
}
