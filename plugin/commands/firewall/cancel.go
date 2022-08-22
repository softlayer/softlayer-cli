package firewall

import (
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	*metadata.SoftlayerCommand
	FirewallManager managers.FirewallManager
	Command         *cobra.Command
	Force           bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) (cmd *CancelCommand) {
	thisCmd := &CancelCommand{
		SoftlayerCommand: sl,
		FirewallManager:  managers.NewFirewallManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancels a firewall."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Force, "force", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {

	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to parse firewall ID : {{.FirewallID}}.", map[string]interface{}{"FirewallID": args[0]}), err.Error(), 1)
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action will cancel the firewall {{.ID}} from your account. Continue?", map[string]interface{}{"ID": args[0]}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	if firewallType == "vlan" {
		err = cmd.FirewallManager.CancelFirewall(firewallID, true)
	} else {
		err = cmd.FirewallManager.CancelFirewall(firewallID, false)
	}
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find firewall with ID: {{.ID}}", map[string]interface{}{"ID": args[0]}), err.Error(), 2)
		}
		return errors.NewAPIError(T("Failed to cancel firewall: {{.ID}}", map[string]interface{}{"ID": args[0]}), err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Firewall {{.ID}} is being cancelled!", map[string]interface{}{"ID": args[0]}))
	return nil
}
