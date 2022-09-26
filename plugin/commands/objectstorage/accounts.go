package objectstorage

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	SWITF      = "Swift"
	CLEVERSAFE = "Cleversafe"
)

type AccountsCommand struct {
	*metadata.SoftlayerCommand
	ObjectStorageManager managers.ObjectStorageManager
	Command              *cobra.Command
}

func NewAccountsCommand(sl *metadata.SoftlayerCommand) *AccountsCommand {
	thisCmd := &AccountsCommand{
		SoftlayerCommand:     sl,
		ObjectStorageManager: managers.NewObjectStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "accounts",
		Short: T("List Object Storage accounts."),
		Long:  T("${COMMAND_NAME} sl object-storage accounts"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AccountsCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := ""
	accounts, err := cmd.ObjectStorageManager.GetAccounts(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get account’s associated Virtual Storage volumes."), err.Error(), 2)
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
		if account.VendorName != nil && strings.Contains(utils.FormatStringPointerName(account.VendorName), SWITF) {
			apiType = SWITF
		}
		if strings.Contains(utils.FormatStringPointerName(account.ServiceResource.Name), CLEVERSAFE) {
			apiType = "S3"
		}

		table.Add(
			utils.FormatIntPointer(account.Id),
			utils.FormatStringPointerName(account.Username),
			utils.FormatStringPointerName(&apiType),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}
