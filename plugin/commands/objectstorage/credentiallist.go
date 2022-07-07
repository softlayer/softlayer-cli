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

type CredentialListCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewCredentialListCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *CredentialListCommand) {
	return &CredentialListCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func CredentialListMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "credential-list",
		Description: T("Retrieve credentials used for generating an AWS signature. Max of 2."),
		Usage:       T(`${COMMAND_NAME} sl object-storage credential-list IDENTIFIER`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *CredentialListCommand) Run(c *cli.Context) error {
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
	credentialList, err := cmd.ObjectStorageManager.ListCredential(StorageID, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list credentials.")+err.Error(), 2)
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
