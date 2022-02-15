package file

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

func FileNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "file",
		Description: T("Classic infrastructure File Storage"),
	}
}

func FileMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "file",
		Description: T("Classic infrastructure File Storage"),
		Usage:       "${COMMAND_NAME} sl file",
		Subcommands: []cli.Command{
			FileAccessAuthorizeMetaData(),
			FileAccessListMetaData(),
			FileAccessRevokeMetaData(),
			FileReplicaFailbackMetaData(),
			FileReplicaFailoverMetaData(),
			FileReplicaLocationsMetaData(),
			FileReplicaOrderMetaData(),
			FileReplicaPartnersMetaData(),
			FileSnapshotCancelMetaData(),
			FileSnapshotCreateMetaData(),
			FileSnapshotDisableMetaData(),
			FileSnapshotEnableMetaData(),
			FileSnapshotDeleteMetaData(),
			FileSnapshotListMetaData(),
			FileSnapshotOrderMetaData(),
			FileSnapshotScheduleListMetaData(),
			FileSnapshotRestoreMetaData(),
			FileVolumeCancelMetaData(),
			FileVolumeCountMetaData(),
			FileVolumeListMetaData(),
			FileVolumeDetailMetaData(),
			FileVolumeDuplicateMetaData(),
			FileVolumeModifyMetaData(),
			FileVolumeOrderMetaData(),
			FileVolumeOptionsMetaData(),
			FileVolumeLimitsMetaData(),
			FileVolumeRefreshMetaData(),
			FileVolumeConvertMetaData(),
			FileDisasterRecoveryFailoverMetaData(),
			FileVolumeSnapshotSetNotificationMetaData(),
			FileVolumeSnapshotGetNotificationStatusMetaData(),
		},
	}
}

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		//file - 25
		"file-access-authorize": func(c *cli.Context) error {
			return NewAccessAuthorizeCommand(ui, storageManager, networkManager).Run(c)
		},
		"file-access-list": func(c *cli.Context) error {
			return NewAccessListCommand(ui, storageManager).Run(c)
		},
		"file-access-revoke": func(c *cli.Context) error {
			return NewAccessRevokeCommand(ui, storageManager, networkManager).Run(c)
		},
		"file-replica-failback": func(c *cli.Context) error {
			return NewReplicaFailbackCommand(ui, storageManager).Run(c)
		},
		"file-replica-failover": func(c *cli.Context) error {
			return NewReplicaFailoverCommand(ui, storageManager).Run(c)
		},
		"file-replica-locations": func(c *cli.Context) error {
			return NewReplicaLocationsCommand(ui, storageManager).Run(c)
		},
		"file-replica-order": func(c *cli.Context) error {
			return NewReplicaOrderCommand(ui, storageManager, context).Run(c)
		},
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
	}

	return CommandActionBindings
}
