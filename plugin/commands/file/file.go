package file

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func FileNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "file",
		Description: T("Classic infrastructure File Storage"),
	}
}

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	StorageCommand := &metadata.SoftlayerStorageCommand{
		SoftlayerCommand: sl,
		StorageI18n:      map[string]interface{}{"storageType": "file"},
		StorageType:      "file",
	}
	cobraCmd := &cobra.Command{
		Use:   "file",
		Short: T("Classic infrastructure File Storage"),
		RunE:  nil,
	}

	// Commands that are the same as their block version.

	cobraCmd.AddCommand(block.NewDisasterRecoveryFailoverCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewVolumeSetNoteCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewDuplicateConvertStatusCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewAccessListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewReplicaFailbackCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewReplicaFailoverCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewReplicaLocationsCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewReplicaPartnersCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotSetNotificationCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotGetNotificationStatusCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotCreateCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotDeleteCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotDisableCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotEnableCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotRestoreCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotScheduleListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewVolumeLimitCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewVolumeRefreshCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewVolumeConvertCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewSnapshotOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(block.NewVolumeOptionsCommand(StorageCommand).Command)

	// Unique File Commands, even these can likely be merged in a later version.
	cobraCmd.AddCommand(NewAccessAuthorizeCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewAccessRevokeCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotCancelCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeCancelCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeCountCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeDetailCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeDuplicateCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeModifyCommand(StorageCommand).Command)
	return cobraCmd
}
