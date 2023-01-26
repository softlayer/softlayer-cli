package block

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type AccessPasswordCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Password       string
}

func NewAccessPasswordCommand(sl *metadata.SoftlayerStorageCommand) *AccessPasswordCommand {
	thisCmd := &AccessPasswordCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "access-password " + T("IDENTIFIER"),
		Short: T("Changes a password for a volume's access."),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} access-password ACCESS_ID --password PASSWORD
	
	ACCESS_ID is the allowed_host_id from '${COMMAND_NAME} sl {{.storageType}} access-list'`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
		DisableFlagsInUseLine: true,
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Password, "password", "p", "", T("Password you want to set, this command will fail if the password is not strong. [required]"))
	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *AccessPasswordCommand) Run(args []string) error {

	fmt.Printf("===AccessPasswordCommand===")
	hostID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("allowed access host ID")
	}

	if cmd.Password == "" {
		return slErr.NewInvalidUsageError(T("[-p|--password] is required."))
	}

	err = cmd.StorageManager.SetCredentialPassword(hostID, cmd.Password)
	subs := map[string]interface{}{"HostID": hostID}
	if err != nil {
		return slErr.NewAPIError(T("Failed to set password for host {{.HostID}}.\n", subs), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Password is updated for host {{.HostID}}.", subs))
	return nil
}
