package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7RuleListCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7RuleListCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7RuleListCommand) {
	return &L7RuleListCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7RuleListCommand) Run(c *cli.Context) error {
	policyid := c.Int("policy-id")
	if policyid == 0 {
		return bxErr.NewMissingInputError("--policy-id")
	}

	l7Rules, err := cmd.LoadBalancerManager.ListL7Rule(policyid)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 rules: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	if len(l7Rules) == 0 {
		cmd.UI.Say(T("No l7 rules was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("UUID"), T("Type"), T("Compare Type"), T("Value"), T("Key"), T("Invert")})
		for _, l7Rule := range l7Rules {
			table.Add(utils.FormatIntPointer(l7Rule.Id),
				utils.FormatStringPointer(l7Rule.Uuid),
				utils.FormatStringPointer(l7Rule.Type),
				utils.FormatStringPointer(l7Rule.ComparisonType),
				utils.FormatStringPointer(l7Rule.Value),
				utils.FormatStringPointer(l7Rule.Key),
				utils.FormatIntPointer(l7Rule.Invert),
			)
		}
		table.Print()
	}
	return nil
}
