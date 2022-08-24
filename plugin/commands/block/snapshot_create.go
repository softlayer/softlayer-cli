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

type SnapshotCreateCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Note           string
}

func NewSnapshotCreateCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotCreateCommand {
	thisCmd := &SnapshotCreateCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-create " + T("IDENTIFIER"),
		Short: T("Create a snapshot on a given volume"),
		Long: T(`${COMMAND_NAME} sl block snapshot-create VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block snapshot-create 12345678 --note snapshotforibmcloud
   This command creates a snapshot for volume with ID 12345678 and with addition note as snapshotforibmcloud.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Note, "note", "n", "", T("Notes to set on the new snapshot"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotCreateCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	snapshot, err := cmd.StorageManager.CreateSnapshot(volumeID, cmd.Note)
	if err != nil {
		return slErr.NewAPIError(T("Error occurred while creating snapshot for volume {{.VolumeID}}.Ensure volume is not failed over or in another state which prevents taking snapshots.\n",
			map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, snapshot)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("New snapshot {{.SnapshotId}} was created.", map[string]interface{}{"SnapshotId": *snapshot.Id}))
	return nil
}
