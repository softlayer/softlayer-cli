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

type VolumeSetNoteCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Note           string
}

func NewVolumeSetNoteCommand(sl *metadata.SoftlayerStorageCommand) *VolumeSetNoteCommand {
	thisCmd := &VolumeSetNoteCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-set-note " + T("IDENTIFIER"),
		Short: T("Set note for an existing {{.storageType}} storage volume.", sl.StorageI18n),
		Long: T("Set note for an existing {{.storageType}} storage volume.", sl.StorageI18n) + " " + T(`${COMMAND_NAME} sl {{.storageType}} volume-set-note [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} volume-set-note 12345678 --note 'this is my note'`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Note, "note", "n", "", T("Public notes related to a Storage volume  [required]"))
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("note")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeSetNoteCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat := cmd.GetOutputFlag()

	successful, err := cmd.StorageManager.VolumeSetNote(volumeID, cmd.Note)
	if err != nil {
		return slErr.NewAPIError(T("Error occurred while adding note to volume: {{.VolumeID}}",
			map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, successful)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The note was set successfully"))
	return nil
}
