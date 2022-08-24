package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ReplicaFailoverCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewReplicaFailoverCommand(sl *metadata.SoftlayerStorageCommand) *ReplicaFailoverCommand {
	thisCmd := &ReplicaFailoverCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "replica-failove " + T("IDENTIFIER") + " " + T("REPLICA_ID"),
		Short: T("Failover a block volume to the given replica volume"),
		Long: T(`${COMMAND_NAME} sl block replica-failover VOLUME_ID REPLICA_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-failover 12345678 87654321
   This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReplicaFailoverCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	replicaID, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Replica ID")
	}
	err = cmd.StorageManager.FailOverToReplicant(volumeID, replicaID)
	subs := map[string]interface{}{"VolumeID": volumeID, "ReplicaID": replicaID}
	if err != nil {
		return slErr.NewAPIError(T("Failover operation could not be initiated for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failover of volume {{.VolumeID}} to replica {{.ReplicaID}} is now in progress.", subs))
	return nil
}
