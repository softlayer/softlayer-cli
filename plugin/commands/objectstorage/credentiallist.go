package objectstorage

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialListCommand struct {
	*metadata.SoftlayerCommand
	ObjectStorageManager managers.ObjectStorageManager
	Command              *cobra.Command
}

func NewCredentialListCommand(sl *metadata.SoftlayerCommand) *CredentialListCommand {
	thisCmd := &CredentialListCommand{
		SoftlayerCommand:     sl,
		ObjectStorageManager: managers.NewObjectStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "credential-list",
		Short: T("Retrieve credentials used for generating an AWS signature. Max of 2."),
		Long:  T(`${COMMAND_NAME} sl object-storage credential-list IDENTIFIER [OPTIONS]`),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialListCommand) Run(args []string) error {
	StorageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := ""
	credentialList, err := cmd.ObjectStorageManager.ListCredential(StorageID, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to list credentials. "), err.Error(), 2)
	}
	PrintCredentialList(credentialList, cmd.UI, outputFormat)
	return nil
}

func PrintCredentialList(credentialList []datatypes.Network_Storage_Credential, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Id"),
		T("Password"),
		T("Username"),
		T("Type Name"),
	})

	for _, credential := range credentialList {
		table.Add(
			utils.FormatIntPointer(credential.Id),
			utils.FormatStringPointerName(credential.Password),
			utils.FormatStringPointerName(credential.Username),
			utils.FormatStringPointerName(credential.Type.Name),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}
