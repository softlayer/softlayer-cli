package firewall

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AddCommand struct {
	*metadata.SoftlayerCommand
	FirewallManager  managers.FirewallManager
	Command          *cobra.Command
	Type             string
	HighAvailability bool
	Force            bool
}

func NewAddCommand(sl *metadata.SoftlayerCommand) (cmd *AddCommand) {
	thisCmd := &AddCommand{
		SoftlayerCommand: sl,
		FirewallManager:  managers.NewFirewallManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "add " + T("TargetID"),
		Short: T("Create a new firewall."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Type, "type", "", T("Firewall type  [required]. Options are: vlan,vs,hardware"))
	cobraCmd.Flags().BoolVar(&thisCmd.HighAvailability, "high-availability", false, T("High available firewall option"))
	cobraCmd.Flags().BoolVar(&thisCmd.Force, "force", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AddCommand) Run(args []string) error {

	if cmd.Type == "" {
		return errors.NewMissingInputError("--type")
	}
	target, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Target ID")
	}
	firewallType := cmd.Type

	var packages []datatypes.Product_Item
	if firewallType == "vlan" {
		packages, err = cmd.FirewallManager.GetDedicatedPackage(cmd.HighAvailability)
	} else if firewallType == "vs" {
		packages, err = cmd.FirewallManager.GetStandardPackage(target, true)
	} else if firewallType == "hardware" {
		packages, err = cmd.FirewallManager.GetStandardPackage(target, false)
	}
	if err != nil {
		return errors.NewAPIError(T("Failed to get package for {{.Type}} firewall.", map[string]interface{}{"Type": firewallType}), err.Error(), 2)
	}
	if len(packages) == 0 {
		return cli.NewExitError(T("Failed to find package for firewall."), 2)
	}

	cmd.UI.Print(fmt.Sprintf("Product: %s", *packages[0].Description))
	cmd.UI.Print(fmt.Sprintf("Price: $%.2f monthly", float64(*packages[0].Prices[0].RecurringFee)))

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action  will incur charges on your account. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	var orderReciept datatypes.Container_Product_Order_Receipt
	if firewallType == "vlan" {
		orderReciept, err = cmd.FirewallManager.AddVlanFirewall(target, cmd.HighAvailability)
	} else if firewallType == "vs" {
		orderReciept, err = cmd.FirewallManager.AddStandardFirewall(target, true)
	} else if firewallType == "hardware" {
		orderReciept, err = cmd.FirewallManager.AddStandardFirewall(target, false)
	}
	if err != nil {
		return errors.NewAPIError(T("Failed to create firewall."), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.ID}} was placed to create a firewall.", map[string]interface{}{"ID": *orderReciept.OrderId}))
	return nil
}
