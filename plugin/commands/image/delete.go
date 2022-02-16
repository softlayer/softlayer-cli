package image

import (
	"strconv"
	"strings"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type DeleteCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewDeleteCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *DeleteCommand) {
	return &DeleteCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}

func (cmd *DeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	imageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	err = cmd.ImageManager.DeleteImage(imageID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find image with ID {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to delete image: {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID})+err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Image {{.ImageID}} was deleted.", map[string]interface{}{"ImageID": imageID}))
	return nil
}

func ImageDelMetaData() cli.Command {
	return cli.Command{
		Category:    "image",
		Name:        "delete",
		Description: T("Delete an image "),
		Usage: T(`${COMMAND_NAME} sl image delete IDENTIFIER

EXAMPLE: 
   ${COMMAND_NAME} sl image delete 12345678
   This command deletes image with ID 12345678.`),
	}
}
