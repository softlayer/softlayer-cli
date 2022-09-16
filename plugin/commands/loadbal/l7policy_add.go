package loadbal

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	REJECT         = "REJECT"
	REDIRECT_POOL  = "REDIRECT_POOL"
	REDIRECT_URL   = "REDIRECT_URL"
	REDIRECT_HTTPS = "REDIRECT_HTTPS"
)

type L7PolicyAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	ProtocolUuid        string
	Name                string
	Action              string
	Redirect            string
	Priority            int
}

func NewL7PolicyAddCommand(sl *metadata.SoftlayerCommand) *L7PolicyAddCommand {
	thisCmd := &L7PolicyAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7policy-add",
		Short: T("Add a new L7 policy."),
		Long:  T("${COMMAND_NAME} sl loadbal l7policy-add (--protocol-uuid PROTOCOL_UUID) (-n, --name NAME) (-a,--action REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS) [-r,--redirect REDIRECT] [-p,--priority PRIORITY]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.ProtocolUuid, "protocol-uuid", "", T("UUID for the load balancer protocol [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Policy name"))
	cobraCmd.Flags().StringVarP(&thisCmd.Action, "action", "a", "", T("Policy action: REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
	cobraCmd.Flags().StringVarP(&thisCmd.Redirect, "redirect", "r", "", T("POOL_UUID, URL or HTTPS_PROTOCOL_UUID . It's only available in REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS action"))
	cobraCmd.Flags().IntVarP(&thisCmd.Priority, "priority", "p", 1, T("Policy priority"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PolicyAddCommand) Run(args []string) error {

	policyUUID := cmd.ProtocolUuid
	if utils.IsEmptyString(policyUUID) {
		return errors.NewMissingInputError("--protocol-uuid")
	}

	name := cmd.Name
	if utils.IsEmptyString(name) {
		return bxErr.NewMissingInputError("-n, --name")
	}

	action := cmd.Action
	if utils.IsEmptyString(action) {
		return bxErr.NewMissingInputError("-a, --action")
	}
	actionUpperCase := strings.ToUpper(action)

	if !IsValidAction(actionUpperCase) {
		return bxErr.NewInvalidUsageError(
			T("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	redirect := cmd.Redirect
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

	priority := cmd.Priority

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
