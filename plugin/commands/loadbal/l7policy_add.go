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

const (
	REJECT         = "REJECT"
	REDIRECT_POOL  = "REDIRECT_POOL"
	REDIRECT_URL   = "REDIRECT_URL"
	REDIRECT_HTTPS = "REDIRECT_HTTPS"
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
	if utils.IsEmptyString(policyUUID) {
		return errors.NewMissingInputError("--protocol-uuid")
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

	policy := datatypes.Network_LBaaS_L7Policy{
		Name:   &name,
		Action: &actionUpperCase,
	}

	if actionUpperCase != REDIRECT_HTTPS {
		policy.Priority = &priority
	}

	if actionUpperCase == REDIRECT_POOL {
		policy.RedirectL7PoolUuid = &redirect
	}
	if actionUpperCase == REDIRECT_URL || actionUpperCase == REDIRECT_HTTPS {
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

func IsValidAction(action string) bool {
	switch action {
	case REJECT, REDIRECT_URL, REDIRECT_POOL, REDIRECT_HTTPS:
		return true
	}
	return false
}

func LoadbalL7PolicyAddMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7policy-add",
		Description: T("Add a new L7 policy"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policy-add (--protocol-uuid PROTOCOL_UUID) (-n, --name NAME) (-a,--action REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS) [-r,--redirect REDIRECT] [-p,--priority PRIORITY]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "protocol-uuid",
				Usage: T("UUID for the load balancer protocol [required]"),
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
				Value: 1,
			},
		},
	}
}
