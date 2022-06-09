package nas

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI                       terminal.UI
	NasNetworkStorageManager managers.NasNetworkStorageManager
}

func NewListCommand(ui terminal.UI, nasNetworkStorageManager managers.NasNetworkStorageManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:                       ui,
		NasNetworkStorageManager: nasNetworkStorageManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	nasNetworkStorages, err := cmd.NasNetworkStorageManager.ListNasNetworkStorages("")
	if err != nil {
		return cli.NewExitError(T("Failed to get NAS Network Storages.")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Datacenter"), T("Size"), T("Server")})
	for _, nasNetworkStorage := range nasNetworkStorages {
		location := "-"
		if nasNetworkStorage.ServiceResource.Datacenter != nil {
			location = *nasNetworkStorage.ServiceResource.Datacenter.Name
		}
		table.Add(
			utils.FormatIntPointer(nasNetworkStorage.Id),
			location,
			fmt.Sprintf("%dGB", *nasNetworkStorage.CapacityGb),
			utils.FormatStringPointer(nasNetworkStorage.ServiceResourceBackendIpAddress),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func NasListMetaData() cli.Command {
	return cli.Command{
		Category:    "nas",
		Name:        "list",
		Description: T("List NAS accounts."),
		Usage: T(`${COMMAND_NAME} sl nas list

EXAMPLE: 
   ${COMMAND_NAME} sl nas list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
