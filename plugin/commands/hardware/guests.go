package hardware

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type GuestsCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewGuestsCommand(sl *metadata.SoftlayerCommand) (cmd *GuestsCommand) {
	thisCmd := &GuestsCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "guests " + T("IDENTIFIER"),
		Short: T("Lists the Virtual Guests running on this server."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *GuestsCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	hardwareGuestsResult, err := cmd.HardwareManager.GetHardwareGuests(hardwareId)
	if err != nil {
		return errors.NewAPIError(T("Failed to get the guests instances for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hardwareGuestsResult)
	}

	cmd.UI.Print("Hardware Guests")
	tableHardwareGuest := cmd.UI.Table([]string{T("id"), T("Hostname"), T("CPU"), T("Memory"), T("Start Date"), T("Status"), T("powerState")})
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
