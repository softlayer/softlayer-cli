package nas

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialsCommand struct {
	UI                       terminal.UI
	NasNetworkStorageManager managers.NasNetworkStorageManager
}

func NewCredentialsCommand(ui terminal.UI, nasNetworkStorageManager managers.NasNetworkStorageManager) (cmd *CredentialsCommand) {
	return &CredentialsCommand{
		UI:                       ui,
		NasNetworkStorageManager: nasNetworkStorageManager,
	}
}

func (cmd *CredentialsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	nasNetworkStorageId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

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

func NasCredentialsMetaData() cli.Command {
	return cli.Command{
		Category:    "nas",
		Name:        "credentials",
		Description: T("List NAS account credentials."),
		Usage: T(`${COMMAND_NAME} sl nas credentials IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl nas credentials 123456`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
