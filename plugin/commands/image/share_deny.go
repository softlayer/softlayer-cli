package image

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ShareDenyCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewShareDenyCommand(sl *metadata.SoftlayerCommand) (cmd *ShareDenyCommand) {
	thisCmd := &ShareDenyCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "share-deny " + T("IDENTIFIER") + " " + T("ACCOUNT ID"),
		Short: T("Deny share an image template with another account."),
		Args:  metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ShareDenyCommand) Run(args []string) error {
	imageId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Image Id")
	}

	accountId, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Account Id")
	}

	image, err := cmd.ImageManager.ShareDenyImage(imageId, accountId)
	if err != nil {
		return errors.NewAPIError(T("Failed to deny share image: {{.ImageId}} with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": accountId}), err.Error(), 2)
	}
	if image {
		cmd.UI.Ok()
		cmd.UI.Print(T("Image {{.ImageId}} was deny shared with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": accountId}))
	}
	return nil
}
