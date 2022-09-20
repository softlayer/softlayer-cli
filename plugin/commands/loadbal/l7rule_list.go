package loadbal

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7RuleListCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PolicyId            int
}

func NewL7RuleListCommand(sl *metadata.SoftlayerCommand) *L7RuleListCommand {
	thisCmd := &L7RuleListCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7rules",
		Short: T("List l7 rules"),
		Long:  T("${COMMAND_NAME} sl loadbal l7rules (--policy-id Policy_ID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.PolicyId, "policy-id", 0, T("ID for the load balancer policy [required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7RuleListCommand) Run(args []string) error {
	policyid := cmd.PolicyId
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
