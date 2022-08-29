package nas

import (
	"fmt"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	NasNetworkStorageManager managers.NasNetworkStorageManager
	Command                  *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand:         sl,
		NasNetworkStorageManager: managers.NewNasNetworkStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List NAS accounts."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	nasNetworkStorages, err := cmd.NasNetworkStorageManager.ListNasNetworkStorages("")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get NAS Network Storages."), err.Error(), 2)
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
