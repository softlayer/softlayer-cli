package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
	// "github.ibm.com/Bluemix/resource-catalog-cli/plugin/i18n"
)

type CreateCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewCreateCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	name := c.String("n")
	if name == "" {
		return errors.NewMissingInputError("-n, --name")
	}

	dataCenter := c.String("d")
	if dataCenter == "" {
		return errors.NewMissingInputError("-d, --datacenter")
	}

	lbType := c.String("t")
	if lbType == "" {
		return errors.NewMissingInputError("-t, --type")
	}
	if strings.ToLower(lbType) != "publictoprivate" && strings.ToLower(lbType) != "privatetoprivate" && strings.ToLower(lbType) != "publictopublic" {
		return errors.NewInvalidUsageError(T("The value of option '-t, --type' should be PublicToPrivate | PrivateToPrivate | PublicToPublic"))
	}
	var lbTypeRequest int
	if strings.ToLower(lbType) == "publictoprivate" {
		lbTypeRequest = 1
	}
	if strings.ToLower(lbType) == "privatetoprivate" {
		lbTypeRequest = 0
	}
	if strings.ToLower(lbType) == "publictopublic" {
		lbTypeRequest = 2
	}

	subnet := c.Int("s")
	if strings.ToLower(lbType) != "publictopublic" && subnet == 0 {
		return errors.NewMissingInputError("-s, --subnet")
	}
	if strings.ToLower(lbType) == "publictopublic" && subnet != 0 {
		return errors.NewInvalidUsageError(T("-s, --subnet is only available in PublicToPrivate and PrivateToPrivate load balancer type."))
	}

	if c.IsSet("use-public-subnet") && strings.ToLower(lbType) != "publictoprivate" {
		return errors.NewInvalidUsageError(T("--use-public-subnet is only available in PublicToPrivate."))
	}

	frontProtocol := c.String("frontend-protocol")
	if frontProtocol == "" {
		frontProtocol = "HTTP"
	}
	frontPort := c.Int("frontend-port")
	if frontPort == 0 {
		frontPort = 80
	}

	backendProtocol := c.String("backend-protocol")
	if backendProtocol == "" {
		backendProtocol = "HTTP"
	}
	backendPort := c.Int("backend-port")
	if backendPort == 0 {
		backendPort = 80
	}

	label := c.String("l")

	method := c.String("m")
	if method == "" {
		method = "ROUNDROBIN"
	}

	if !c.IsSet("f") && !c.IsSet("verify") {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	protocol := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
		BackendPort:         &backendPort,
		BackendProtocol:     &backendProtocol,
		FrontendPort:        &frontPort,
		FrontendProtocol:    &frontProtocol,
		LoadBalancingMethod: &method,
	}

	connections := c.Int("c")
	if c.IsSet("c") {
		protocol.MaxConn = &connections
	}

	var sessionType string
	if strings.ToLower(c.String("sticky")) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocol.SessionType = &sessionType
	} else if strings.ToLower(c.String("sticky")) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocol.SessionType = &sessionType
	} else if c.String("sticky") != "" {
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	protocols := []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
		datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
			BackendPort:         &backendPort,
			BackendProtocol:     &backendProtocol,
			FrontendPort:        &frontPort,
			FrontendProtocol:    &frontProtocol,
			LoadBalancingMethod: &method,
		},
	}
	if c.IsSet("verify") {
		orderReceipt, err := cmd.LoadBalancerManager.CreateLoadBalancerVerify(dataCenter, name, lbTypeRequest, label, protocols, subnet, c.IsSet("use-public-subnet"))
		if err != nil {
			return cli.NewExitError(T("Failed to verify load balancer with name {{.Name}} on {{.Location}}.\n",
				map[string]interface{}{"Name": name, "Location": dataCenter})+err.Error(), 2)
		}
		cmd.UI.Ok()
		table := cmd.UI.Table([]string{T("Item"), T("Cost")})
		if orderReceipt.Prices != nil {
			for _, price := range orderReceipt.Prices {
				if price.Item != nil {
					table.Add(utils.FormatStringPointer(price.Item.Description), utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee))
				}
			}
		}
		table.Print()
		return nil
	}
	orderReceipt, err := cmd.LoadBalancerManager.CreateLoadBalancer(dataCenter, name, lbTypeRequest, label, protocols, subnet, c.IsSet("use-public-subnet"))
	if err != nil {
		return cli.NewExitError(T("Failed to create load balancer with name {{.Name}} on {{.Location}}.\n",
			map[string]interface{}{"Name": name, "Location": dataCenter})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Order ID: {{.OrderID}}", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	table := cmd.UI.Table([]string{T("Item"), T("Cost")})
	if orderReceipt.OrderDetails != nil {
		for _, price := range orderReceipt.OrderDetails.Prices {
			if price.Item != nil {
				table.Add(utils.FormatStringPointer(price.Item.Description), utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee))
			}
		}
	}
	table.Print()
	return nil
}
