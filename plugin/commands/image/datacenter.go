package image

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"regexp"
	"strconv"
	"strings"
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