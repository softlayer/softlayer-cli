package image

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ShareCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewShareCommand(sl *metadata.SoftlayerCommand) (cmd *ShareCommand) {
	thisCmd := &ShareCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "share " + T("IDENTIFIER") + " " + T("ACCOUNT ID"),
		Short: T("Share an image template with another account."),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ShareCommand) Run(args []string) error {
	imageId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Image Id")
	}

	accountId, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Account Id")
	}

	image, err := cmd.ImageManager.ShareImage(imageId, accountId)
	if err != nil {
		return errors.NewAPIError(T("Failed to share image: {{.ImageId}} with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": accountId}), err.Error(), 2)
	}
	if image {
		cmd.UI.Ok()
		cmd.UI.Print(T("Image {{.ImageId}} was shared with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": accountId}))
	}
	return nil
}
