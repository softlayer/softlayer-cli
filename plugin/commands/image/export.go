package image

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ExportCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewExportCommand(sl *metadata.SoftlayerCommand) (cmd *ExportCommand) {
	thisCmd := &ExportCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "export " + T("IDENTIFIER") + " " + T("URI") + " " + T("API_KEY"),
		Short: T("Export an image to an object storage"),
		Long: T(`
EXAMPLE:
	${COMMAND_NAME} sl image export IDENTIFIER URI API_KEY
	IDENTIFIER: ID of the image
	URI: The URI for an object storage object (.vhd/.iso file) of the format: cos://<regionName>/<bucketName>/<objectPath>
	API_KEY: The IBM Cloud API Key with access to IBM Cloud Object Storage instance.`),
		Args: metadata.ThreeArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ExportCommand) Run(args []string) error {
	imageId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Image ID")
	}
	uri := args[1]
	ibmApiKey := args[2]

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
