package objectstorage

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EndpointData struct {
	LocationRegion string
	Url            string
	EndPointType   string
	PublicPrivate  string
	Legacy         string
}

type EndpointsCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewEndpointsCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *EndpointsCommand) {
	return &EndpointsCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func EndpointsMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "endpoints",
		Description: T("List object storage endpoints."),
		Usage:       T(`${COMMAND_NAME} sl object-storage endpoints IDENTIFIER`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *EndpointsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument"))
	}

	HubNetworkStorageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Invoice ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	endpoints, err := cmd.ObjectStorageManager.GetEndpoints(HubNetworkStorageID)
	if err != nil {
		return cli.NewExitError(T("Failed to get list object storage endpoints.")+err.Error(), 2)
	}
	PrintEndpoints(endpoints, cmd.UI, outputFormat)
	return nil
}

func PrintEndpoints(endpoints []datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint, ui terminal.UI, outputFormat string) {
	table := ui.Table([]string{
		T("Location/Region"),
		T("Url"),
		T("EndPoint Type"),
		T("Public/Private"),
		T("Legacy"),
	})

	allArrays := []EndpointData{}
	for _, endpoint := range endpoints {
		data := EndpointData{
			LocationRegion: LocationRegion(endpoint),
			Url:            *endpoint.Url,
			EndPointType:   EndPointTypeReturn(*endpoint.Region),
			PublicPrivate:  PublicPrivate(*endpoint.Type),
			Legacy:         LegacyReturn(*endpoint.Legacy),
		}
		allArrays = append(allArrays, data)
	}

	allArrays = SortEndpoint(allArrays)
	for _, array := range allArrays {
		table.Add(
			array.LocationRegion,
			array.Url,
			array.EndPointType,
			array.PublicPrivate,
			array.Legacy,
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}

func LegacyReturn(data bool) string {
	if data {
		return T("True")
	}
	return T("False")
}

func EndPointTypeReturn(endpoint string) string {
	if endpoint == "singleSite" {
		return T("Single Site")
	}
	if endpoint == "regional" {
		return T("Region")
	}
	return T("Cross Region")
}

func PublicPrivate(data string) string {
	if data == "public" {
		return T("Public")
	}
	return T("Private")
}

func LocationRegion(endpoint datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint) string {
	if endpoint.Location != nil {
		return *endpoint.Location
	}
	return *endpoint.Region
}

func SortEndpoint(endpoints []EndpointData) []EndpointData {
	endpoint_type := ""
	firstItem := 0
	if len(endpoints) > 0 {
		endpoint_type = endpoints[firstItem].EndPointType
	}
	public := []EndpointData{}
	private := []EndpointData{}
	array_final := []EndpointData{}
	for _, endpoint := range endpoints {
		if endpoint.EndPointType != endpoint_type {
			endpoint_type = endpoint.EndPointType
			array_final = append(array_final, public...)
			array_final = append(array_final, private...)
			public = []EndpointData{}
			private = []EndpointData{}
		}
		if endpoint.PublicPrivate == T("Public") {
			public = append(public, endpoint)
		} else {
			private = append(private, endpoint)
		}
	}
	array_final = append(array_final, public...)
	array_final = append(array_final, private...)
	return array_final
}
