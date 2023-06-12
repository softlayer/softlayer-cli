package hardware

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type StorageCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewStorageCommand(sl *metadata.SoftlayerCommand) (cmd *StorageCommand) {
	thisCmd := &StorageCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "storage " + T("IDENTIFIER"),
		Short: T("Get storage details for a hardware server."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *StorageCommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	iscsiStorageData, err := cmd.HardwareManager.GetStorageDetails(hardwareID, "ISCSI")
	if err != nil {
		return errors.NewAPIError(T("Failed to get iscsi storage detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	nasStorageData, err := cmd.HardwareManager.GetStorageDetails(hardwareID, "NAS")
	if err != nil {
		return errors.NewAPIError(T("Failed to get nas storage detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	storageCredentials, err := cmd.HardwareManager.GetStorageCredentials(hardwareID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get the storage credential detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	hardDrives, err := cmd.HardwareManager.GetHardDrives(hardwareID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get the hard drives detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
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
		tableIscsi.Add(
			*iscsi.Username,
			utils.FormatIntPointer(iscsi.CapacityGb),
			*iscsi.ServiceResourceBackendIpAddress,
			*iscsi.AllowedHardware[0].Datacenter.LongName,
			*iscsi.Notes)
	}
	utils.PrintTable(cmd.UI, tableIscsi, outputFormat)

	cmd.UI.Print("\nFile Storage Details")
	tableNas := cmd.UI.Table([]string{T("Volume name"), T("capacity"), T("Hostname"), T("Location"), T("Notes")})
	for _, nas := range nasStorageData {
		tableNas.Add(
			*nas.Username,
			utils.FormatIntPointer(nas.CapacityGb),
			*nas.ServiceResourceBackendIpAddress,
			*nas.AllowedHardware[0].Datacenter.LongName,
			*nas.Notes)
	}
	utils.PrintTable(cmd.UI, tableNas, outputFormat)

	cmd.UI.Print("\nOther storage details")
	tableHardDrives := cmd.UI.Table([]string{T("Type"), T("Name"), T("Capacity"), T("Serial #")})
	for _, drive := range hardDrives {
		typeDrive := drive.HardwareComponentModel.HardwareGenericComponentModel.HardwareComponentType.Type
		name := fmt.Sprintf("%s %s", *drive.HardwareComponentModel.Manufacturer, *drive.HardwareComponentModel.Name)
		capacity := fmt.Sprintf("%.2f %s", *drive.HardwareComponentModel.HardwareGenericComponentModel.Capacity, *drive.HardwareComponentModel.HardwareGenericComponentModel.Units)
		serial := drive.SerialNumber

		tableHardDrives.Add(*typeDrive, name, capacity, *serial)
	}
	utils.PrintTable(cmd.UI, tableHardDrives, outputFormat)

	return nil
}
