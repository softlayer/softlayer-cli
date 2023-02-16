package virtual

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CaptureCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Name                 string
	All                  bool
	Note                 string
	BlockDevices         []int
}

func NewCaptureCommand(sl *metadata.SoftlayerCommand) (cmd *CaptureCommand) {
	thisCmd := &CaptureCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capture " + T("IDENTIFIER"),
		Short: T("Capture virtual server instance into an image"),
		Long: T(`${COMMAND_NAME} sl vs capture IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs capture 12345678 -n mycloud --all --note testing
   ${COMMAND_NAME} sl vs capture 12345678 -n mycloud --device 111111 --device 222222 --note testing
   --all and --device options can not be set at the same time and one is required.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name of the image [required]"))
	cobraCmd.Flags().BoolVar(&thisCmd.All, "all", false, T("Capture all block devices that belong to the virtual server"))
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("Add a note to be associated with the image"))
	cobraCmd.Flags().IntSliceVar(&thisCmd.BlockDevices, "device", []int{}, T("The block device ID´s to archive, multiple occurrence allowed"))
	cobraCmd.MarkFlagsMutuallyExclusive("all", "device")
	return thisCmd
}
func (cmd *CaptureCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	if cmd.Name == "" {
		return slErrors.NewMissingInputError("-n|--name")
	}
	if !cmd.All && len(cmd.BlockDevices) == 0 {
		return slErrors.NewMissingInputError("--all|--device")
	}
	blockDevices := []datatypes.Virtual_Guest_Block_Device{}
	if cmd.All {
		virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, "id,blockDevices[id,device,mountType,diskImage[id,metadataFlag,type[keyName]]]")
		if err != nil {
			subs := map[string]interface{}{"VsID": vsID}
			return slErrors.NewAPIError(T("Failed to get virtual server instance: {{.VsID}}.", subs), err.Error(), 2)
		}
		blockDevices = getDisks(virtualGuest)

	} else {
		for _, blockDeviceId := range cmd.BlockDevices {
			blockDevice := datatypes.Virtual_Guest_Block_Device{
				Id: sl.Int(blockDeviceId),
			}
			blockDevices = append(blockDevices, blockDevice)
		}
	}

	outputFormat := cmd.GetOutputFlag()

	image, err := cmd.VirtualServerManager.CaptureImage(vsID, cmd.Name, cmd.Note, blockDevices)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to capture image for virtual server instance: {{.VsID}}.\n",
			map[string]interface{}{"VsID": vsID}), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add(T("Virtual guest ID"), strconv.Itoa(vsID))
	table.Add(T("Image ID"), utils.FormatIntPointer(image.Id))
	table.Add(T("Date time"), utils.FormatSLTimePointer(image.CreateDate))
	table.Add(T("Note"), utils.FormatStringPointer(image.Note))
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func getDisks(vs datatypes.Virtual_Guest) []datatypes.Virtual_Guest_Block_Device {
	disks := []datatypes.Virtual_Guest_Block_Device{}
	for _, disk := range vs.BlockDevices {
		//We never want metadata disks
		if disk.DiskImage != nil && disk.DiskImage.MetadataFlag != nil && *disk.DiskImage.MetadataFlag == true {
			continue
		}
		//We never want swap devices
		if disk.DiskImage != nil && disk.DiskImage.Type != nil && disk.DiskImage.Type.KeyName != nil && *disk.DiskImage.Type.KeyName == "SWAP" {
			continue
		}
		//We never want CD images
		if disk.MountType != nil && *disk.MountType == "CD" {
			continue
		}
		disks = append(disks, disk)
	}
	return disks
}
