package virtual

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/datatypes"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type StorageCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewStorageCommand(sl *metadata.SoftlayerCommand) (cmd *StorageCommand) {
	thisCmd := &StorageCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "storage " + T("IDENTIFIER"),
		Short: T("Get storage details for a virtual server."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *StorageCommand) Run(args []string) error {

	vsID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()
	subs := map[string]interface{}{"VsId": vsID, "VsID": vsID, "ID": vsID}

	iscsiStorageData, err := cmd.VirtualServerManager.GetStorageDetails(vsID, "ISCSI")
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get iscsi storage detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	nasStorageData, err := cmd.VirtualServerManager.GetStorageDetails(vsID, "NAS")
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get nas storage detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	storageCredentials, err := cmd.VirtualServerManager.GetStorageCredentials(vsID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get the storage credential detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	portableStorage, err := cmd.VirtualServerManager.GetPortableStorage(vsID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get the portable storage detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	localDisks, err := cmd.VirtualServerManager.GetLocalDisks(vsID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get the local disks detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	cmd.UI.Print("Block Storage Details\niSCSI")
	tableCredentials := cmd.UI.Table([]string{T("Username"), T("Password"), T("IQN")})
	if storageCredentials.Credential != nil && storageCredentials.Credential.Password != nil {
		tableCredentials.Add(
			*storageCredentials.Credential.Username,
			*storageCredentials.Credential.Password,
			*storageCredentials.Name)
	}
	utils.PrintTable(cmd.UI, tableCredentials, outputFormat)

	tableIscsi := cmd.UI.Table([]string{T("\nLUN name"), T("capacity"), T("Target address"), T("Location"), T("Notes")})
	for _, iscsi := range iscsiStorageData {
		notes := ""
		if iscsi.Notes != nil {
			notes = *iscsi.Notes
		}
		tableIscsi.Add(
			*iscsi.Username,
			utils.FormatIntPointer(iscsi.CapacityGb),
			*iscsi.ServiceResourceBackendIpAddress,
			*iscsi.AllowedVirtualGuests[0].Datacenter.LongName,
			notes)
	}
	utils.PrintTable(cmd.UI, tableIscsi, outputFormat)

	cmd.UI.Print("\nPortable Storage")
	tablePortableStorage := cmd.UI.Table([]string{T("Description"), T("Capacity"), T("Location")})
	for _, portable := range portableStorage {
		tablePortableStorage.Add(
			*portable.Description,
			utils.FormatIntPointer(portable.Capacity),
			*portable.BillingItem.Location.LongName)
	}
	utils.PrintTable(cmd.UI, tablePortableStorage, outputFormat)

	cmd.UI.Print("\nFile Storage Details")
	tableNas := cmd.UI.Table([]string{T("Volume name"), T("capacity"), T("Hostname"), T("Location"), T("Notes")})
	for _, nas := range nasStorageData {
		notes := ""
		if nas.Notes != nil {
			notes = *nas.Notes
		}
		tableNas.Add(
			*nas.Username,
			utils.FormatIntPointer(nas.CapacityGb),
			*nas.ServiceResourceBackendIpAddress,
			*nas.AllowedVirtualGuests[0].Datacenter.LongName,
			notes)
	}
	utils.PrintTable(cmd.UI, tableNas, outputFormat)

	cmd.UI.Print("\nSystem storage details")
	tableLocalDisks := cmd.UI.Table([]string{T("Type"), T("Name"), T("Drive"), T("Capacity")})
	for _, disk := range localDisks {
		if disk.DiskImage != nil {
			capacity := fmt.Sprintf("%d %s", *disk.DiskImage.Capacity, *disk.DiskImage.Units)

			tableLocalDisks.Add(cmd.getLocalType(disk), *disk.MountType, *disk.Device, capacity)
		}
	}
	utils.PrintTable(cmd.UI, tableLocalDisks, outputFormat)

	return nil
}

// Returns the virtual server local disk type.
// param disks: virtual server local disks.
func (cmd *StorageCommand) getLocalType(disk datatypes.Virtual_Guest_Block_Device) string {
	diskType := "System"
	swapType := disk.DiskImage.Description
	if strings.Contains(*swapType, "SWAP") {
		diskType = "Swap"
	}
	return diskType
}
