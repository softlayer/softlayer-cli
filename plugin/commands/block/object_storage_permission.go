package block

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ObjectStoragePermissionCommand struct {
	*metadata.SoftlayerStorageCommand
	Command              *cobra.Command
	StorageManager       managers.StorageManager
	ObjectStorageManager managers.ObjectStorageManager
}

func NewObjectStoragePermissionCommand(sl *metadata.SoftlayerStorageCommand) (cmd *ObjectStoragePermissionCommand) {
	thisCmd := &ObjectStoragePermissionCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
		ObjectStorageManager:    managers.NewObjectStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "object-storage-permission" + " " + T("IDENTIFIER"),
		Short: T("Display permission details for a cloud object storage."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}
func (cmd *ObjectStoragePermissionCommand) Run(args []string) error {

	storageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	outputFormat := cmd.GetOutputFlag()

	cloudObjectStorages, err := cmd.StorageManager.GetNetworkMessageDeliveryAccounts(storageID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get permissions."), err.Error(), 2)
	}

	endPoints, err := cmd.ObjectStorageManager.GetEndpoints(storageID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get endPoints."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{
		T("Name"),
		T("Value"),
	})

	table.Add("UUID", utils.FormatStringPointer(cloudObjectStorages.Uuid))

	bufTableCredentials := new(bytes.Buffer)
	tableCredentials := terminal.NewTable(bufTableCredentials, []string{
		T("Id"),
		T("Access Key ID"),
		T("Secret Access Key"),
		T("Description"),
	})

	for _, credential := range cloudObjectStorages.Credentials {
		tableCredentials.Add(
			utils.FormatIntPointer(credential.Id),
			utils.FormatStringPointer(credential.Username),
			utils.FormatStringPointer(credential.Password),
			utils.FormatStringPointer(credential.Type.Description),
		)
	}

	utils.PrintTable(cmd.UI, tableCredentials, outputFormat)
	table.Add("Credentials", bufTableCredentials.String())

	bufTableEndPoints := new(bytes.Buffer)
	tableEndPoints := terminal.NewTable(bufTableEndPoints, []string{
		T("Region"),
		T("Location"),
		T("Type"),
		T("URL"),
	})

	for _, endPoint := range endPoints {
		tableEndPoints.Add(
			utils.FormatStringPointer(endPoint.Region),
			utils.FormatStringPointer(endPoint.Location),
			utils.FormatStringPointer(endPoint.Type),
			utils.FormatStringPointer(endPoint.Url),
		)
	}

	utils.PrintTable(cmd.UI, tableEndPoints, outputFormat)
	table.Add("EndPoint URL´s", bufTableEndPoints.String())

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
