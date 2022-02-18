package firewall

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	UI              terminal.UI
	FirewallManager managers.FirewallManager
}

func NewCancelCommand(ui terminal.UI, firewallManager managers.FirewallManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:              ui,
		FirewallManager: firewallManager,
	}
}

func FirewallCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "firewall",
		Name:        "cancel",
		Description: T("Cancels a firewall"),
		Usage:       "${COMMAND_NAME} sl firewall cancel IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(c.Args()[0])
	if err != nil {
		return cli.NewExitError(T("Failed to parse firewall ID : {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": c.Args()[0]})+err.Error(), 1)
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This action will cancel the firewall {{.ID}} from your account. Continue?", map[string]interface{}{"ID": c.Args()[0]}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
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
			return cli.NewExitError(T("Unable to find firewall with ID: {{.ID}}", map[string]interface{}{"ID": c.Args()[0]})+err.Error(), 2)
		}
		return cli.NewExitError(T("Failed to cancel firewall: {{.ID}}", map[string]interface{}{"ID": c.Args()[0]})+err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Firewall {{.ID}} is being cancelled!", map[string]interface{}{"ID": c.Args()[0]}))
	return nil
}
