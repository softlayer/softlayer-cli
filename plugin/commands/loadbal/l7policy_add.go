package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type L7PolicyAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PolicyAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PolicyAddCommand) {
	return &L7PolicyAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PolicyAddCommand) Run(c *cli.Context) error {
	policyUUID := c.String("protocol-uuid")
	if policyUUID == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}

	name := c.String("n")
	if name == "" {
		return bxErr.NewMissingInputError("-n, --name")
	}

	action := c.String("a")
	if action == "" {
		return bxErr.NewMissingInputError("-a, --action")
	}
	actionUpperCase := strings.ToUpper(action)

	if actionUpperCase != "REJECT" && actionUpperCase != "REDIRECT_POOL" && actionUpperCase != "REDIRECT_URL" && actionUpperCase != "REDIRECT_HTTPS" {
		return bxErr.NewInvalidUsageError(T("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
	}

	redirect := c.String("r")
	if redirect != "" && actionUpperCase == "REJECT" {
		return bxErr.NewInvalidUsageError(T("-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL"))
	}

	if (actionUpperCase == "REDIRECT_POOL" || actionUpperCase == "REDIRECT_URL" || actionUpperCase == "REDIRECT_HTTPS") && redirect == "" {
		return bxErr.NewInvalidUsageError(T("-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
	}

	priority := c.Int("p")

	policy := datatypes.Network_LBaaS_L7Policy{
		Name:   &name,
		Action: &actionUpperCase,
	}

	if strings.ToUpper(actionUpperCase) != "REDIRECT_HTTPS" {
		policy.Priority = &priority
	}

	if strings.ToUpper(actionUpperCase) == "REDIRECT_POOL" {
		policy.RedirectL7PoolUuid = &redirect
	}
	if strings.ToUpper(actionUpperCase) == "REDIRECT_URL" || strings.ToUpper(actionUpperCase) == "REDIRECT_HTTPS" {
		policy.RedirectUrl = &redirect
	}

	policyRule := datatypes.Network_LBaaS_PolicyRule{
		L7Policy: &policy,
	}

	_, err := cmd.LoadBalancerManager.AddL7Policy(&policyUUID, policyRule)
	if err != nil {
		return cli.NewExitError(T("Failed to add l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy added"))
	return nil
}
