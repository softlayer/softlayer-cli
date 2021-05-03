package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	bxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type L7RuleAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7RuleAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7RuleAddCommand) {
	return &L7RuleAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7RuleAddCommand) Run(c *cli.Context) error {
	policyUUID := c.String("policy-uuid")
	if policyUUID == "" {
		return errors.NewMissingInputError("--policy-uuid")
	}

	policyType := c.String("t")
	if policyType == "" {
		return bxErr.NewMissingInputError("-t, --type")
	}
	if strings.ToUpper(policyType) != "HOST_NAME" && strings.ToUpper(policyType) != "FILE_TYPE" && strings.ToUpper(policyType) != "HEADER" && strings.ToUpper(policyType) != "COOKIE" && strings.ToUpper(policyType) != "PATH" {
		return bxErr.NewInvalidUsageError(T("The value of option -t, --type should be HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH."))
	}

	compareType := c.String("c")
	if compareType == "" {
		return bxErr.NewMissingInputError("-c, --compare-type")
	}
	if strings.ToUpper(compareType) != "EQUAL_TO" && strings.ToUpper(compareType) != "ENDS_WITH" && strings.ToUpper(compareType) != "STARTS_WITH" && strings.ToUpper(compareType) != "REGEX" && strings.ToUpper(compareType) != "CONTAINS" {
		return bxErr.NewInvalidUsageError(T("The value of option -c, --compare-type should be EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS."))
	}

	value := c.String("v")
	if value == "" {
		return bxErr.NewMissingInputError("-v, --value")
	}

	key := c.String("k")

	if key != "" && (strings.ToUpper(policyType) != "HEADER" && strings.ToUpper(policyType) != "COOKIE") {
		return bxErr.NewInvalidUsageError(T("-k, --key is only available in HEADER or COOKIE type."))
	}

	invert := c.Int("invert")

	rule := datatypes.Network_LBaaS_L7Rule{
		Type:           &policyType,
		ComparisonType: &compareType,
		Value:          &value,
		Invert:         &invert,
	}
	if c.IsSet("k") {
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
