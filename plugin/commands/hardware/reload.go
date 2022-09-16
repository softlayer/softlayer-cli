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

type ReloadCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Postinstall     string
	Key             []int
	UpgradeBios     bool
	UpgradeFirmware bool
	ForceFlag       bool
}

func NewReloadCommand(sl *metadata.SoftlayerCommand) (cmd *ReloadCommand) {
	thisCmd := &ReloadCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "reload " + T("IDENTIFIER"),
		Short: T("Reload operating system on a server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Postinstall, "postinstall", "i", "", T("Post-install script to download, only HTTPS executes, HTTP leaves file in /root"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Key, "key", "k", []int{}, T("IDs of SSH key to add to the root user, multiple occurrence allowed"))
	cobraCmd.Flags().BoolVarP(&thisCmd.UpgradeBios, "upgrade-bios", "b", false, T("Upgrade BIOS"))
	cobraCmd.Flags().BoolVarP(&thisCmd.UpgradeFirmware, "upgrade-firmware", "w", false, T("Upgrade all hard drives' firmware"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReloadCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will reload operating system for hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.HardwareManager.Reload(hardwareId, cmd.Postinstall, cmd.Key, cmd.UpgradeBios, cmd.UpgradeFirmware)
	if err != nil {
		return errors.NewAPIError(T("Failed to reload operating system for hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Started to reload operating system for hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
