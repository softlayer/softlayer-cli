package objectstorage

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AccountsCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewAccountsCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *AccountsCommand) {
	return &AccountsCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func AccountsMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "accounts",
		Description: T("List Object Storage accounts."),
		Usage:       T(`${COMMAND_NAME} sl object-storage accounts`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "limit",
				Usage: T("Result limit [Default: 50]"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *AccountsCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	limit := 50
	if c.IsSet("limit") {
		limit = c.Int("limit")
	}

	mask := "mask[id,username,notes,vendorName,serviceResource]"
	accounts, err := cmd.ObjectStorageManager.GetAccounts(mask, limit)
	if err != nil {
		return cli.NewExitError(T("Failed to get account’s associated Virtual Storage volumes.")+err.Error(), 2)
	}
	PrintAccounts(accounts, cmd.UI, outputFormat)
	return nil
}

func PrintAccounts(accounts []datatypes.Network_Storage, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Id"),
		T("Name"),
		T("ApiType"),
	})

	for _, account := range accounts {
		apiType := ""
		if account.VendorName != nil && strings.Contains(utils.FormatStringPointerName(account.VendorName), "Swift") {
			apiType = "Swift"
		} else {
			if strings.Contains(utils.FormatStringPointerName(account.ServiceResource.Name), "Cleversafe") {
				apiType = "S3"
			}
		}
		table.Add(
			utils.FormatIntPointer(account.Id),
			utils.FormatStringPointerName(account.Username),
			utils.FormatStringPointerName(&apiType),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}
