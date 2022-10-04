package hardware

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ReflashFirmwareCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	ForceFlag       bool
}

func NewReflashFirmwareCommand(sl *metadata.SoftlayerCommand) (cmd *ReflashFirmwareCommand) {
	thisCmd := &ReflashFirmwareCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "reflash-firmware " + T("IDENTIFIER"),
		Short: T("Reflash server firmware."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReflashFirmwareCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !cmd.ForceFlag {
		hardwareMapValue := map[string]interface{}{"hardwareID": hardwareId}
		confirm, err := cmd.UI.Confirm(T("This will power off the server with id {{.hardwareID}} and reflash device firmware. Continue?", hardwareMapValue))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	response, err := cmd.HardwareManager.CreateFirmwareReflashTransaction(hardwareId)
	if err != nil {
		return errors.NewAPIError(T("Failed to reflash firmware."), err.Error(), 2)
	}
	if response {
		cmd.UI.Print(T("Successfully device firmware reflashed"))
	}

	return nil
}
