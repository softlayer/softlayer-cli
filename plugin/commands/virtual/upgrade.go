package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UpgradeCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewUpgradeCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *UpgradeCommand) {
	return &UpgradeCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *UpgradeCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if c.IsSet("private") && !c.IsSet("cpu") {
		return bmxErr.NewInvalidUsageError(T("Must specify [--cpu] when using [--private]."))
	}

	if !c.IsSet("cpu") && !c.IsSet("memory") && !c.IsSet("network") && !c.IsSet("flavor") {
		return bmxErr.NewInvalidUsageError(T("Must provide [--cpu], [--memory], [--network] or [--flavor] to upgrade."))
	}

	if c.IsSet("flavor") && (c.IsSet("cpu") || c.IsSet("memory") || c.IsSet("private")) {
		return bmxErr.NewInvalidUsageError(T("Option [--flavor] is exclusive with [--cpu], [--memory] and [--private]."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("f") && !c.IsSet("force") && outputFormat != "JSON" {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	orderReceipt, err := cmd.VirtualServerManager.UpgradeInstance(vsID, c.Int("cpu"), c.Int("memory"), c.Int("network"), c.Bool("private"), c.String("flavor"))
	if err != nil {
		return cli.NewExitError(T("Failed to upgrade virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderId}} to upgrade virtual server instance: {{.VsId}} was placed.",
		map[string]interface{}{"OrderId": *orderReceipt.OrderId, "VsId": vsID}))

	return nil
}
