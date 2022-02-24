package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func LoadbalOrderMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "order",
		Description: T("Order a load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal order (-n, --name NAME) (-d, --datacenter DATACENTER) (-t, --type PublicToPrivate | PrivateToPrivate | PublicToPublic ) [-l, --label LABEL] [ -s, --subnet SUBNET_ID] [--frontend-protocol PROTOCOL] [--frontend-port PORT] [--backend-protocol PROTOCOL] [--backend-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--use-public-subnet] [--verify]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name for this load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter name. It can be found from the keyName in the command '${COMMAND_NAME} sl order package-locations LBAAS' output. [required]"),
			},
			cli.StringFlag{
				Name:  "t,type",
				Usage: T("Load balancer type: PublicToPrivate | PrivateToPrivate | PublicToPublic [required]"),
			},
			cli.IntFlag{
				Name:  "s,subnet",
				Usage: T("Private subnet Id to order the load balancer. See '${COMMAND_NAME} sl loadbal order-options'. Only available in PublicToPrivate and PrivateToPrivate load balancer type"),
			},
			cli.StringFlag{
				Name:  "l,label",
				Usage: T("A descriptive label for this load balancer"),
			},
			cli.StringFlag{
				Name:  "frontend-protocol",
				Usage: T("Frontend protocol [default: HTTP]"),
			},
			cli.IntFlag{
				Name:  "frontend-port",
				Usage: T("Frontend port [default: 80]"),
			},
			cli.StringFlag{
				Name:  "backend-protocol",
				Usage: T("Backend protocol [default: HTTP]"),
			},
			cli.IntFlag{
				Name:  "backend-port",
				Usage: T("Backend port [default: 80]"),
			},
			cli.StringFlag{
				Name:  "m,method",
				Usage: T("Balancing Method: ROUNDROBIN | LEASTCONNECTION | WEIGHTED_RR, default: ROUNDROBIN"),
			},
			cli.IntFlag{
				Name:  "c, connections",
				Usage: T("Maximum number of connections to allow"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
			cli.BoolFlag{
				Name:  "use-public-subnet",
				Usage: T("If this option is specified, the public ip will be allocated from a public subnet in this account. Otherwise, it will be allocated form IBM system pool. Only available in PublicToPrivate load balancer type."),
			},
			cli.BoolFlag{
				Name:  "verify",
				Usage: T("Only verify an order, dont actually create one"),
			},
			metadata.ForceFlag(),
		},
	}
}
