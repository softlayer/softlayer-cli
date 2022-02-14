package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type AccessPasswordCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewAccessPasswordCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *AccessPasswordCommand) {
	return &AccessPasswordCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockAccessPasswordMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "access-password",
		Description: T("Changes a password for a volume's access"),
		Usage: T(`${COMMAND_NAME} sl block access-password ACCESS_ID PASSWORD
	
	ACCESS_ID is the allowed_host_id from '${COMMAND_NAME} sl block access-list'
	`),
		Flags: []cli.Flag{},
	}
}

func (cmd *AccessPasswordCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	hostID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("allowed access host ID")
	}
	err = cmd.StorageManager.SetCredentialPassword(hostID, c.Args()[1])
	if err != nil {
		return cli.NewExitError(T("Failed to set password for host {{.HostID}}.\n", map[string]interface{}{"HostID": hostID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Password is updated for host {{.HostID}}.", map[string]interface{}{"HostID": hostID}))
	return nil
}
