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
	// Access
	cobraCmd.AddCommand(NewAccessAuthorizeCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessPasswordCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessListCommand(sl).Command)
	cobraCmd.AddCommand(NewAccessRevokeCommand(sl).Command)
	// Replica
	cobraCmd.AddCommand(NewReplicaFailbackCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaFailoverCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaLocationsCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaOrderCommand(sl).Command)
	cobraCmd.AddCommand(NewReplicaPartnersCommand(sl).Command)
	// Snapshot
	cobraCmd.AddCommand(NewSnapshotCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotSetNotificationCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotGetNotificationStatusCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotCreateCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotDeleteCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotDisableCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotEnableCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotListCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotOrderCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotRestoreCommand(sl).Command)
	cobraCmd.AddCommand(NewSnapshotScheduleListCommand(sl).Command)
	// Volume
	cobraCmd.AddCommand(NewVolumeCancelCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeCountCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeDetailCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeDuplicateCommand(sl).Command)

	return cobraCmd
}

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	storageManager := managers.NewStorageManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

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
