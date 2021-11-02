package file

import (
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)


func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		//file - 25
		NS_FILE_NAME + "-" + CMD_FILE_ACCESS_AUTHORIZE_NAME: func(c *cli.Context) error {
			return NewAccessAuthorizeCommand(ui, storageManager, networkManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_ACCESS_LIST_NAME: func(c *cli.Context) error {
			return NewAccessListCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_ACCESS_REVOKE_NAME: func(c *cli.Context) error {
			return NewAccessRevokeCommand(ui, storageManager, networkManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_REPLICA_FAILBACK_NAME: func(c *cli.Context) error {
			return NewReplicaFailbackCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_REPLICA_FAILOVER_NAME: func(c *cli.Context) error {
			return NewReplicaFailoverCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_REPLICA_LOCATIONS_NAME: func(c *cli.Context) error {
			return NewReplicaLocationsCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_REPLICA_ORDER_NAME: func(c *cli.Context) error {
			return NewReplicaOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_REPLICA_PARTNERS_NAME: func(c *cli.Context) error {
			return NewReplicaPartnersCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_CANCEL_NAME: func(c *cli.Context) error {
			return NewSnapshotCancelCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_CREATE_NAME: func(c *cli.Context) error {
			return NewSnapshotCreateCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_DELETE_NAME: func(c *cli.Context) error {
			return NewSnapshotDeleteCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_DISABLE_NAME: func(c *cli.Context) error {
			return NewSnapshotDisableCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_ENABLE_NAME: func(c *cli.Context) error {
			return NewSnapshotEnableCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_LIST_NAME: func(c *cli.Context) error {
			return NewSnapshotListCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_ORDER_NAME: func(c *cli.Context) error {
			return NewSnapshotOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_RESTORE_NAME: func(c *cli.Context) error {
			return NewSnapshotRestoreCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_SNAPSHOT_SCHEDULE_LIST_NAME: func(c *cli.Context) error {
			return NewSnapshotScheduleListCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_CANCEL_NAME: func(c *cli.Context) error {
			return NewVolumeCancelCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_COUNT_NAME: func(c *cli.Context) error {
			return NewVolumeCountCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_DETAIL_NAME: func(c *cli.Context) error {
			return NewVolumeDetailCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_DUPLICATE_NAME: func(c *cli.Context) error {
			return NewVolumeDuplicateCommand(ui, storageManager, context).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_LIST_NAME: func(c *cli.Context) error {
			return NewVolumeListCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_ORDER_NAME: func(c *cli.Context) error {
			return NewVolumeOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_MODIFY_NAME: func(c *cli.Context) error {
			return NewVolumeModifyCommand(ui, storageManager, context).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_OPTIONS_NAME: func(c *cli.Context) error {
			return NewVolumeOptionsCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_LIMITS_NAME: func(c *cli.Context) error {
			return NewVolumeLimitCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_REFRESH_NAME: func(c *cli.Context) error {
			return NewVolumeRefreshCommand(ui, storageManager).Run(c)
		},
		NS_FILE_NAME + "-" + CMD_FILE_VOLUME_CONVERT_NAME: func(c *cli.Context) error {
			return NewVolumeConvertCommand(ui, storageManager).Run(c)
		},
		// sl file disaster-recovery-failover is in commands/block/block.go
	}

	return CommandActionBindings
}