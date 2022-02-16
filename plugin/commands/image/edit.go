package image

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewEditCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	imageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	if c.NumFlags() == 0 {
		return bmxErr.NewInvalidUsageError(T("One of --name, --note and --tag must be specified."))
	}

	succeeses, messages := cmd.ImageManager.EditImage(imageID, c.String("name"), c.String("note"), c.String("tag"))
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

func ImageEditMetaData() cli.Command {
	return cli.Command{
		Category:    "image",
		Name:        "edit",
		Description: T("Edit details of an image"),
		Usage: T(`${COMMAND_NAME} sl image edit IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl image edit 12345678 --name ubuntu16 --note testing --tag staging
   This command edits an image with ID 12345678 and set its name to "ubuntu16", note to "testing", and tag to "staging".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("Name of the image"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Add notes for the image"),
			},
			cli.StringFlag{
				Name:  "tag",
				Usage: T("Tags for the image"),
			},
		},
	}
}
