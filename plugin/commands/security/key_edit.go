package security

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type KeyEditCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Label           string
	Note            string
}

func NewKeyEditCommand(sl *metadata.SoftlayerCommand) *KeyEditCommand {
	thisCmd := &KeyEditCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "sshkey-edit " + T("IDENTIFIER"),
		Short: T("Edit an SSH key"),
		Long: T(`${COMMAND_NAME} sl security sshkey-edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-edit 12345678 --label IBMCloud --note testing
   This command updates the SSH key with ID 12345678 and sets label to "IBMCloud" and note to "testing".`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Label, "label", "", T("The new label for the key"))
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("New notes for the key"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *KeyEditCommand) Run(args []string) error {
	keyID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSH Key ID")
	}
	if cmd.Label == "" && cmd.Note == "" {
		return errors.NewInvalidUsageError(T("either [--label] or [--note] must be specified to edit SSH key."))
	}
	err = cmd.SecurityManager.EditSSHKey(keyID, cmd.Label, cmd.Note)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit SSH key: {{.ID}}.\n", map[string]interface{}{"ID": keyID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSH key {{.ID}} was updated.", map[string]interface{}{"ID": keyID}))
	return nil
}
