package objectstorage

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CredentialDeleteCommand struct {
	*metadata.SoftlayerCommand
	ObjectStorageManager managers.ObjectStorageManager
	Command              *cobra.Command
	CredentialId         int
}

func NewCredentialDeleteCommand(sl *metadata.SoftlayerCommand) *CredentialDeleteCommand {
	thisCmd := &CredentialDeleteCommand{
		SoftlayerCommand:     sl,
		ObjectStorageManager: managers.NewObjectStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "credential-delete",
		Short: T("Delete the credential of an Object Storage Account."),
		Long: T(`${COMMAND_NAME} sl object-storage credential-delete IDENTIFIER [OPTIONS]

Examples:
	${COMMAND_NAME} sl object-storage credential-delete 123456 --credential-id 654321`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.CredentialId, "credential-id", 0, T("This is the credential id associated with the volume. [Required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialDeleteCommand) Run(args []string) error {
	storageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	if cmd.CredentialId == 0 {
		return slErr.NewMissingInputError("--credential-id")
	}

	credentialID := cmd.CredentialId

	err = cmd.ObjectStorageManager.DeleteCredential(storageID, credentialID)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find object-storage with ID: {{.storageID}}.\n", map[string]interface{}{"storageID": storageID}), err.Error(), 0)
		}
		if strings.Contains(err.Error(), "ObjectNotFound") {
			return errors.NewAPIError(T("Unable to find credential with ID: {{.credentialID}}.\n", map[string]interface{}{"credentialID": credentialID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to delete credential: {{.storageID}}.\n", map[string]interface{}{"storageID": storageID}), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Credential: {{.credentialID}} was deleted.", map[string]interface{}{"credentialID": credentialID}))
	return nil
}
