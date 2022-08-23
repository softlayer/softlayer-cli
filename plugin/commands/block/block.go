package block

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
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
	cobraCmd.AddCommand(NewVolumeListCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeLunCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeOrderCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeModifyCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeOptionsCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeLimitCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeRefreshCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeConvertCommand(sl).Command)
	cobraCmd.AddCommand(NewVolumeSetNoteCommand(sl).Command)
	// Miscellaneous Commands
	cobraCmd.AddCommand(NewObjectListCommand(sl).Command)
	cobraCmd.AddCommand(NewDisasterRecoveryFailoverCommand(sl).Command)
	cobraCmd.AddCommand(NewDuplicateConvertStatusCommand(sl).Command)

	// Subnets
	cobraCmd.AddCommand(NewSubnetsListCommand(sl).Command)
	cobraCmd.AddCommand(NewSubnetsAssignCommand(sl).Command)
	cobraCmd.AddCommand(NewSubnetsRemoveCommand(sl).Command)

	return cobraCmd
}

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {

	CommandActionBindings := map[string]func(c *cli.Context) error{

		// "block-duplicate-convert-status": func(c *cli.Context) error {
		// 	return NewDuplicateConvertStatusCommand(ui, storageManager).Run(c)
		// },
		// Commands that are the same for file and block go here.
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! MOVE THESE TO FILE TODO!!!!!!!
		// "file-disaster-recovery-failover": func(c *cli.Context) error {
		// 	return NewDisasterRecoveryFailoverCommand(ui, storageManager).Run(c)
		// },
		// "file-volume-set-note": func(c *cli.Context) error {
		// 	return NewVolumeSetNoteCommand(ui, storageManager).Run(c)
		// },
		// "file-duplicate-convert-status": func(c *cli.Context) error {
		// 	return NewDuplicateConvertStatusCommand(ui, storageManager).Run(c)
		// },
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
