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

const (
	LEGACY_TRUE  = "True"
	LEGACY_FALSE = "False"
	CROS         = "Cross Region"
	REGION       = "Region"
	SINGLE       = "Single Site"
	PUBLIC       = "Public"
	PRIVATE      = "Private"
)

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
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
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
		T("Legacy"),
		T("EndPoint Type"),
		T("Public/Private"),
		T("Location/Region"),
		T("Url"),
	})

	allArrays := [][]string{}

	for _, endpoint := range endpoints {
		data := []string{LegacyReturn(*endpoint.Legacy), EndPointTypeReturn(*endpoint.Region), PublicPrivate(*endpoint.Type), LocationRegion(endpoint), *endpoint.Url}
		allArrays = append(allArrays, data)
	}

	allArrays = SortEndpoint(allArrays)
	for _, array := range allArrays {
		table.Add(
			array[0], array[1], array[2], array[3], array[4],
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}

func LegacyReturn(data bool) string {
	if data {
		return "True"
	}
	return "False"
}

func EndPointTypeReturn(endpoint string) string {
	if endpoint == "singleSite" {
		return "Single Site"
	}
	if endpoint == "regional" {
		return "Region"
	}
	return "Cross Region"
}

func PublicPrivate(data string) string {
	if data == "public" {
		return "Public"
	}
	return "Private"
}

func LocationRegion(endpoint datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint) string {
	if endpoint.Location != nil {
		return *endpoint.Location
	}
	return *endpoint.Region
}

func SortEndpoint(endpoints [][]string) [][]string {
	endpoint_type := ""
	if len(endpoints) > 0 {
		endpoint_type = endpoints[0][1]
	}
	public := [][]string{}
	private := [][]string{}
	array_final := [][]string{}
	for _, endpoint := range endpoints {
		if endpoint[1] != endpoint_type {
			endpoint_type = endpoint[1]
			array_final = append(array_final, public...)
			array_final = append(array_final, private...)
			public = [][]string{}
			private = [][]string{}
		}
		if endpoint[2] == "Public" {
			public = append(public, endpoint)
		} else {
			private = append(private, endpoint)
		}
	}
	array_final = append(array_final, public...)
	array_final = append(array_final, private...)
	return array_final
}
