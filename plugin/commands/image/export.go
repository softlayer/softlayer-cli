package image

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ExportCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewExportCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *ExportCommand) {
	return &ExportCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}

func (cmd *ExportCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	imageId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Image ID")
	}
	uri := c.Args()[1]
	ibmApiKey := c.Args()[2]

	config := datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration{
		Uri:       &uri,
		IbmApiKey: &ibmApiKey,
	}

	_, err = cmd.ImageManager.ExportImage(imageId, config)
	if err != nil {
		return err
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("The image {{.ImageId}} was exported successfully!", map[string]interface{}{"ImageId": imageId}))
	return nil
}

func ImageExportMetaData() cli.Command {
	return cli.Command{
		Category:    "image",
		Name:        "export",
		Description: T("Export an image to an object storage"),
		Usage:       T("${COMMAND_NAME} sl image export IDENTIFIER URI API_KEY\n  IDENTIFIER: ID of the image\n  URI: The URI for an object storage object (.vhd/.iso file) of the format: cos://<regionName>/<bucketName>/<objectPath>\n  API_KEY: The IBM Cloud API Key with access to IBM Cloud Object Storage instance."),
	}
}
