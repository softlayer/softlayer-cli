package objectstorage

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Reason         string
	Immediate      bool
	Force          bool
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) *DeleteCommand {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "delete " + T("IDENTIFIER"),
		Short: T("Cancel an existing block storage volume"),
		Long: T(`${COMMAND_NAME} sl object-storage delete VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl object-storage delete 12345678 --immediate -f 
   This command cancels volume with ID 12345678 immediately and without asking for confirmation.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Reason, "reason", "", T("An optional reason for cancellation"))
	cobraCmd.Flags().BoolVar(&thisCmd.Immediate, "immediate", false, T("Cancel the block storage volume immediately instead of on the billing anniversary"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {

	objectStorageID, err := strconv.Atoi(args[0])
	subs := map[string]interface{}{"ID": objectStorageID, "objectStorageID": objectStorageID}
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("objectStorageID")
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.StorageManager.CancelVolume("block", objectStorageID, cmd.Reason, cmd.Immediate)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return slErr.NewAPIError(T("Unable to find object-storage with objectStorageID {{.objectStorageID}}.", subs), err.Error(), 0)
		}
		return slErr.NewAPIError(T("Failed to cancel object-storage: {{.objectStorageID}}.", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	if cmd.Immediate {
		cmd.UI.Print(T("Object-storage {{.objectStorageID}} has been marked for immediate cancellation.", subs))
	} else {
		cmd.UI.Print(T("Object-storage {{.objectStorageID}} has been marked for cancellation.", subs))
	}
	return nil
}
