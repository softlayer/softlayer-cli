package hardware

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type StorageCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewStorageCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *StorageCommand) {
	return &StorageCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *StorageCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	iscsiStorageData, err := cmd.HardwareManager.GetStorageDetails(hardwareID, "ISCSI")
	if err != nil {
		return cli.NewExitError(T("Failed to get iscsi storage detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	nasStorageData, err := cmd.HardwareManager.GetStorageDetails(hardwareID, "NAS")
	if err != nil {
		return cli.NewExitError(T("Failed to get nas storage detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	storageCredentials, err := cmd.HardwareManager.GetStorageCredentials(hardwareID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the storage credential detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	hardDrives, err := cmd.HardwareManager.GetHardDrives(hardwareID)
	if err != nil {
		return cli.NewExitError(T("Failed to get the hard drives detail for the hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID})+err.Error(), 2)
	}

	var storageDetailList []interface{}
	storageDetailList = append(storageDetailList, storageCredentials)
	storageDetailList = append(storageDetailList, iscsiStorageData)
	storageDetailList = append(storageDetailList, nasStorageData)
	storageDetailList = append(storageDetailList, hardDrives)

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
			*iscsi.AllowedHardware[0].Datacenter.LongName,
			*iscsi.Notes)
	}
	tableIscsi.Print()

	cmd.UI.Print("\nFile Storage Details")
	tableNas := cmd.UI.Table([]string{T("Volume name"), T("capacity"), T("Host name"), T("Location"), T("Notes")})
	for _, nas := range nasStorageData {
		tableNas.Add(
			*nas.Username,
			utils.FormatIntPointer(nas.CapacityGb),
			*nas.ServiceResourceBackendIpAddress,
			*nas.AllowedHardware[0].Datacenter.LongName,
			*nas.Notes)
	}
	tableNas.Print()

	cmd.UI.Print("\nOther storage details")
	tableHardDrives := cmd.UI.Table([]string{T("Type"), T("Name"), T("Capacity"), T("Serial #")})
	for _, drive := range hardDrives {
		typeDrive := drive.HardwareComponentModel.HardwareGenericComponentModel.HardwareComponentType.Type
		name := fmt.Sprintf("%s %s", *drive.HardwareComponentModel.Manufacturer, *drive.HardwareComponentModel.Name)
		capacity := fmt.Sprintf("%.2f %s", *drive.HardwareComponentModel.HardwareGenericComponentModel.Capacity, *drive.HardwareComponentModel.HardwareGenericComponentModel.Units)
		serial := drive.SerialNumber

		tableHardDrives.Add(*typeDrive, name, capacity, *serial)
	}
	tableHardDrives.Print()

	return nil
}
