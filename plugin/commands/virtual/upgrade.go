package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
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

func VSUpgradeMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "upgrade",
		Description: T("Upgrade a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs upgrade IDENTIFIER [OPTIONS]
	Note: Classic infrastructure service automatically reboots the instance once upgrade request is
  	placed. The instance is halted until the upgrade transaction is completed.
  	However for Network, no reboot is required.

EXAMPLE:
   ${COMMAND_NAME} sl vs upgrade 12345678 -c 8 -m 8192 --network 1000
   This commands upgrades virtual server instance with ID 12345678 and set number of CPU cores to 8, memory to 8192M, network port speed to 1000 Mbps.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Number of CPU cores"),
			},
			cli.BoolFlag{
				Name:  "private",
				Usage: T("CPU core will be on a dedicated host server"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Memory in megabytes"),
			},
			cli.IntFlag{
				Name:  "network",
				Usage: T("Network port speed in Mbps"),
			},
			cli.StringFlag{
				Name:  "flavor",
				Usage: T("Flavor key name"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}
