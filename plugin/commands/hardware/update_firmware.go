package hardware

import (
	"strconv"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type UpdateFirmwareCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	ForceFlag       bool
	IPMIFlag bool
	RAIDFlag bool
	BIOSFlag bool
	HDFlag bool
	NetworkFlag bool
}

func NewUpdateFirmwareCommand(sl *metadata.SoftlayerCommand) (cmd *UpdateFirmwareCommand) {
	thisCmd := &UpdateFirmwareCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "update-firmware " + T("IDENTIFIER"),
		Short: T("Update server firmware"),
		Long: T("Update server firmware. By default will update all available server components."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))
	cobraCmd.Flags().BoolVarP(&thisCmd.IPMIFlag, "ipmi", "i", false, T("Update IPMI firmware"))
	cobraCmd.Flags().BoolVarP(&thisCmd.RAIDFlag, "raid", "r", false, T("Update RAID firmware"))
	cobraCmd.Flags().BoolVarP(&thisCmd.BIOSFlag, "bios", "b", false, T("Update BIOS firmware"))
	cobraCmd.Flags().BoolVarP(&thisCmd.HDFlag, "harddrive", "d", false, T("Update Hard Drive firmware"))
	cobraCmd.Flags().BoolVarP(&thisCmd.NetworkFlag, "network", "n", false, T("Update Network Card firmware"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UpdateFirmwareCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	// No options specified, set them all to true
	if !(cmd.IPMIFlag || cmd.RAIDFlag || cmd.BIOSFlag || cmd.HDFlag || cmd.NetworkFlag) {
		cmd.IPMIFlag = true
		cmd.RAIDFlag = true
		cmd.BIOSFlag = true
		cmd.HDFlag = true
		cmd.NetworkFlag = true
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

	err = cmd.HardwareManager.UpdateFirmware(hardwareId, cmd.IPMIFlag, cmd.RAIDFlag, cmd.BIOSFlag, cmd.HDFlag, cmd.NetworkFlag)
	if err != nil {
		return err
	}
	cmd.UI.Print(T("Started to update firmware for hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
