package objectstorage

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialCreateCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewCredentialCreateCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *CredentialCreateCommand) {
	return &CredentialCreateCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func CredentialCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "credential-create",
		Description: T("Create credentials for an IBM Cloud Object Storage Account."),
		Usage: T(`${COMMAND_NAME} sl object-storage credential-create IDENTIFIER

Examples:
	${COMMAND_NAME} sl object-storage credential-create 123456`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *CredentialCreateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	StorageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := ""
	credentialCreate, err := cmd.ObjectStorageManager.CreateCredential(StorageID, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to create credential. ")+err.Error(), 2)
	}
	PrintCredentialCreated(credentialCreate, cmd.UI, outputFormat)
	return nil
}

func PrintCredentialCreated(credentialCreate []datatypes.Network_Storage_Credential, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Id"),
		T("Password"),
		T("Username"),
		T("Type Name"),
	})

	for _, credential := range credentialCreate {
		table.Add(
			utils.FormatIntPointer(credential.Id),
			utils.FormatStringPointerName(credential.Password),
			utils.FormatStringPointerName(credential.Username),
			utils.FormatStringPointerName(credential.Type.Name),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}
