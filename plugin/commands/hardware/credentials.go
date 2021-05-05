package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialsCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewCredentialsCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *CredentialsCommand) {
	return &CredentialsCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *CredentialsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hardware, err := cmd.HardwareManager.GetHardware(hardwareID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		if hardware.OperatingSystem != nil {
			return utils.PrintPrettyJSON(cmd.UI, hardware.OperatingSystem.Passwords)
		}
		return utils.PrintPrettyJSON(cmd.UI, hardware.OperatingSystem)
	}

	if hardware.OperatingSystem != nil && hardware.OperatingSystem.Passwords != nil {
		table := cmd.UI.Table([]string{T("Username"), T("Password")})
		for _, item := range hardware.OperatingSystem.Passwords {
			if item.Username != nil && item.Password != nil {
				table.Add(*item.Username, *item.Password)
			}
		}
		table.Print()
		return nil
	}
	return cli.NewExitError(T("Failed to find credentials of hardware server {{.ID}}.", map[string]interface{}{"ID": hardwareID}), 2)
}
