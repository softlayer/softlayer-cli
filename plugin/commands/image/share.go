package image

import (
	"strconv"

	"github.com/spf13/cobra"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ShareCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
	AccountId    int
}

func NewShareCommand(sl *metadata.SoftlayerCommand) (cmd *ShareCommand) {
	thisCmd := &ShareCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "share " + T("IDENTIFIER"),
		Short: T("Permit share an image template to another account."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.AccountId, "account-id", 0, T("Account Id for another account to share image template."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ShareCommand) Run(args []string) error {
	imageId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image Id")
	}

	if cmd.AccountId == 0 {
		return slErrors.NewMissingInputError(T("--account-id"))
	}

	image, err := cmd.ImageManager.ShareImage(imageId, cmd.AccountId)
	if err != nil {
		return bmxErr.NewAPIError(T("Failed to share image: {{.ImageId}} with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": cmd.AccountId}), err.Error(), 2)
	}
	if image {
		cmd.UI.Ok()
		cmd.UI.Print(T("Image {{.ImageId}} was shared with account {{.AccountId}}.", map[string]interface{}{"ImageId": imageId, "AccountId": cmd.AccountId}))
	}
	return nil
}
