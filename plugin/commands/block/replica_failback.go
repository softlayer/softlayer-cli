package block

import (
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ReplicaFailbackCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewReplicaFailbackCommand(sl *metadata.SoftlayerStorageCommand) *ReplicaFailbackCommand {
	thisCmd := &ReplicaFailbackCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "replica-failback " + T("IDENTIFIER"),
		Short: T("Failback a {{.storageType}} volume from replica", sl.StorageI18n),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} replica-failback VOLUME_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} replica-failback 12345678
   This command performs failback operation for volume with ID 12345678.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReplicaFailbackCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}
	err = cmd.StorageManager.FailBackFromReplicant(volumeID)
	subs := map[string]interface{}{"VolumeID": volumeID}
	if err != nil {
		return slErr.NewAPIError(T("Failback operation could not be initiated for volume {{.VolumeID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Failback of volume {{.VolumeID}} is now in progress.", subs))
	return nil
}
