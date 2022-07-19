package objectstorage

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialLimitCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewCredentialLimitCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *CredentialLimitCommand) {
	return &CredentialLimitCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func CredentialLimitMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "credential-limit",
		Description: T("Credential limits for this IBM Cloud Object Storage account."),
		Usage: T(`${COMMAND_NAME} sl object-storage credential-limit IDENTIFIER [OPTIONS]

Examples:
	${COMMAND_NAME} sl object-storage credential-limit 123456`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *CredentialLimitCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	storageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	credentialLimit, err := cmd.ObjectStorageManager.LimitCredential(storageID)
	if err != nil {
		return cli.NewExitError(T("Failed to get credential limit. ")+err.Error(), 2)
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
