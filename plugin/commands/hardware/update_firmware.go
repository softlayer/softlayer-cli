package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type UpdateFirmwareCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	ForceFlag       bool
}

func NewUpdateFirmwareCommand(sl *metadata.SoftlayerCommand) (cmd *UpdateFirmwareCommand) {
	thisCmd := &UpdateFirmwareCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "update-firmware " + T("IDENTIFIER"),
		Short: T("Update server firmware"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UpdateFirmwareCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will power off hardware server: {{.ID}} and update device firmware. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.UpdateFirmware(hardwareId, true, true, true, true)
	if err != nil {
		return errors.NewAPIError(T("Failed to update firmware for hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Started to update firmware for hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
