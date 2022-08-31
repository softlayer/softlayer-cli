package virtual

import (
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCancelCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	VsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the virtual server instance: {{.VsID}} and cannot be undone. Continue?", map[string]interface{}{"VsID": VsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.VirtualServerManager.CancelInstance(VsID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find virtual server instance with ID: {{.VsID}}.\n", map[string]interface{}{"VsID": VsID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": VsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was cancelled.", map[string]interface{}{"VsId": VsID}))
	return nil
}

func VSCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "cancel",
		Description: T("Cancel virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs cancel IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs cancel 12345678
   This command cancels virtual server instance with ID of 12345678.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
