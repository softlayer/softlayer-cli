package objectstorage

import (
	"fmt"
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
		Usage:       T(`${COMMAND_NAME} sl object-storage endpoints`),
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

	arrayCrossPublic := [][]string{}
	arrayCrossPrivate := [][]string{}
	arrayRegionPublic := [][]string{}
	arrayRegionPrivate := [][]string{}
	arraySinglePublic := [][]string{}
	arraySinglePrivate := [][]string{}

	legacyArrayCrossPublic := [][]string{}
	legacyArrayCrossPrivate := [][]string{}
	legacyArrayRegionPublic := [][]string{}
	legacyArrayRegionPrivate := [][]string{}
	legacyArraySinglePublic := [][]string{}
	legacyArraySinglePrivate := [][]string{}
	allArrays := [][]string{}

	for _, endpoint := range endpoints {
		if !*endpoint.Legacy {
			if *endpoint.Region == "singleSite" {
				if endpoint.Location != nil {
					if *endpoint.Type == "public" {
						data := []string{LEGACY_FALSE, SINGLE, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
						arraySinglePublic = append(arraySinglePublic, data)
					} else {
						data := []string{LEGACY_FALSE, SINGLE, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
						arraySinglePrivate = append(arraySinglePrivate, data)
					}
				} else {
					fmt.Println("paso?")
					if *endpoint.Type == "public" {
						data := []string{LEGACY_FALSE, SINGLE, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
						arraySinglePublic = append(arraySinglePublic, data)
					} else {
						data := []string{LEGACY_FALSE, SINGLE, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
						arraySinglePrivate = append(arraySinglePrivate, data)
					}
				}
			} else {
				if *endpoint.Region == "regional" {
					if endpoint.Location != nil {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_FALSE, REGION, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							arrayRegionPublic = append(arrayRegionPublic, data)
						} else {
							data := []string{LEGACY_FALSE, REGION, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							arrayRegionPrivate = append(arrayRegionPrivate, data)
						}
					} else {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_FALSE, REGION, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							arrayRegionPublic = append(arrayRegionPublic, data)
						} else {
							data := []string{LEGACY_FALSE, REGION, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							arrayRegionPrivate = append(arrayRegionPrivate, data)
						}
					}
				} else {

					if endpoint.Location != nil {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_FALSE, CROS, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							arrayCrossPublic = append(arrayCrossPublic, data)
						} else {
							data := []string{LEGACY_FALSE, CROS, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							arrayCrossPrivate = append(arrayCrossPrivate, data)
						}
					} else {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_FALSE, CROS, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							arrayCrossPublic = append(arrayCrossPublic, data)
						} else {
							data := []string{LEGACY_FALSE, CROS, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							arrayCrossPrivate = append(arrayCrossPrivate, data)
						}
					}
				}
			}
		} else {
			if *endpoint.Region == "singleSite" {
				if endpoint.Location != nil {
					if *endpoint.Type == "public" {
						data := []string{LEGACY_TRUE, SINGLE, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
						legacyArraySinglePublic = append(legacyArraySinglePublic, data)
					} else {
						data := []string{LEGACY_TRUE, SINGLE, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
						legacyArraySinglePrivate = append(legacyArraySinglePrivate, data)
					}
				} else {
					if *endpoint.Type == "public" {
						data := []string{LEGACY_TRUE, SINGLE, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
						legacyArraySinglePublic = append(legacyArraySinglePublic, data)
					} else {
						data := []string{LEGACY_TRUE, SINGLE, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
						legacyArraySinglePrivate = append(legacyArraySinglePrivate, data)
					}
				}
			} else {
				if *endpoint.Region == "regional" {
					if endpoint.Location != nil {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_TRUE, REGION, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayRegionPublic = append(legacyArrayRegionPublic, data)
						} else {
							data := []string{LEGACY_TRUE, REGION, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayRegionPrivate = append(legacyArrayRegionPrivate, data)
						}
					} else {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_TRUE, REGION, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayRegionPublic = append(legacyArrayRegionPublic, data)
						} else {
							data := []string{LEGACY_TRUE, REGION, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayRegionPrivate = append(legacyArrayRegionPrivate, data)
						}
					}
				} else {

					if endpoint.Location != nil {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_TRUE, CROS, PUBLIC, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayCrossPublic = append(legacyArrayCrossPublic, data)
						} else {
							data := []string{LEGACY_TRUE, CROS, PRIVATE, utils.FormatStringPointerName(endpoint.Location), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayCrossPrivate = append(legacyArrayCrossPrivate, data)
						}
					} else {
						if *endpoint.Type == "public" {
							data := []string{LEGACY_TRUE, CROS, PUBLIC, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayCrossPublic = append(legacyArrayCrossPublic, data)
						} else {
							data := []string{LEGACY_TRUE, CROS, PRIVATE, utils.FormatStringPointerName(endpoint.Region), utils.FormatStringPointerName(endpoint.Url)}
							legacyArrayCrossPrivate = append(legacyArrayCrossPrivate, data)
						}
					}
				}
			}
		}
	}

	allArrays = append(allArrays, arrayCrossPublic...)
	allArrays = append(allArrays, arrayCrossPrivate...)
	allArrays = append(allArrays, arrayRegionPublic...)
	allArrays = append(allArrays, arrayRegionPrivate...)
	allArrays = append(allArrays, arraySinglePublic...)
	allArrays = append(allArrays, arraySinglePrivate...)
	allArrays = append(allArrays, legacyArrayCrossPublic...)
	allArrays = append(allArrays, legacyArrayCrossPrivate...)
	allArrays = append(allArrays, legacyArrayRegionPublic...)
	allArrays = append(allArrays, legacyArrayRegionPrivate...)
	allArrays = append(allArrays, legacyArraySinglePublic...)
	allArrays = append(allArrays, legacyArraySinglePrivate...)

	for _, array := range allArrays {
		table.Add(
			array[0],
			array[1],
			array[2],
			array[3],
			array[4],
		)
	}
	utils.PrintTable(ui, table, outputFormat)
}
