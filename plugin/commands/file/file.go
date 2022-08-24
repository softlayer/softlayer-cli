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

	// Unique File Commands
	cobraCmd.AddCommand(NewAccessAuthorizeCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewAccessRevokeCommand(StorageCommand).Command)
	cobraCmd.AddCommand(NewReplicaOrderCommand(StorageCommand).Command)

	return cobraCmd
}

/*
func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		//file - 25

		"file-replica-partners": func(c *cli.Context) error {
			return NewReplicaPartnersCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-cancel": func(c *cli.Context) error {
			return NewSnapshotCancelCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-set-notification": func(c *cli.Context) error {
			return NewSnapshotSetNotificationCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-get-notification-status": func(c *cli.Context) error {
			return NewSnapshotGetNotificationStatusCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-create": func(c *cli.Context) error {
			return NewSnapshotCreateCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-delete": func(c *cli.Context) error {
			return NewSnapshotDeleteCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-disable": func(c *cli.Context) error {
			return NewSnapshotDisableCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-enable": func(c *cli.Context) error {
			return NewSnapshotEnableCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-list": func(c *cli.Context) error {
			return NewSnapshotListCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-order": func(c *cli.Context) error {
			return NewSnapshotOrderCommand(ui, storageManager, context).Run(c)
		},
		"file-snapshot-restore": func(c *cli.Context) error {
			return NewSnapshotRestoreCommand(ui, storageManager).Run(c)
		},
		"file-snapshot-schedule-list": func(c *cli.Context) error {
			return NewSnapshotScheduleListCommand(ui, storageManager).Run(c)
		},
		"file-volume-cancel": func(c *cli.Context) error {
			return NewVolumeCancelCommand(ui, storageManager).Run(c)
		},
		"file-volume-count": func(c *cli.Context) error {
			return NewVolumeCountCommand(ui, storageManager).Run(c)
		},
		"file-volume-detail": func(c *cli.Context) error {
			return NewVolumeDetailCommand(ui, storageManager).Run(c)
		},
		"file-volume-duplicate": func(c *cli.Context) error {
			return NewVolumeDuplicateCommand(ui, storageManager, context).Run(c)
		},
		"file-volume-list": func(c *cli.Context) error {
			return NewVolumeListCommand(ui, storageManager).Run(c)
		},
		"file-volume-order": func(c *cli.Context) error {
			return NewVolumeOrderCommand(ui, storageManager, context).Run(c)
		},
		"file-volume-modify": func(c *cli.Context) error {
			return NewVolumeModifyCommand(ui, storageManager, context).Run(c)
		},
		"file-volume-options": func(c *cli.Context) error {
			return NewVolumeOptionsCommand(ui, storageManager).Run(c)
		},
		"file-volume-limits": func(c *cli.Context) error {
			return NewVolumeLimitCommand(ui, storageManager).Run(c)
		},
		"file-volume-refresh": func(c *cli.Context) error {
			return NewVolumeRefreshCommand(ui, storageManager).Run(c)
		},
		"file-volume-convert": func(c *cli.Context) error {
			return NewVolumeConvertCommand(ui, storageManager).Run(c)
		},
		// sl file disaster-recovery-failover is in commands/block/block.go
		// sl file volume-set-note is in commands/block/block.go
		// sl duplicate-convert-status is in commands/block/block.go
	}

	return CommandActionBindings
}

*/
