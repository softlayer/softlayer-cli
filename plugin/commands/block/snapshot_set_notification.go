package block

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SnapshotSetNotificationCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Enable         bool
	Disable        bool
}

func NewSnapshotSetNotificationCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotSetNotificationCommand {
	thisCmd := &SnapshotSetNotificationCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-set-notification " + T("IDENTIFIER"),
		Short: T("Enables/Disables snapshot space usage threshold warning for a given volume."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Enable, "enable", false, T("Enable snapshot notification."))
	cobraCmd.Flags().BoolVar(&thisCmd.Disable, "disable", false, T("Disable snapshot notification."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotSetNotificationCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	if cmd.Enable && cmd.Disable {
		return slErr.NewExclusiveFlagsErrorWithDetails([]string{"--enable", "--disable"}, "")
	}

	if !cmd.Enable && !cmd.Disable {
		return slErr.NewInvalidUsageError(T("Either '--enable' or '--disable' is required."))
	}

	enabled := !cmd.Disable
	subs := map[string]interface{}{"ID": volumeID, "ENABLE": enabled}
	if err = cmd.StorageManager.SetSnapshotNotification(volumeID, enabled); err != nil {
		return slErr.NewAPIError(T("Failed to set the snapshort notification  for volume '{{.ID}}'.\n", subs), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Snapshots space usage threshold warning notification has been set to '{{.ENABLE}}' for volume '{{.ID}}'.", subs))
	return nil
}
