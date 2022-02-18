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

	policyId := c.Int("policy-id")
	if policyId == 0 {
		return errors.NewMissingInputError("--policy-id")
	}

	currentPolicy, err := cmd.LoadBalancerManager.GetL7Policy(policyId)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	name := c.String("n")
	if !utils.IsEmptyString(name) {
		currentPolicy.Name = &name
	}

	action := c.String("a")
	actionToUpdate := strings.ToUpper(action)

	if !utils.IsEmptyString(actionToUpdate) && !IsValidAction(actionToUpdate) {
		return bxErr.NewInvalidUsageError(
			T("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	if !utils.IsEmptyString(actionToUpdate) && IsValidAction(actionToUpdate) {
		currentPolicy.Action = &actionToUpdate
	}

	redirect := c.String("r")
	if !utils.IsEmptyString(redirect) && utils.FormatStringPointer(currentPolicy.Action) == REJECT {
		return bxErr.NewInvalidUsageError(
			T("-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	if IsValidAction(actionToUpdate) && utils.IsEmptyString(redirect) && actionToUpdate != REJECT {
		return bxErr.NewInvalidUsageError(
			T("-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	priority := c.Int("p")

	if priority > 0 && utils.FormatStringPointer(currentPolicy.Action) != REDIRECT_HTTPS {
		currentPolicy.Priority = &priority
	}

	if utils.FormatStringPointer(currentPolicy.Action) == REDIRECT_POOL {
		currentPolicy.RedirectL7PoolUuid = &redirect
	}
	if utils.FormatStringPointer(currentPolicy.Action) == REDIRECT_URL || utils.FormatStringPointer(currentPolicy.Action) == REDIRECT_HTTPS {
		currentPolicy.RedirectUrl = &redirect
	}

	_, err = cmd.LoadBalancerManager.EditL7Policy(policyId, &currentPolicy)
	if err != nil {
		return cli.NewExitError(T("Failed to edit l7 policy: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy edited"))

	policyEdited, err := cmd.LoadBalancerManager.GetL7Policy(policyId)
	if err != nil {
		return cli.NewExitError(T("Failed to get l7 policy detail: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	PrintPolicies([]datatypes.Network_LBaaS_L7Policy{policyEdited}, cmd.UI)
	return nil
}

func LoadbalL7PolicyEditMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7policy-edit",
		Description: T("Edit a L7 policy"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policy-edit (--policy-d POLICY_ID) (-n, --name NAME) (-a,--action REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS) [-r,--redirect REDIRECT] [-p,--priority PRIORITY]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "policy-id",
				Usage: T("ID for the load balancer policy [required]"),
			},
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Policy name"),
			},
			cli.StringFlag{
				Name:  "a,action",
				Usage: T("Policy action: REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
			},
			cli.StringFlag{
				Name:  "r,redirect",
				Usage: T("POOL_UUID, URL or HTTPS_PROTOCOL_UUID . It's only available in REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS action"),
			},
			cli.IntFlag{
				Name:  "p,priority",
				Usage: T("Policy priority"),
			},
		},
	}
}
