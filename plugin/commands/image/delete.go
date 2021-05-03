package image

import (
	"strconv"
	"strings"

	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
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
