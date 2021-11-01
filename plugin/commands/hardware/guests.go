package hardware

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type GuestsCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewGuestsCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *GuestsCommand) {
	return &GuestsCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *GuestsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hardwareGuestsResult, err := cmd.HardwareManager.GetHardwareGuests(hardwareId)
	if err != nil {
		return cli.NewExitError(T("Failed to get the guests instances for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hardwareGuestsResult)
	}

	cmd.UI.Print("Hardware Guests")
	tableHardwareGuest := cmd.UI.Table([]string{T("id"), T("hostname"), T("CPU"), T("Memory"), T("Start Date"), T("Status"), T("powerState")})
	for _, guest := range hardwareGuestsResult {
		tableHardwareGuest.Add(
			utils.FormatIntPointer(guest.Id),
			*guest.Hostname,
			fmt.Sprintf("%d %s", *guest.MaxCpu, *guest.MaxCpuUnits),
			utils.FormatIntPointer(guest.MaxMemory),
			utils.FormatSLTimePointer(guest.CreateDate),
			*guest.Status.KeyName,
			*guest.PowerState.KeyName)
	}
	tableHardwareGuest.Print()

	return nil
}
