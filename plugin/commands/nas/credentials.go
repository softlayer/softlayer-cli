package nas

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialsCommand struct {
	*metadata.SoftlayerCommand
	NasNetworkStorageManager managers.NasNetworkStorageManager
	Command                  *cobra.Command
}

func NewCredentialsCommand(sl *metadata.SoftlayerCommand) *CredentialsCommand {
	thisCmd := &CredentialsCommand{
		SoftlayerCommand:         sl,
		NasNetworkStorageManager: managers.NewNasNetworkStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "credentials " + T("IDENTIFIER"),
		Short: T("List NAS account credentials."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialsCommand) Run(args []string) error {

	nasNetworkStorageId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[id,username,password]"
	nasNetworkStorage, err := cmd.NasNetworkStorageManager.GetNasNetworkStorage(nasNetworkStorageId, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get NAS Network Storage.")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Username"), T("Password")})

	password := "-"
	if nasNetworkStorage.Password != nil {
		password = *nasNetworkStorage.Password
	}
	table.Add(
		utils.FormatStringPointer(nasNetworkStorage.Username),
		password,
	)

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
