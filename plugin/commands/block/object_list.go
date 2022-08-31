package block

import (
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ObjectListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewObjectListCommand(sl *metadata.SoftlayerStorageCommand) (cmd *ObjectListCommand) {
	thisCmd := &ObjectListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-create",
		Short: T("List cloud block storage."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}
func (cmd *ObjectListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	cloudObjectStorages, err := cmd.StorageManager.GetHubNetworkStorage("")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get Cloud Object Storages.\n"), err.Error(), 2)
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
