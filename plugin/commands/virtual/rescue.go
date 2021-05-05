package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RescueCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewRescueCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *RescueCommand) {
	return &RescueCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *RescueCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will reboot virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.VirtualServerManager.RescueInstance(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to rescue virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was rebooted.", map[string]interface{}{"VsId": vsID}))
	return nil
}
