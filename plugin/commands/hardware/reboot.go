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

type RebootCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Hard            bool
	Soft            bool
	ForceFlag       bool
}

func NewRebootCommand(sl *metadata.SoftlayerCommand) (cmd *RebootCommand) {

	thisCmd := &RebootCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}


	cobraCmd := &cobra.Command{
		Use:   "reboot " + T("IDENTIFIER"),
		Short: T("Reboot an active server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Hard, "hard", false, T("Perform a hard reboot"))
	cobraCmd.Flags().BoolVar(&thisCmd.Soft, "soft", false, T("Perform a soft reboot"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RebootCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if cmd.Hard && cmd.Soft {
		return errors.NewInvalidUsageError(T("Can only specify either --hard or --soft."))
	}

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will reboot hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.Reboot(hardwareId, cmd.Soft, cmd.Hard)
	if err != nil {
		return errors.NewAPIError(T("Failed to reboot hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} was rebooted.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
