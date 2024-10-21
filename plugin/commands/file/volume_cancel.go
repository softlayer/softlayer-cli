package file

import (
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VolumeCancelCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Reason         string
	Immediate      bool
	Force          bool
}

func NewVolumeCancelCommand(sl *metadata.SoftlayerStorageCommand) *VolumeCancelCommand {
	thisCmd := &VolumeCancelCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-cancel " + T("IDENTIFIER"),
		Short: T("Cancel an existing file storage volume"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} volume-cancel VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} volume-cancel 12345678 --immediate -f 
   This command cancels volume with ID 12345678 immediately and without asking for confirmation.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Reason, "reason", "", T("An optional reason for cancellation"))
	cobraCmd.Flags().BoolVar(&thisCmd.Immediate, "immediate", false, T("Cancel the file storage volume immediately instead of on the billing anniversary"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeCancelCommand) Run(args []string) error {
	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}
	subs := map[string]interface{}{"ID": volumeID, "VolumeId": volumeID}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This will cancel the file volume: {{.ID}} and cannot be undone. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.StorageManager.CancelVolume("file", volumeID, cmd.Reason, cmd.Immediate)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return slErr.NewAPIError(T("Unable to find volume with ID {{.ID}}.\n", subs), err.Error(), 0)
		}
		return slErr.NewAPIError(T("Failed to cancel file volume: {{.ID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	if cmd.Immediate {
		cmd.UI.Print(T("File volume {{.VolumeId}} has been marked for immediate cancellation.", subs))
	} else {
		cmd.UI.Print(T("File volume {{.VolumeId}} has been marked for cancellation.", subs))
	}
	return nil
}
