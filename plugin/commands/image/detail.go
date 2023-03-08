package image

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details for an image"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	outputFormat := cmd.GetOutputFlag()

	image, err := cmd.ImageManager.GetImage(imageID)
	if err != nil {
		return bmxErr.NewAPIError(T("Failed to get image: {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, image)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(image.Id))
	table.Add(T("global_identifier"), utils.FormatStringPointer(image.GlobalIdentifier))
	table.Add(T("name"), utils.FormatStringPointer(image.Name))
	if image.Status != nil {
		table.Add(T("status"), utils.FormatStringPointer(image.Status.Name))
	} else {
		table.Add(T("status"), "-")
	}

	table.Add(T("account"), utils.FormatIntPointer(image.AccountId))
	table.Add(T("created"), utils.FormatSLTimePointer(image.CreateDate))

	diskSpace := "-"
	if image.FirstChild != nil {
		diskSpace = utils.B2GB(int(*image.FirstChild.BlockDevicesDiskSpaceTotal))
	}
	table.Add(T("disk space"), diskSpace)

	if image.PublicFlag != nil && *image.PublicFlag == 1 {
		table.Add(T("visibility"), T("Public"))
	} else {
		table.Add(T("visibility"), T("Private"))
	}
	if image.ImageType != nil {
		table.Add(T("type"), utils.FormatStringPointer(image.ImageType.Name))
	} else {
		table.Add(T("type"), "-")
	}

	table.Add(T("flex"), utils.FormatBoolPointer(image.FlexImageFlag))
	table.Add(T("note"), utils.FormatStringPointer(image.Note))

	os := "-"
	if image.Children[0].BlockDevices != nil && len(image.Children[0].BlockDevices) != 0 {
		for _, blockDevice := range image.Children[0].BlockDevices {
			if blockDevice.DiskImage.SoftwareReferences != nil && len(blockDevice.DiskImage.SoftwareReferences) != 0 {
				os = *blockDevice.DiskImage.SoftwareReferences[0].SoftwareDescription.LongDescription
			}
		}
	}
	table.Add(T("os"), os)
	table.Add(T("datacenters"), getDatacenters(image.Children, cmd.UI, outputFormat))
	table.Add(T("virtual disks"), getVirtualDisks(image.Children, cmd.UI, outputFormat))
	table.Add(T("share image"), getShareImages(image, cmd.UI, outputFormat))
	table.Print()
	return nil
}

func getDatacenters(childrens []datatypes.Virtual_Guest_Block_Device_Template_Group, ui terminal.UI, outputFormat string) string {
	bufTable := new(bytes.Buffer)
	table := terminal.NewTable(bufTable, []string{
		T("Data Center"),
		T("Size"),
	})
	for _, child := range childrens {
		table.Add(
			utils.FormatStringPointer(child.Datacenter.Name),
			utils.B2GB(int(*child.BlockDevicesDiskSpaceTotal)),
		)
	}
	utils.PrintTable(ui, table, outputFormat)
	return bufTable.String()
}

func getVirtualDisks(childrens []datatypes.Virtual_Guest_Block_Device_Template_Group, ui terminal.UI, outputFormat string) string {
	bufTable := new(bytes.Buffer)
	table := terminal.NewTable(bufTable, []string{
		T("Device"),
		T("Capacity"),
		T("Size on disk"),
	})
	for _, blockDevices := range childrens[0].BlockDevices {
		device_name := utils.FormatStringPointer(blockDevices.DiskImage.Name)
		if len(blockDevices.DiskImage.SoftwareReferences) > 0 {
			device_name = *blockDevices.DiskImage.SoftwareReferences[0].SoftwareDescription.LongDescription
		}
		sizeOnDisk := "N/A"
		if blockDevices.DiskSpace != nil {
			sizeOnDisk = utils.B2GB(int(*blockDevices.DiskSpace))
		}

		table.Add(
			device_name,
			utils.FormatIntPointer(blockDevices.DiskImage.Capacity)+
				utils.FormatStringPointer(blockDevices.DiskImage.Units),
			sizeOnDisk,
		)
	}
	utils.PrintTable(ui, table, outputFormat)
	return bufTable.String()
}

func getShareImages(image datatypes.Virtual_Guest_Block_Device_Template_Group, ui terminal.UI, outputFormat string) string {
	bufTable := new(bytes.Buffer)
	table := terminal.NewTable(bufTable, []string{
		T("Account"),
		T("Shared On"),
	})
	for _, account := range image.AccountReferences {
		if *account.AccountId != *image.AccountId {
			table.Add(
				utils.FormatIntPointer(account.AccountId),
				utils.FormatSLTimePointer(account.CreateDate),
			)
		}

	}
	utils.PrintTable(ui, table, outputFormat)
	return bufTable.String()
}
