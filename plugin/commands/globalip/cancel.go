package globalip

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewCancelCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	globalIPID, err := utils.ResolveGloablIPId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Globalip ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the IP address: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": globalIPID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.NetworkManager.CancelGlobalIP(globalIPID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find global IP with ID: {{.ID}}.\n", map[string]interface{}{"ID": globalIPID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel global IP: {{.ID}}.\n", map[string]interface{}{"ID": globalIPID})+err.Error(), 2)

	}

	cmd.UI.Ok()
	cmd.UI.Print(T("IP address {{.ID}} was cancelled.", map[string]interface{}{"ID": globalIPID}))
	return nil
}

func GlobalIpCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "globalip",
		Name:        "cancel",
		Description: T("Cancel a global IP"),
		Usage: T(`${COMMAND_NAME} sl globalip cancel IDENTIFIER [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip cancel 12345678
	This command cancels IP address with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
