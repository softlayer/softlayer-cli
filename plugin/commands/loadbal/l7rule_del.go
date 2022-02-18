package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7RuleDelCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7RuleDelCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7RuleDelCommand) {
	return &L7RuleDelCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7RuleDelCommand) Run(c *cli.Context) error {
	l7PolicyID := c.String("policy-uuid")
	if l7PolicyID == "" {
		return errors.NewMissingInputError("--policy-uuid")
	}

	l7RuleID := c.String("rule-uuid")
	if l7RuleID == "" {
		return errors.NewMissingInputError("--rule-uuid")
	}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer L7 rule: {{.RuleID}} and cannot be undone. Continue?", map[string]interface{}{"RuleID": l7RuleID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.LoadBalancerManager.DeleteL7Rule(&l7PolicyID, l7RuleID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete L7Rule {{.L7RuleID}}: {{.Error}}.\n",
			map[string]interface{}{"L7RuleID": l7RuleID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7Rule {{.L7RuleID}} removed", map[string]interface{}{"L7RuleID": l7RuleID}))
	return nil
}

func LoadbalL7RuleDelMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7rule-delete",
		Description: T("Delete a L7 rule"),
		Usage:       "${COMMAND_NAME} sl loadbal l7rule-delete (--policy-uuid L7POLICY_UUID) (--rule-uuid L7RULE_UUID) [-f, --force]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "policy-uuid",
				Usage: T("UUID for the load balancer policy [required]"),
			},
			cli.StringFlag{
				Name:  "rule-uuid",
				Usage: T("UUID for the load balancer rule [required]"),
			},
			metadata.ForceFlag(),
		},
	}
}
