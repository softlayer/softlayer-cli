package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ObjectStorageDetailCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewObjectStorageDetailCommand(sl *metadata.SoftlayerStorageCommand) (cmd *ObjectStorageDetailCommand) {
	thisCmd := &ObjectStorageDetailCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "object-storage-detail " + T("IDENTIFIER"),
		Short: T("Display details for a cloud object storage."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}
func (cmd *ObjectStorageDetailCommand) Run(args []string) error {

	storageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat := cmd.GetOutputFlag()

	networkStorageDetail, err := cmd.StorageManager.GetNetworkStorageDetail(storageID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get details of storage {{.StorageID}}.\n",
			map[string]interface{}{"StorageID": storageID}), err.Error(), 2)
	}

	bucket, err := cmd.StorageManager.GetBuckets(storageID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get bucket of storage {{.StorageID}}.\n",
			map[string]interface{}{"StorageID": storageID}), err.Error(), 2)
	}



	table := cmd.UI.Table([]string{
		T("Name"),
		T("Value"),
	})
	table.Add("Id", utils.FormatIntPointer(networkStorageDetail.Id))
	table.Add("Username", utils.FormatStringPointer(networkStorageDetail.Username))
	table.Add("Name Service Resource", utils.FormatStringPointer(networkStorageDetail.ServiceResource.Name))
	table.Add("Type Service Resource", utils.FormatStringPointer(networkStorageDetail.ServiceResource.Type.Type))
	table.Add("Datacenter", utils.FormatStringPointer(networkStorageDetail.ServiceResource.Datacenter.Name))
	table.Add("Storage Type", utils.FormatStringPointer(networkStorageDetail.StorageType.KeyName))
	table.Add("Bytes Used", utils.B2GB(*bucket[0].BytesUsed))
	table.Add("Bucket Name", utils.FormatStringPointer(bucket[0].Name))

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
