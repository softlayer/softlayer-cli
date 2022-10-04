package security

import (
	"strconv"

	"github.com/spf13/cobra"
	
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type KeyRemoveCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Force           bool
}

func NewKeyRemoveCommand(sl *metadata.SoftlayerCommand) *KeyRemoveCommand {
	thisCmd := &KeyRemoveCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "sshkey-remove " + T("IDENTIFIER"),
		Short: T("Permanently removes an SSH key"),
		Long: T(`${COMMAND_NAME} sl security sshkey-remove IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-remove 12345678 -f 
   This command removes the SSH key with ID 12345678 without asking for confirmation.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *KeyRemoveCommand) Run(args []string) error {
	keyID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will remove SSH key: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": keyID}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.SecurityManager.DeleteSSHKey(keyID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to remove SSH key: {{.ID}}.\n", map[string]interface{}{"ID": keyID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key {{.ID}} was removed.", map[string]interface{}{"ID": keyID}))
	return nil
}
