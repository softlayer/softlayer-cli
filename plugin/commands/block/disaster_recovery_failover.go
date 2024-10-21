package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DisasterRecoveryFailoverCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewDisasterRecoveryFailoverCommand(sl *metadata.SoftlayerStorageCommand) *DisasterRecoveryFailoverCommand {
	thisCmd := &DisasterRecoveryFailoverCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "disaster-recovery-failover " + T("IDENTIFIER") + " " + T("REPLICA_ID"),
		Short: T("Failover an inaccessible volume to its available replicant volume."),
		Long: T(`If a volume (with replication) becomes inaccessible due to a disaster event, this method can be used to immediately
failover to an available replica in another location. This method does not allow for fail back via the API.
To fail back to the original volume after using this method, open a support ticket.
To test failover, use '${COMMAND_NAME} sl {{.storageType}} replica-failover' instead.

EXAMPLE:
	${COMMAND_NAME} sl {{.storageType}} disaster-recovery-failover 12345678 87654321
	This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`, sl.StorageI18n),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DisasterRecoveryFailoverCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}
	replicaID, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Replica ID")
	}
	err = cmd.StorageManager.DisasterRecoveryFailover(volumeID, replicaID)
	subs := map[string]interface{}{"VolumeID": volumeID, "ReplicaID": replicaID}
	if err != nil {
		return slErr.NewAPIError(T("Failover operation could not be initiated for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failover of volume {{.VolumeID}} to replica {{.ReplicaID}} is now in progress.", subs))
	return nil
}
