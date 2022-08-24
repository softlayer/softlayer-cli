package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotGetNotificationStatusCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewSnapshotGetNotificationStatusCommand(sl *metadata.SoftlayerCommand) *SnapshotGetNotificationStatusCommand {
	thisCmd := &SnapshotGetNotificationStatusCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-get-notification-status " + T("IDENTIFIER"),
		Short: T("Get snapshots space usage threshold warning flag setting for a given volume."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotGetNotificationStatusCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()
	subs := map[string]interface{}{"ID": volumeID}
	enabled, err := cmd.StorageManager.GetSnapshotNotificationStatus(volumeID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get the snapshot notification status for volume '{{.ID}}'.\n", subs), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, enabled)
	}

	if enabled == 0 {
		cmd.UI.Print(T("Disabled: Snapshots space usage threshold is disabled for volume '{{.ID}}'.", subs))
	} else {
		cmd.UI.Print(T("Enabled: Snapshots space usage threshold is enabled for volume '{{.ID}}'.", subs))
	}

	return nil
}
