package firewall

import (
	"fmt"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AddCommand struct {
	UI              terminal.UI
	FirewallManager managers.FirewallManager
}

func NewAddCommand(ui terminal.UI, firewallManager managers.FirewallManager) (cmd *AddCommand) {
	return &AddCommand{
		UI:              ui,
		FirewallManager: firewallManager,
	}
}

func FirewallAddMetaData() cli.Command {
	return cli.Command{
		Category:    "firewall",
		Name:        "add",
		Description: T("Create a new firewall"),
		Usage:       "${COMMAND_NAME} sl firewall add TARGET [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Usage: T("Firewall type  [required]. Options are: vlan,vs,hardware"),
			},
			cli.BoolFlag{
				Name:  "ha,high-availability",
				Usage: T("High available firewall option"),
			},
			metadata.ForceFlag(),
		},
	}
}

func (cmd *AddCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	if !c.IsSet("type") {
		return errors.NewMissingInputError("--type")
	}
	target, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Target ID")
	}
	firewallType := c.String("type")

	var packages []datatypes.Product_Item
	if firewallType == "vlan" {
		packages, err = cmd.FirewallManager.GetDedicatedPackage(c.IsSet("ha"))
	} else if firewallType == "vs" {
		packages, err = cmd.FirewallManager.GetStandardPackage(target, true)
	} else if firewallType == "hardware" {
		packages, err = cmd.FirewallManager.GetStandardPackage(target, false)
	}
	if err != nil {
		return cli.NewExitError(T("Failed to get package for {{.Type}} firewall.\n", map[string]interface{}{"Type": firewallType})+err.Error(), 2)
	}
	if len(packages) == 0 {
		return cli.NewExitError(T("Failed to find package for firewall."), 2)
	}

	cmd.UI.Print(fmt.Sprintf("Product: %s", *packages[0].Description))
	cmd.UI.Print(fmt.Sprintf("Price: $%.2f monthly", float64(*packages[0].Prices[0].RecurringFee)))

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This action  will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	var orderReciept datatypes.Container_Product_Order_Receipt
	if firewallType == "vlan" {
		orderReciept, err = cmd.FirewallManager.AddVlanFirewall(target, c.IsSet("ha"))
	} else if firewallType == "vs" {
		orderReciept, err = cmd.FirewallManager.AddStandardFirewall(target, true)
	} else if firewallType == "hardware" {
		orderReciept, err = cmd.FirewallManager.AddStandardFirewall(target, false)
	}
	if err != nil {
		return cli.NewExitError(T("Failed to create firewall.\n")+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.ID}} was placed to create a firewall.", map[string]interface{}{"ID": *orderReciept.OrderId}))
	return nil
}
