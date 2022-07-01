package block

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ObjectListCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewObjectListCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *ObjectListCommand) {
	return &ObjectListCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *ObjectListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	cloudObjectStorages, err := cmd.StorageManager.GetHubNetworkStorage("")
	if err != nil {
		return cli.NewExitError(T("Failed to get Cloud Object Storages.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Account name"), T("Description"), T("Create Date"), T("Type")})
	for _, objectStorage := range cloudObjectStorages {
		table.Add(
			utils.FormatIntPointer(objectStorage.Id),
			utils.FormatStringPointer(objectStorage.Username),
			utils.FormatStringPointer(objectStorage.StorageType.Description),
			utils.FormatSLTimePointer(objectStorage.BillingItem.CreateDate),
			utils.FormatStringPointer(objectStorage.StorageType.KeyName),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func BlockObjectListMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "object-list",
		Description: T("List cloud block storage."),
		Usage: T(`${COMMAND_NAME} sl block object-list [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block object-list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
