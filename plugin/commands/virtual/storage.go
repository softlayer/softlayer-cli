package virtual

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type StorageCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewStorageCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *StorageCommand) {
	return &StorageCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *StorageCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	iscsiStorageData, err := cmd.VirtualServerManager.GetStorageDetails(vsID, "ISCSI")
	if err != nil {
		return cli.NewExitError(T("Failed to get iscsi storage detail for the virtual server {{.ID}}.\n", map[string]interface{}{"ID": vsID})+err.Error(), 2)
	}

	nasStorageData, err := cmd.VirtualServerManager.GetStorageDetails(vsID, "NAS")
	if err != nil {
		return cli.NewExitError(T("Failed to get nas storage detail for the virtual server {{.ID}}.\n", map[string]interface{}{"ID": vsID})+err.Error(), 2)
	}

	storageCredentials, err := cmd.VirtualServerManager.GetStorageCredentials(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the storage credential detail for the virtual server {{.ID}}.\n", map[string]interface{}{"ID": vsID})+err.Error(), 2)
	}

	portableStorage, err := cmd.VirtualServerManager.GetPortableStorage(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the portable storage detail for the virtual server {{.ID}}.\n", map[string]interface{}{"ID": vsID})+err.Error(), 2)
	}

	localDisks, err := cmd.VirtualServerManager.GetLocalDisks(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the local disks detail for the virtual server {{.ID}}.\n", map[string]interface{}{"ID": vsID})+err.Error(), 2)
	}

	var storageDetailList []interface{}
	storageDetailList = append(storageDetailList, storageCredentials)
	storageDetailList = append(storageDetailList, iscsiStorageData)
	storageDetailList = append(storageDetailList, nasStorageData)
	storageDetailList = append(storageDetailList, portableStorage)
	storageDetailList = append(storageDetailList, localDisks)

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSONList(cmd.UI, storageDetailList)
	}

	cmd.UI.Print("Block Storage Details\niSCSI")
	tableCredentials := cmd.UI.Table([]string{T("Username"), T("Password"), T("IQN")})
	if storageCredentials.Credential != nil && storageCredentials.Credential.Password != nil {
		tableCredentials.Add(
			*storageCredentials.Credential.Username,
			*storageCredentials.Credential.Password,
			*storageCredentials.Name)
	}
	tableCredentials.Print()

	tableIscsi := cmd.UI.Table([]string{T("\nLUN name"), T("capacity"), T("Target address"), T("Location"), T("Notes")})
	for _, iscsi := range iscsiStorageData {
		tableIscsi.Add(
			*iscsi.Username,
			utils.FormatIntPointer(iscsi.CapacityGb),
			*iscsi.ServiceResourceBackendIpAddress,
			*iscsi.AllowedVirtualGuests[0].Datacenter.LongName,
			*iscsi.Notes)
	}
	tableIscsi.Print()

	cmd.UI.Print("\nPortable Storage")
	tablePortableStorage := cmd.UI.Table([]string{T("Description"), T("Capacity"), T("Location")})
	for _, portable := range portableStorage {
		tablePortableStorage.Add(
			*portable.Description,
			utils.FormatIntPointer(portable.Capacity),
			*portable.BillingItem.Location.LongName)
	}
	tablePortableStorage.Print()

	cmd.UI.Print("\nFile Storage Details")
	tableNas := cmd.UI.Table([]string{T("Volume name"), T("capacity"), T("Host name"), T("Location"), T("Notes")})
	for _, nas := range nasStorageData {
		tableNas.Add(
			*nas.Username,
			utils.FormatIntPointer(nas.CapacityGb),
			*nas.ServiceResourceBackendIpAddress,
			*nas.AllowedVirtualGuests[0].Datacenter.LongName,
			*nas.Notes)
	}
	tableNas.Print()

	cmd.UI.Print("\nSystem storage details")
	tableLocalDisks := cmd.UI.Table([]string{T("Type"), T("Name"), T("Drive"), T("Capacity")})
	for _, disk := range localDisks {
		if disk.DiskImage != nil {
			capacity := fmt.Sprintf("%d %s", *disk.DiskImage.Capacity, *disk.DiskImage.Units)

			tableLocalDisks.Add(cmd.getLocalType(disk), *disk.MountType, *disk.Device, capacity)
		}
	}
	tableLocalDisks.Print()

	return nil
}

//Returns the virtual server local disk type.
//param disks: virtual server local disks.
func (cmd *StorageCommand) getLocalType(disk datatypes.Virtual_Guest_Block_Device) string {
	diskType := "System"
	swapType := disk.DiskImage.Description
	if strings.Contains(*swapType, "SWAP") {
		diskType = "Swap"
	}
	return diskType
}

func VSStorageMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "storage",
		Description: T("Get storage details for a virtual server."),
		Usage: T(`${COMMAND_NAME} sl vs storage [OPTIONS] IDENTIFIER
	
EXAMPLE:
   ${COMMAND_NAME} sl vs storage 1234567
   Get storage details for a virtual server.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
