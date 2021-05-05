package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7PolicyListCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PolicyListCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PolicyListCommand) {
	return &L7PolicyListCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PolicyListCommand) Run(c *cli.Context) error {
	protocolID := c.Int("protocol-id")
	if protocolID == 0 {
		return errors.NewMissingInputError("--protocol-id")
	}

	l7Policies, err := cmd.LoadBalancerManager.GetL7Policies(protocolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policies: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	if len(l7Policies) == 0 {
		cmd.UI.Say(T("No l7 policies was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("UUID"), T("Name"), T("Action"), T("Redirect"), T("Priority")})
		for _, l7Policy := range l7Policies {
			if l7Policy.Action != nil && *l7Policy.Action == "REDIRECT_URL" {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatStringPointer(l7Policy.RedirectUrl),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == "REDIRECT_POOL" {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatIntPointer(l7Policy.RedirectL7PoolId),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == "REJECT" {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					"-",
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
		}
		table.Print()
	}
	return nil
}
