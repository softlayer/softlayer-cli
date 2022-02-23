package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
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
	PrintPolicies(l7Policies, cmd.UI)

	return nil
}

func PrintPolicies(l7Policies []datatypes.Network_LBaaS_L7Policy, cmdUI terminal.UI) {

	if len(l7Policies) == 0 {
		cmdUI.Say(T("No l7 policies was found."))
	} else {
		table := cmdUI.Table([]string{T("ID"), T("UUID"), T("Name"), T("Action"), T("Redirect"), T("Priority"), T("Create Date")})
		for _, l7Policy := range l7Policies {
			if l7Policy.Action != nil && (*l7Policy.Action == REDIRECT_URL || *l7Policy.Action == REDIRECT_HTTPS) {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatStringPointer(l7Policy.RedirectUrl),
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == REDIRECT_POOL {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					utils.FormatIntPointer(l7Policy.RedirectL7PoolId),
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
			if l7Policy.Action != nil && *l7Policy.Action == REJECT {
				table.Add(utils.FormatIntPointer(l7Policy.Id),
					utils.FormatStringPointer(l7Policy.Uuid),
					utils.FormatStringPointer(l7Policy.Name),
					utils.FormatStringPointer(l7Policy.Action),
					"-",
					utils.FormatIntPointer(l7Policy.Priority),
					utils.FormatSLTimePointer(l7Policy.CreateDate),
				)
			}
		}
		table.Print()
	}
}

func LoadbalL7PolicyListMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7policies",
		Description: T("List L7 policies"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policies (--protocol-id PROTOCOL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "protocol-id",
				Usage: T("ID for the load balancer protocol [required]"),
			},
		},
	}
}
