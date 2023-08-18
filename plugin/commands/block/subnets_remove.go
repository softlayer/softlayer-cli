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

type SubnetsRemoveCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	SubnetIds      []int
}

func NewSubnetsRemoveCommand(sl *metadata.SoftlayerStorageCommand) *SubnetsRemoveCommand {
	thisCmd := &SubnetsRemoveCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "subnets-remove " + T("IDENTIFIER"),
		Short: T("Remove block storage subnets to the given host id."),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} subnets-remove ACCESS_ID [OPTIONS]

access_id is the host_id obtained by: slcli block access-list <volume_id>
SoftLayer_Account::iscsiisolationdisabled must be False to use this command

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} subnets-remove 111111 --subnet-id 222222
   ${COMMAND_NAME} sl {{.storageType}} subnets-remove 111111 --subnet-id 222222 --subnet-id 333333
   ACCESS_ID is the host_id obtained by: ibmcloud sl {{.storageType}} access-list <volume_id>`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntSliceVar(&thisCmd.SubnetIds, "subnet-id", []int{}, T("IDs of the subnets to remove"))
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("subnet-id")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SubnetsRemoveCommand) Run(args []string) error {

	accessID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Access ID")
	}

	subnetsToRemove := cmd.SubnetIds

	subnetsResponse, err := cmd.StorageManager.RemoveSubnetsFromAcl(accessID, subnetsToRemove)
	if err != nil {
		return slErr.NewAPIError(T("Failed to remove subnets."), err.Error(), 2)
	}

	for _, subnet := range subnetsToRemove {
		values := map[string]interface{}{"subnetID": subnet, "accessID": accessID}
		if utils.IntInSlice(subnet, subnetsResponse) != -1 {
			cmd.UI.Print(T("Successfully removed subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		} else {
			cmd.UI.Print(T("Failed to remove subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		}
	}

	return nil
}
