package objectstorage

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialLimitCommand struct {
	*metadata.SoftlayerCommand
	ObjectStorageManager managers.ObjectStorageManager
	Command              *cobra.Command
}

func NewCredentialLimitCommand(sl *metadata.SoftlayerCommand) *CredentialLimitCommand {
	thisCmd := &CredentialLimitCommand{
		SoftlayerCommand:     sl,
		ObjectStorageManager: managers.NewObjectStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "credential-limit",
		Short: T("Credential limits for this IBM Cloud Object Storage account."),
		Long: T(`${COMMAND_NAME} sl object-storage credential-limit IDENTIFIER [OPTIONS]

Examples:
	${COMMAND_NAME} sl object-storage credential-limit 123456`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialLimitCommand) Run(args []string) error {
	storageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat := cmd.GetOutputFlag()

	credentialLimit, err := cmd.ObjectStorageManager.LimitCredential(storageID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get credential limit. "), err.Error(), 2)
	}
	PrintCredentialLimit(credentialLimit, cmd.UI, outputFormat)
	return nil
}

func PrintCredentialLimit(credentialLimit int, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Limit"),
	})

	table.Add(
		utils.FormatIntPointer(&credentialLimit),
	)
	utils.PrintTable(ui, table, outputFormat)
}
