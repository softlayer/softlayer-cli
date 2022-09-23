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
)

type L7RuleAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PolicyUuid          string
	Type                string
	CompareType         string
	Value               string
	Key                 string
	Invert              int
}

func NewL7RuleAddCommand(sl *metadata.SoftlayerCommand) *L7RuleAddCommand {
	thisCmd := &L7RuleAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7rule-add",
		Short: T("Add a new L7 rule"),
		Long:  T("${COMMAND_NAME} sl loadbal l7rule-add (--policy-uuid L7POLICY_UUID) (-t, --type HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH ) (-c, --compare-type EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS) (-v,--value VALUE) [-k,--key KEY] [--invert 0 | 1]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.PolicyUuid, "policy-uuid", "", T("UUID for the load balancer policy [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Type, "type", "t", "", T("Rule type: HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH. [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.CompareType, "compare-type", "c", "", T("Compare type: EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS. [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Value, "value", "v", "", T("Compared Value [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Key, "key", "k", "", T("Key name. It's only available in HEADER or COOKIE type"))
	cobraCmd.Flags().IntVar(&thisCmd.Invert, "invert", 0, T("Invert rule: 0 | 1."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7RuleAddCommand) Run(args []string) error {
	policyUUID := cmd.PolicyUuid
	if policyUUID == "" {
		return errors.NewMissingInputError("--policy-uuid")
	}

	policyType := cmd.Type
	if policyType == "" {
		return bxErr.NewMissingInputError("-t, --type")
	}
	if strings.ToUpper(policyType) != "HOST_NAME" && strings.ToUpper(policyType) != "FILE_TYPE" && strings.ToUpper(policyType) != "HEADER" && strings.ToUpper(policyType) != "COOKIE" && strings.ToUpper(policyType) != "PATH" {
		return bxErr.NewInvalidUsageError(T("The value of option -t, --type should be HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH."))
	}

	compareType := cmd.CompareType
	if compareType == "" {
		return bxErr.NewMissingInputError("-c, --compare-type")
	}
	if strings.ToUpper(compareType) != "EQUAL_TO" && strings.ToUpper(compareType) != "ENDS_WITH" && strings.ToUpper(compareType) != "STARTS_WITH" && strings.ToUpper(compareType) != "REGEX" && strings.ToUpper(compareType) != "CONTAINS" {
		return bxErr.NewInvalidUsageError(T("The value of option -c, --compare-type should be EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS."))
	}

	value := cmd.Value
	if value == "" {
		return bxErr.NewMissingInputError("-v, --value")
	}

	key := cmd.Key

	if key != "" && (strings.ToUpper(policyType) != "HEADER" && strings.ToUpper(policyType) != "COOKIE") {
		return bxErr.NewInvalidUsageError(T("-k, --key is only available in HEADER or COOKIE type."))
	}

	invert := cmd.Invert

	rule := datatypes.Network_LBaaS_L7Rule{
		Type:           &policyType,
		ComparisonType: &compareType,
		Value:          &value,
		Invert:         &invert,
	}
	if cmd.Key != "" {
		rule.Key = &key
	}

	_, err := cmd.LoadBalancerManager.AddL7Rule(&policyUUID, rule)
	if err != nil {
		return cli.NewExitError(T("Failed to add l7 rule: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 rule added"))
	return nil
}
