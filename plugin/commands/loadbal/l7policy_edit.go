package loadbal

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7PolicyEditCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PolicyId            int
	Name                string
	Action              string
	Redirect            string
	Priority            int
}

func NewL7PolicyEditCommand(sl *metadata.SoftlayerCommand) *L7PolicyEditCommand {
	thisCmd := &L7PolicyEditCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7policy-edit",
		Short: T("Edit a L7 policy"),
		Long:  T("${COMMAND_NAME} sl loadbal l7policy-edit (--policy-d POLICY_ID) (-n, --name NAME) (-a,--action REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS) [-r,--redirect REDIRECT] [-p,--priority PRIORITY]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.PolicyId, "policy-id", 0, T("ID for the load balancer policy [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Policy name"))
	cobraCmd.Flags().StringVarP(&thisCmd.Action, "action", "a", "", T("Policy action: REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
	cobraCmd.Flags().StringVarP(&thisCmd.Redirect, "redirect", "r", "", T("POOL_UUID, URL or HTTPS_PROTOCOL_UUID . It's only available in REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS action"))
	cobraCmd.Flags().IntVarP(&thisCmd.Priority, "priority", "p", 0, T("Policy priority"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PolicyEditCommand) Run(args []string) error {

	policyId := cmd.PolicyId
	if policyId == 0 {
		return errors.NewMissingInputError("--policy-id")
	}

	currentPolicy, err := cmd.LoadBalancerManager.GetL7Policy(policyId)
	if err != nil {
		return errors.New(T("Failed to get l7 policy: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}

	name := cmd.Name
	if !utils.IsEmptyString(name) {
		currentPolicy.Name = &name
	}

	action := cmd.Action
	actionToUpdate := strings.ToUpper(action)

	if !utils.IsEmptyString(actionToUpdate) && !IsValidAction(actionToUpdate) {
		return errors.NewInvalidUsageError(
			T("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	if !utils.IsEmptyString(actionToUpdate) && IsValidAction(actionToUpdate) {
		currentPolicy.Action = &actionToUpdate
	}

	redirect := cmd.Redirect
	if !utils.IsEmptyString(redirect) && utils.FormatStringPointer(currentPolicy.Action) == REJECT {
		return errors.NewInvalidUsageError(
			T("-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	if IsValidAction(actionToUpdate) && utils.IsEmptyString(redirect) && actionToUpdate != REJECT {
		return errors.NewInvalidUsageError(
			T("-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
		)
	}

	priority := cmd.Priority

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
		return errors.New(T("Failed to edit l7 policy: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}

	cmd.UI.Ok()
	cmd.UI.Say(T("L7 policy edited"))

	policyEdited, err := cmd.LoadBalancerManager.GetL7Policy(policyId)
	if err != nil {
		return errors.New(T("Failed to get l7 policy detail: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}
	PrintPolicies([]datatypes.Network_LBaaS_L7Policy{policyEdited}, cmd.UI)
	return nil
}
