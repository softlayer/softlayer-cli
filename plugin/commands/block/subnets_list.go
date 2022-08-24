package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SubnetsListCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewSubnetsListCommand(sl *metadata.SoftlayerStorageCommand) *SubnetsListCommand {
	thisCmd := &SubnetsListCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "subnets-list " + T("IDENTIFIER"),
		Short: T("List block storage assigned subnets for the given host id."),
		Long: T(`${COMMAND_NAME} sl block subnets-list ACCESS_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block subnets-list 12345678 
   ACCESS_ID is the host_id obtained by: ibmcloud sl block access-list <volume_id>`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SubnetsListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()
	accessID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Access ID")
	}

	subnets, err := cmd.StorageManager.GetSubnetsInAcl(accessID, "")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get subnets."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Network Identifier"), T("CIDR")})
	for _, subnet := range subnets {
		table.Add(
			utils.FormatIntPointer(subnet.Id),
			utils.FormatStringPointer(subnet.NetworkIdentifier),
			utils.FormatIntPointer(subnet.Cidr),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
