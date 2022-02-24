package subnet

import (
	"strconv"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
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
	subnetID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the subnet: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": subnetID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.CancelSubnet(subnetID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find subnet with ID: {{.ID}}.\n", map[string]interface{}{"ID": subnetID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Subnet {{.ID}} was cancelled.", map[string]interface{}{"ID": subnetID}))
	return nil
}

func SubnetCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "subnet",
		Name:        "cancel",
		Description: T("Cancel a subnet"),
		Usage: T(`${COMMAND_NAME} sl subnet cancel IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet cancel 12345678 -f
   This command cancels subnet with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
