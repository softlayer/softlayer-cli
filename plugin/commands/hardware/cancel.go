package hardware

import (
	"strconv"
	"strings"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Immediate       bool
	Reason          string
	Comment         string
	ForceFlag       bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) (cmd *CancelCommand) {
	thisCmd := &CancelCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel a hardware server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.Immediate, "immediate", "i", false, T("Cancels the server immediately (instead of on the billing anniversary)"))
	cobraCmd.Flags().StringVarP(&thisCmd.Reason, "reason", "r", "", T("An optional cancellation reason. See '${COMMAND_NAME} sl hardware cancel-reasons' for a list of available options"))
	cobraCmd.Flags().StringVarP(&thisCmd.Comment, "comment", "c", "", T("An optional comment to add to the cancellation ticket"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will cancel the hardware server: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": hardwareID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.HardwareManager.CancelHardware(hardwareID, cmd.Reason, cmd.Comment, cmd.Immediate)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find hardware server with ID: {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to cancel hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Hardware server {{.ID}} was cancelled.", map[string]interface{}{"ID": hardwareID}))
	return nil
}
