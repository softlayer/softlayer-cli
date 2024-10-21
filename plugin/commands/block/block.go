package block

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	StorageCommand := &metadata.SoftlayerStorageCommand{
		SoftlayerCommand: sl,
		StorageI18n:      map[string]interface{}{"storageType": "block"},
		StorageType:      "block",
	}
	cobraCmd := &cobra.Command{
		Use:   "block",
		Short: T("Classic infrastructure Block Storage"),
		RunE:  nil,
	}
	// Access
	cobraCmd.AddCommand(NewAccessAuthorizeCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewAccessPasswordCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewAccessListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewAccessRevokeCommand(StorageCommand).Command)
	// Replica
	cobraCmd.AddCommand(NewReplicaFailbackCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaFailoverCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaLocationsCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaPartnersCommand(StorageCommand).Command)
	// Snapshot
	cobraCmd.AddCommand(NewSnapshotCancelCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotSetNotificationCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotGetNotificationStatusCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotCreateCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotDeleteCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotDisableCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotEnableCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotRestoreCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSnapshotScheduleListCommand(StorageCommand).Command)
	// Volume
	cobraCmd.AddCommand(NewVolumeCancelCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeCountCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeDetailCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeDuplicateCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeLunCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeOrderCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeModifyCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeOptionsCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeLimitCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeRefreshCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeConvertCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewVolumeSetNoteCommand(StorageCommand).Command)
	// Miscellaneous Commands
	cobraCmd.AddCommand(NewObjectListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewObjectStorageDetailCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewObjectStoragePermissionCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewDisasterRecoveryFailoverCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewDuplicateConvertStatusCommand(StorageCommand).Command)

	// Subnets
	cobraCmd.AddCommand(NewSubnetsListCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSubnetsAssignCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewSubnetsRemoveCommand(StorageCommand).Command)

	return cobraCmd
}

func BlockNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "block",
		Description: T("Classic infrastructure Block Storage"),
	}
}
