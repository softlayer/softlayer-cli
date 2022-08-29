package block

import (
	"fmt"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"strconv"
)

type AccessPasswordCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
}

func NewAccessPasswordCommand(sl *metadata.SoftlayerStorageCommand) *AccessPasswordCommand {
	thisCmd := &AccessPasswordCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-password " + T("IDENTIFIER") + " " + T("PASSWORD"),
		Short: T("Changes a password for a volume's access"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} access-password ACCESS_ID PASSWORD
	
	ACCESS_ID is the allowed_host_id from '${COMMAND_NAME} sl {{.storageType}} access-list'`, sl.StorageI18n),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
		DisableFlagsInUseLine: true,
	}

	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *AccessPasswordCommand) Run(args []string) error {

	fmt.Printf("===AccessPasswordCommand===")
	hostID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("allowed access host ID")
	}
	err = cmd.StorageManager.SetCredentialPassword(hostID, args[1])
	subs := map[string]interface{}{"HostID": hostID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to set password for host {{.HostID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Password is updated for host {{.HostID}}.", subs))
	return nil
}
