package image

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DatacenterCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewDatacenterCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *DatacenterCommand) {
	return &DatacenterCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}
func (cmd *DatacenterCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 && (c.String("add") == "" || c.String("remove") == "") {
		return slErr.NewInvalidUsageError(T("This command requires one indentifier, the option and the option arguments."))
	}
	if !c.IsSet("add") && !c.IsSet("remove") {
		return slErr.NewInvalidUsageError(T("This command requires --add or --remove option."))
	}

	imageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Image ID")
	}
	if c.IsSet("add") {
		datacenter := buildLocation(c.String("add"))
		_, err = cmd.ImageManager.AddLocation(imageID, datacenter)
		if err != nil {
			return err
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The location was added successfully!"))

	}
	if c.IsSet("remove") {
		datacenter := buildLocation(c.String("remove"))
		_, err = cmd.ImageManager.DeleteLocation(imageID, datacenter)
		if err != nil {
			return err
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The location was removed successfully!"))

	}
	return nil
}

func buildLocation(location string) []datatypes.Location {
	locations := datatypes.Location{}
	datacenter := []datatypes.Location{}
	match, _ := regexp.MatchString("[a-z]", strings.ToLower(location))
	if match {
		locations.Name = &location
	} else {
		identifier, _ := strconv.Atoi(location)
		locations.Id = &identifier
	}
	datacenter = append(datacenter, locations)
	return datacenter

}

func ImageDatacenterMetaData() cli.Command {
	return cli.Command{
		Category:    "image",
		Name:        "datacenter",
		Description: T("Add/Remove datacenter of an image."),
		Usage: T(`${COMMAND_NAME} sl image datacenter IDENTIFIER [OPTIONS] 

EXAMPLE:
	${COMMAND_NAME} sl image datacenter 12345678 --add dal05 --remove sjc03
	This command Add/Remove datacenter of an image.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "add",
				Usage: T("To add Datacenter"),
			},
			cli.StringFlag{
				Name:  "remove",
				Usage: T("Datacenter to remove"),
			},
			metadata.OutputFlag(),
		},
	}
}
