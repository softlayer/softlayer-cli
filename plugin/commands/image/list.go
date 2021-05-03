package image

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewListCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}

type Image struct {
	Visibility string
	datatypes.Virtual_Guest_Block_Device_Template_Group
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	var publicImages, privateImages []datatypes.Virtual_Guest_Block_Device_Template_Group
	var err error
	mask := "mask[id,name,accountId,imageType.name]"

	if c.IsSet("public") && c.IsSet("private") {
		return bmxErr.NewInvalidUsageError(T("[--public] is not allowed with [--private]."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if c.IsSet("public") && !c.IsSet("private") {
		publicImages, err = cmd.ImageManager.ListPublicImages(c.String("name"), mask)
		if err != nil {
			return cli.NewExitError(T("Failed to list public images.")+err.Error(), 2)
		}
	} else if c.IsSet("private") && !c.IsSet("public") {
		privateImages, err = cmd.ImageManager.ListPrivateImages(c.String("name"), mask)
		if err != nil {
			return cli.NewExitError(T("Failed to list private images.")+err.Error(), 2)
		}
	} else {
		publicImages, err = cmd.ImageManager.ListPublicImages(c.String("name"), mask)
		if err != nil {
			return cli.NewExitError(T("Failed to list public images.")+err.Error(), 2)
		}
		privateImages, err = cmd.ImageManager.ListPrivateImages(c.String("name"), mask)
		if err != nil {
			return cli.NewExitError(T("Failed to list private images.")+err.Error(), 2)
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

	table := cmd.UI.Table([]string{T("id"), T("name"), T("type"), T("visibility"), T("account")})

	for _, image := range allImages {
		var typeName string
		if image.ImageType != nil {
			typeName = utils.FormatStringPointer(image.ImageType.Name)
		}
		table.Add(utils.FormatIntPointer(image.Id),
			utils.FormatStringPointer(image.Name),
			typeName, image.Visibility,
			utils.FormatIntPointer(image.AccountId))
	}
	table.Print()
	return nil
}
