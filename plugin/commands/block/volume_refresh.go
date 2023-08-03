package block

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VolumeRefreshCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Force          bool
}

func NewVolumeRefreshCommand(sl *metadata.SoftlayerStorageCommand) *VolumeRefreshCommand {
	thisCmd := &VolumeRefreshCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-refresh " + T("IDENTIFIER") + " " + T("SNAPSHOT_ID"),
		Short: T("Refresh a duplicate volume with a snapshot from its parent."),
		Long: T("Refresh a duplicate volume with a snapshot from its parent.") + " " + T(`${COMMAND_NAME} sl {{.storageType}} volume-refresh VOLUME_ID SNAPSHOT_ID

EXAMPLE:
	${COMMAND_NAME} sl {{.storageType}} volume-refresh VOLUME_ID SNAPSHOT_ID
	Refresh a duplicate VOLUME_ID with a snapshot from its parent SNAPSHOT_ID.`, sl.StorageI18n),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force-refresh", "f", false, T("Force the volume refresh, will cancel any ongoing transactions."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeRefreshCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	snapshotId, err := strconv.Atoi(args[1])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Snapshot ID")
	}

	err = cmd.StorageManager.VolumeRefresh(volumeID, snapshotId, cmd.Force)
	if err != nil {
		return err
	}
	cmd.UI.Ok()
	return nil
}
