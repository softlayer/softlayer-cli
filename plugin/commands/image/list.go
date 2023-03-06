package image

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
	Name         string
	Public       bool
	Private      bool
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all images on your account"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("Filter on image name"))
	cobraCmd.Flags().BoolVar(&thisCmd.Public, "public", false, T("Display only public images"))
	cobraCmd.Flags().BoolVar(&thisCmd.Private, "private", false, T("Display only private images"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

type Image struct {
	Visibility string
	datatypes.Virtual_Guest_Block_Device_Template_Group
}

func (cmd *ListCommand) Run(args []string) error {
	var publicImages, privateImages []datatypes.Virtual_Guest_Block_Device_Template_Group
	var err error
	mask := "mask[id,name,accountId,imageType.name,children[blockDevices[diskImage[softwareReferences[softwareDescription]]]]]"

	if cmd.Public && cmd.Private {
		return bmxErr.NewInvalidUsageError(T("[--public] is not allowed with [--private]."))
	}

	outputFormat := cmd.GetOutputFlag()

	if cmd.Public && !cmd.Private {
		publicImages, err = cmd.ImageManager.ListPublicImages(cmd.Name, mask)
		if err != nil {
			return bmxErr.NewAPIError(T("Failed to list public images."), err.Error(), 2)
		}
	} else if cmd.Private && !cmd.Public {
		privateImages, err = cmd.ImageManager.ListPrivateImages(cmd.Name, mask)
		if err != nil {
			return bmxErr.NewAPIError(T("Failed to list private images."), err.Error(), 2)
		}
	} else {
		publicImages, err = cmd.ImageManager.ListPublicImages(cmd.Name, mask)
		if err != nil {
			return bmxErr.NewAPIError(T("Failed to list public images."), err.Error(), 2)
		}
		privateImages, err = cmd.ImageManager.ListPrivateImages(cmd.Name, mask)
		if err != nil {
			return bmxErr.NewAPIError(T("Failed to list private images."), err.Error(), 2)
		}
	}

	allImages := []Image{}
	for _, pubImage := range publicImages {
		image := Image{T("Public"), pubImage}
		allImages = append(allImages, image)
	}
	for _, priImage := range privateImages {
		image := Image{T("Private"), priImage}
		allImages = append(allImages, image)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, allImages)
	}

	if len(allImages) == 0 {
		cmd.UI.Print(T("No image found."))
	}

	table := cmd.UI.Table([]string{T("id"), T("name"), T("type"), T("visibility"), T("account"), T("os")})

	for _, image := range allImages {
		var typeName string
		if image.ImageType != nil {
			typeName = utils.FormatStringPointer(image.ImageType.Name)
		}

		os := "-"
		if image.Children != nil && len(image.Children) != 0 {
			if image.Children[0].BlockDevices != nil && len(image.Children[0].BlockDevices) != 0 {
				for _, blockDevice := range image.Children[0].BlockDevices {
					if blockDevice.DiskImage.SoftwareReferences != nil && len(blockDevice.DiskImage.SoftwareReferences) != 0 {
						os = *blockDevice.DiskImage.SoftwareReferences[0].SoftwareDescription.LongDescription
					}
				}
			}
		}
		table.Add(utils.FormatIntPointer(image.Id),
			utils.FormatStringPointer(image.Name),
			typeName, image.Visibility,
			utils.FormatIntPointer(image.AccountId),
			os,
		)
	}
	table.Print()
	return nil
}
