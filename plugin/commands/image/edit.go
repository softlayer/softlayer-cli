package image

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
	Name         string
	Note         string
	Tag          string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit details of an image"),
		Long: T(`${COMMAND_NAME} sl image edit IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl image edit 12345678 --name ubuntu16 --note testing --tag staging
   This command edits an image with ID 12345678 and set its name to "ubuntu16", note to "testing", and tag to "staging".`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("Name of the image"))
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("Add notes for the image"))
	cobraCmd.Flags().StringVar(&thisCmd.Tag, "tag", "", T("Tags for the image"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	if cmd.Name == "" && cmd.Note == "" && cmd.Tag == "" {
		return bmxErr.NewInvalidUsageError(T("One of --name, --note and --tag must be specified."))
	}

	succeeses, messages := cmd.ImageManager.EditImage(imageID, cmd.Name, cmd.Note, cmd.Tag)
	for index, succees := range succeeses {
		if succees == true {
			cmd.UI.Ok()
			cmd.UI.Print(messages[index])
		} else {
			cmd.UI.Print(terminal.FailureColor("FAILED"))
			cmd.UI.Print(messages[index])
		}
	}
	return nil
}
