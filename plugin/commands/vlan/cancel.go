package vlan

import (
	"strconv"
	"strings"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

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
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vlanID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("VLAN ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the VLAN: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": vlanID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	// See if the API will just tell us if this VLAN can't be cancelled for a specific reason
	reasons := cmd.NetworkManager.GetCancelFailureReasons(vlanID)
	if len(reasons) > 0 {
		for _, reason := range reasons {
			cmd.UI.Print(reason)
		}
		return cli.NewExitError(T("Failed to cancel VLAN {{.ID}}.\n", map[string]interface{}{"ID": vlanID}), 2)

	}
	err = cmd.NetworkManager.CancelVLAN(vlanID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find VLAN with ID {{.ID}}.\n", map[string]interface{}{"ID": vlanID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel VLAN {{.ID}}.\n", map[string]interface{}{"ID": vlanID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("VLAN {{.ID}} was cancelled.", map[string]interface{}{"ID": vlanID}))
	return nil
}

func VlanCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "vlan",
		Name:        "cancel",
		Description: T("Cancel a VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan cancel IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan cancel 12345678 -f
   This command cancels vlan with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
