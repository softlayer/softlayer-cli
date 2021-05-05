package image

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
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
