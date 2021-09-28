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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7PolicyEditCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PolicyEditCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PolicyEditCommand) {
	return &L7PolicyEditCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PolicyEditCommand) Run(c *cli.Context) error {

	policyID := c.Int("policy-id")
	if policyID == 0 {
		return errors.NewMissingInputError("--policy-id")
	}

	name := c.String("n")
	if utils.IsEmptyString(name) {
		return bxErr.NewMissingInputError("-n, --name")
	}

	action := c.String("a")
	if utils.IsEmptyString(action) {
		return bxErr.NewMissingInputError("-a, --action")
	}
	actionUpperCase := strings.ToUpper(action)

	if !IsValidAction(actionUpperCase) {
		return bxErr.NewInvalidUsageError(
			T("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	redirect := c.String("r")
	if !utils.IsEmptyString(redirect) && actionUpperCase == REJECT {
		return bxErr.NewInvalidUsageError(
			T("-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	if IsValidAction(actionUpperCase) && utils.IsEmptyString(redirect) && actionUpperCase != REJECT {
		return bxErr.NewInvalidUsageError(
			T("-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	priority := c.Int("p")

	policy, err := cmd.LoadBalancerManager.GetL7Policy(policyID)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	policy.Name = &name
	policy.Action = &actionUpperCase

	if actionUpperCase != REDIRECT_HTTPS {
		policy.Priority = &priority
	}

	if actionUpperCase == REDIRECT_POOL {
		policy.RedirectL7PoolUuid = &redirect
	}
	if actionUpperCase == REDIRECT_URL || actionUpperCase == REDIRECT_HTTPS {
		policy.RedirectUrl = &redirect
	}

	_, err = cmd.LoadBalancerManager.EditL7Policy(policyID, &policy)
	if err != nil {
		return cli.NewExitError(T("Failed to edit l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy edited"))

	policyEdited, err := cmd.LoadBalancerManager.GetL7Policy(policyID)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policy details: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	PrintPolicies([]datatypes.Network_LBaaS_L7Policy{policyEdited}, cmd.UI)
	return nil
}
