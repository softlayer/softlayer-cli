package block

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "block",
		Short: T("Classic infrastructure Block Storage"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAccessAuthorizeCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessPasswordCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessListCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessRevokeCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaFailbackCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaFailoverCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaLocationsCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaOrderCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaPartnersCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotSetNotificationCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotGetNotificationStatusCommand(sl).Command)
	return cobraCmd
}

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		"block-snapshot-create": func(c *cli.Context) error {
			return NewSnapshotCreateCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-delete": func(c *cli.Context) error {
			return NewSnapshotDeleteCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-disable": func(c *cli.Context) error {
			return NewSnapshotDisableCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-enable": func(c *cli.Context) error {
			return NewSnapshotEnableCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-list": func(c *cli.Context) error {
			return NewSnapshotListCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-order": func(c *cli.Context) error {
			return NewSnapshotOrderCommand(ui, storageManager, context).Run(c)
		},
		"block-snapshot-restore": func(c *cli.Context) error {
			return NewSnapshotRestoreCommand(ui, storageManager).Run(c)
		},
		"block-snapshot-schedule-list": func(c *cli.Context) error {
			return NewSnapshotScheduleListCommand(ui, storageManager).Run(c)
		},
		"block-volume-cancel": func(c *cli.Context) error {
			return NewVolumeCancelCommand(ui, storageManager).Run(c)
		},
		"block-volume-count": func(c *cli.Context) error {
			return NewVolumeCountCommand(ui, storageManager).Run(c)
		},
		"block-volume-detail": func(c *cli.Context) error {
			return NewVolumeDetailCommand(ui, storageManager).Run(c)
		},
		"block-volume-duplicate": func(c *cli.Context) error {
			return NewVolumeDuplicateCommand(ui, storageManager, context).Run(c)
		},
		"block-volume-list": func(c *cli.Context) error {
			return NewVolumeListCommand(ui, storageManager).Run(c)
		},
		"block-volume-set-lun-id": func(c *cli.Context) error {
			return NewVolumeLunCommand(ui, storageManager).Run(c)
		},
		"block-volume-order": func(c *cli.Context) error {
			return NewVolumeOrderCommand(ui, storageManager, context).Run(c)
		},
		"block-volume-modify": func(c *cli.Context) error {
			return NewVolumeModifyCommand(ui, storageManager, context).Run(c)
		},
		"block-volume-options": func(c *cli.Context) error {
			return NewVolumeOptionsCommand(ui, storageManager).Run(c)
		},
		"block-volume-limits": func(c *cli.Context) error {
			return NewVolumeLimitCommand(ui, storageManager).Run(c)
		},
		"block-volume-refresh": func(c *cli.Context) error {
			return NewVolumeRefreshCommand(ui, storageManager).Run(c)
		},
		"block-volume-convert": func(c *cli.Context) error {
			return NewVolumeConvertCommand(ui, storageManager).Run(c)
		},
		"block-object-list": func(c *cli.Context) error {
			return NewObjectListCommand(ui, storageManager).Run(c)
		},
		"block-subnets-list": func(c *cli.Context) error {
			return NewSubnetsListCommand(ui, storageManager).Run(c)
		},
		"block-subnets-assign": func(c *cli.Context) error {
			return NewSubnetsAssignCommand(ui, storageManager).Run(c)
		},
		"block-subnets-remove": func(c *cli.Context) error {
			return NewSubnetsRemoveCommand(ui, storageManager).Run(c)
		},
		"block-disaster-recovery-failover": func(c *cli.Context) error {
			return NewDisasterRecoveryFailoverCommand(ui, storageManager).Run(c)
		},
		"block-volume-set-note": func(c *cli.Context) error {
			return NewVolumeSetNoteCommand(ui, storageManager).Run(c)
		},
		"block-duplicate-convert-status": func(c *cli.Context) error {
			return NewDuplicateConvertStatusCommand(ui, storageManager).Run(c)
		},
		// Commands that are the same for file and block go here.
		"file-disaster-recovery-failover": func(c *cli.Context) error {
			return NewDisasterRecoveryFailoverCommand(ui, storageManager).Run(c)
		},
		"file-volume-set-note": func(c *cli.Context) error {
			return NewVolumeSetNoteCommand(ui, storageManager).Run(c)
		},
		"file-duplicate-convert-status": func(c *cli.Context) error {
			return NewDuplicateConvertStatusCommand(ui, storageManager).Run(c)
		},
	}

	return CommandActionBindings
}

func BlockNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "block",
		Description: T("Classic infrastructure Block Storage"),
	}
}
