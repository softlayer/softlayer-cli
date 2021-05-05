package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReloadCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	Context              plugin.PluginContext
}

func NewReloadCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, context plugin.PluginContext) (cmd *ReloadCommand) {
	return &ReloadCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		Context:              context,
	}
}

func (cmd *ReloadCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will reload operating system of virtual server instance: {{.VsId}} and cannot be undone. Continue?",
			map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.ReloadInstance(vsID, c.String("i"), c.IntSlice("k"), c.Int("image"))
	if err != nil {
		return cli.NewExitError(T("Failed to reload virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("System reloading for virtual server instance: {{.VsId}} is in progress. Run '{{.CommandName}} sl vs ready {{.VsId}}' to check whether it is ready later on.",
		map[string]interface{}{
			"VsId":        vsID,
			"CommandName": cmd.Context.CLIName()}))
	return nil
}
