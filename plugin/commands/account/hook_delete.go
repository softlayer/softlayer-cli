package account

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type HookDeleteCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewHookDeleteCommand(sl *metadata.SoftlayerCommand) *HookDeleteCommand {
	thisCmd := &HookDeleteCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "hook-delete " + T("IDENTIFIER"),
		Short: T("Delete a provisioning scriptt."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *HookDeleteCommand) Run(args []string) error {

	hookID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Hook ID")
	}

	success, err := cmd.AccountManager.DeleteProvisioningScript(hookID)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete Provisioning Hook."), err.Error(), 2)
	}

	if success {
		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully removed Provisioning Hook."))
	}
	return nil
}
