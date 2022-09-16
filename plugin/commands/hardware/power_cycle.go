package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type PowerCycleCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	ForceFlag       bool
}

func NewPowerCycleCommand(sl *metadata.SoftlayerCommand) (cmd *PowerCycleCommand) {
	thisCmd := &PowerCycleCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "power-cycle " + T("IDENTIFIER"),
		Short: T("Power cycle a server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PowerCycleCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will power cycle hardware server: {{.ID}}. Continue?", map[string]interface{}{"ID": hardwareId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.HardwareManager.PowerCycle(hardwareId)
	if err != nil {
		return errors.NewAPIError(T("Failed to power cycle hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server: {{.ID}} was power cycle.", map[string]interface{}{"ID": hardwareId}))
	return nil
}
