package image

import (
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DatacenterCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
	Add          string
	Remove       string
}

func NewDatacenterCommand(sl *metadata.SoftlayerCommand) (cmd *DatacenterCommand) {
	thisCmd := &DatacenterCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "datacenter " + T("IDENTIFIER"),
		Short: T("Add/Remove datacenter of an image."),
		Long: T(`${COMMAND_NAME} sl image datacenter IDENTIFIER [OPTIONS] 

EXAMPLE:
	${COMMAND_NAME} sl image datacenter 12345678 --add dal05 --remove sjc03
	This command Add/Remove datacenter of an image.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Add, "add", "", T("To add Datacenter"))
	cobraCmd.Flags().StringVar(&thisCmd.Remove, "remove", "", T("Datacenter to remove"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DatacenterCommand) Run(args []string) error {
	if cmd.Add == "" && cmd.Remove == "" {
		return slErr.NewInvalidUsageError(T("This command requires --add or --remove option."))
	}

	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Image ID")
	}

	storageLocations, err := cmd.ImageManager.GetDatacenters(imageID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get image datacenters.\n"), err.Error(), 2)
	}

	if cmd.Add != "" {
		datacenter := buildLocation(cmd.Add, storageLocations)
		if datacenter[0].Id != nil {
			_, err = cmd.ImageManager.AddLocation(imageID, datacenter)
			if err != nil {
				return err
			}
			cmd.UI.Ok()
			i18nsubs := map[string]interface{}{"imageId": imageID, "datacenter": cmd.Add, "action": "added"}
			cmd.UI.Print(T("{{.imageId}} was {{.action}} from datacenter {{.datacenter}}", i18nsubs))

		} else {
			return slErr.NewInvalidUsageError(T("{{.datacenter}} is invalid", map[string]interface{}{"datacenter": cmd.Add}))
		}
	}
	if cmd.Remove != "" {
		datacenter := buildLocation(cmd.Remove, storageLocations)
		if datacenter[0].Id != nil {
			_, err = cmd.ImageManager.DeleteLocation(imageID, datacenter)
			if err != nil {
				return err
			}
			cmd.UI.Ok()
			i18nsubs := map[string]interface{}{"imageId": imageID, "datacenter": cmd.Remove, "action": "removed"}
			cmd.UI.Print(T("{{.imageId}} was {{.action}} from datacenter {{.datacenter}}", i18nsubs))

		} else {
			return slErr.NewInvalidUsageError(T("{{.datacenter}} is invalid", map[string]interface{}{"datacenter": cmd.Remove}))
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
