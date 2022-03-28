package image

import (
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

	storageLocations, err := cmd.ImageManager.GetDatacenters(imageID)
	if err != nil {
		return cli.NewExitError(T("Failed to get image datacenters.\n")+err.Error(), 2)
	}

	if c.IsSet("add") {
		datacenter := buildLocation(c.String("add"), storageLocations)
		if datacenter[0].Id != nil {
			_, err = cmd.ImageManager.AddLocation(imageID, datacenter)
			if err != nil {
				return err
			}
			cmd.UI.Ok()
			i18nsubs := map[string]interface{}{"imageId": imageID, "datacenter": c.String("add"), "action": "added"}
			cmd.UI.Print(T("{{.imageId}} was {{.action}} from datacenter {{.datacenter}}", i18nsubs))

		} else {
			return slErr.NewInvalidUsageError(T("{{.datacenter}} is invalid", map[string]interface{}{"datacenter": c.String("add")}))
		}
	}
	if c.IsSet("remove") {
		datacenter := buildLocation(c.String("remove"), storageLocations)
		if datacenter[0].Id != nil {
			_, err = cmd.ImageManager.DeleteLocation(imageID, datacenter)
			if err != nil {
				return err
			}
			cmd.UI.Ok()
			i18nsubs := map[string]interface{}{"imageId": imageID, "datacenter": c.String("remove"), "action": "removed"}
			cmd.UI.Print(T("{{.imageId}} was {{.action}} from datacenter {{.datacenter}}", i18nsubs))

		} else {
			return slErr.NewInvalidUsageError(T("{{.datacenter}} is invalid", map[string]interface{}{"datacenter": c.String("remove")}))
		}

	}
	return nil
}

func buildLocation(location string, storageLocations []datatypes.Location) []datatypes.Location {
	locations := datatypes.Location{}
	datacenter := []datatypes.Location{}
	idLocation, err := strconv.Atoi(location)
	if err != nil {
		for _, storageLocation := range storageLocations {
			if *storageLocation.Name == strings.ToLower(location) {
				locations.Id = storageLocation.Id
			}
		}
	} else {
		locations.Id = &idLocation
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
