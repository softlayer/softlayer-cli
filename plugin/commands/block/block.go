package block

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)
	networkManager := managers.NewNetworkManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		NS_BLOCK_NAME + "-" + CMD_BLK_ACCESS_AUTHORIZE_NAME: func(c *cli.Context) error {
			return NewAccessAuthorizeCommand(ui, storageManager, networkManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_ACCESS_LIST_NAME: func(c *cli.Context) error {
			return NewAccessListCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_ACCESS_PASSWORD_NAME: func(c *cli.Context) error {
			return NewAccessPasswordCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_ACCESS_REVOKE_NAME: func(c *cli.Context) error {
			return NewAccessRevokeCommand(ui, storageManager, networkManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_REPLICA_FAILBACK_NAME: func(c *cli.Context) error {
			return NewReplicaFailbackCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_REPLICA_FAILOVER_NAME: func(c *cli.Context) error {
			return NewReplicaFailoverCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_REPLICA_LOCATIONS_NAME: func(c *cli.Context) error {
			return NewReplicaLocationsCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_REPLICA_ORDER_NAME: func(c *cli.Context) error {
			return NewReplicaOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_REPLICA_PARTNERS_NAME: func(c *cli.Context) error {
			return NewReplicaPartnersCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_CANCEL_NAME: func(c *cli.Context) error {
			return NewSnapshotCancelCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_SET_NOTIFICATION_NAME: func(c *cli.Context) error {
			return NewSnapshotSetNotificationCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_GET_NOTIFIACTION_STATUS_NAME: func(c *cli.Context) error {
			return NewSnapshotGetNotificationStatusCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_CREATE_NAME: func(c *cli.Context) error {
			return NewSnapshotCreateCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_DELETE_NAME: func(c *cli.Context) error {
			return NewSnapshotDeleteCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_DISABLE_NAME: func(c *cli.Context) error {
			return NewSnapshotDisableCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_ENABLE_NAME: func(c *cli.Context) error {
			return NewSnapshotEnableCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_LIST_NAME: func(c *cli.Context) error {
			return NewSnapshotListCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_ORDER_NAME: func(c *cli.Context) error {
			return NewSnapshotOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_RESTORE_NAME: func(c *cli.Context) error {
			return NewSnapshotRestoreCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_SNAPSHOT_SCHEDULE_LIST_NAME: func(c *cli.Context) error {
			return NewSnapshotScheduleListCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_CANCEL_NAME: func(c *cli.Context) error {
			return NewVolumeCancelCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_COUNT_NAME: func(c *cli.Context) error {
			return NewVolumeCountCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_DETAIL_NAME: func(c *cli.Context) error {
			return NewVolumeDetailCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_DUPLICATE_NAME: func(c *cli.Context) error {
			return NewVolumeDuplicateCommand(ui, storageManager, context).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_LIST_NAME: func(c *cli.Context) error {
			return NewVolumeListCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_LUN_NAME: func(c *cli.Context) error {
			return NewVolumeLunCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_ORDER_NAME: func(c *cli.Context) error {
			return NewVolumeOrderCommand(ui, storageManager, context).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_MODIFY_NAME: func(c *cli.Context) error {
			return NewVolumeModifyCommand(ui, storageManager, context).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_OPTIONS_NAME: func(c *cli.Context) error {
			return NewVolumeOptionsCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_LIMITS_NAME: func(c *cli.Context) error {
			return NewVolumeLimitCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_REFRESH_NAME: func(c *cli.Context) error {
			return NewVolumeRefreshCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_VOLUME_CONVERT_NAME: func(c *cli.Context) error {
			return NewVolumeConvertCommand(ui, storageManager).Run(c)
		},
		NS_BLOCK_NAME + "-" + CMD_BLK_DISASTER_FAILOVER_NAME: func(c *cli.Context) error {
			return NewDisasterRecoveryFailoverCommand(ui, storageManager).Run(c)
		},
		// Commands that are the same for file and block go here.
		NS_FILE_NAME + "-" + CMD_BLK_DISASTER_FAILOVER_NAME: func(c *cli.Context) error {
			return NewDisasterRecoveryFailoverCommand(ui, storageManager).Run(c)
		},
	}

	return CommandActionBindings
}
