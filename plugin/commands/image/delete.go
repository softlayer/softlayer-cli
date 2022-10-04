package image

import (
	"strconv"
	"strings"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "delete " + T("IDENTIFIER"),
		Short: T("Delete an image "),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {
	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	err = cmd.ImageManager.DeleteImage(imageID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return bmxErr.NewAPIError(T("Unable to find image with ID {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID}), err.Error(), 0)
		}
		return bmxErr.NewAPIError(T("Failed to delete image: {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID}), err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Image {{.ImageID}} was deleted.", map[string]interface{}{"ImageID": imageID}))
	return nil
}
