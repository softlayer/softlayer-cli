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

type CreateCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Name                string
	Datacenter          string
	Type                string
	Subnet              int
	Label               string
	FrontendProtocol    string
	FrontendPort        int
	BackendProtocol     string
	BackendPort         int
	Method              string
	Connections         int
	Sticky              string
	UsePublicSubnet     bool
	Verify              bool
	Force               bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "order",
		Short: T("Order a load balancer"),
		Long:  T("${COMMAND_NAME} sl loadbal order (-n, --name NAME) (-d, --datacenter DATACENTER) (-t, --type PublicToPrivate | PrivateToPrivate | PublicToPublic ) [-l, --label LABEL] [ -s, --subnet SUBNET_ID] [--frontend-protocol PROTOCOL] [--frontend-port PORT] [--backend-protocol PROTOCOL] [--backend-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--use-public-subnet] [--verify]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name for this load balancer [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Datacenter name. It can be found from the keyName in the command '${COMMAND_NAME} sl order package-locations LBAAS' output. [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Type, "type", "t", "", T("Load balancer type: PublicToPrivate | PrivateToPrivate | PublicToPublic [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Subnet, "subnet", "s", 0, T("Private subnet Id to order the load balancer. See '${COMMAND_NAME} sl loadbal order-options'. Only available in PublicToPrivate and PrivateToPrivate load balancer type"))
	cobraCmd.Flags().StringVarP(&thisCmd.Label, "label", "l", "", T("A descriptive label for this load balancer"))
	cobraCmd.Flags().StringVar(&thisCmd.FrontendProtocol, "frontend-protocol", "HTTP", T("Frontend protocol"))
	cobraCmd.Flags().IntVar(&thisCmd.FrontendPort, "frontend-port", 80, T("Frontend port"))
	cobraCmd.Flags().StringVar(&thisCmd.BackendProtocol, "backend-protocol", "HTTP", T("Backend protocol [default: HTTP]"))
	cobraCmd.Flags().IntVar(&thisCmd.BackendPort, "backend-port", 80, T("Backend port [default: 80]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "ROUNDROBIN", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Connections, "connections", "c", 0, T("Maximum number of connections to allow"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	cobraCmd.Flags().BoolVar(&thisCmd.UsePublicSubnet, "use-public-subnet", false, T("If this option is specified, the public ip will be allocated from a public subnet in this account. Otherwise, it will be allocated form IBM system pool. Only available in PublicToPrivate load balancer type."))
	cobraCmd.Flags().BoolVar(&thisCmd.Verify, "verify", false, T("Only verify an order, dont actually create one"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	name := cmd.Name
	if name == "" {
		return errors.NewMissingInputError("-n, --name")
	}

	dataCenter := cmd.Datacenter
	if dataCenter == "" {
		return errors.NewMissingInputError("-d, --datacenter")
	}

	lbType := cmd.Type
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

	subnet := cmd.Subnet
	if strings.ToLower(lbType) != "publictopublic" && subnet == 0 {
		return errors.NewMissingInputError("-s, --subnet")
	}
	if strings.ToLower(lbType) == "publictopublic" && subnet != 0 {
		return errors.NewInvalidUsageError(T("-s, --subnet is only available in PublicToPrivate and PrivateToPrivate load balancer type."))
	}

	if cmd.UsePublicSubnet && strings.ToLower(lbType) != "publictoprivate" {
		return errors.NewInvalidUsageError(T("--use-public-subnet is only available in PublicToPrivate."))
	}

	frontProtocol := cmd.FrontendProtocol
	if frontProtocol == "" {
		frontProtocol = "HTTP"
	}
	frontPort := cmd.FrontendPort
	if frontPort == 0 {
		frontPort = 80
	}

	backendProtocol := cmd.BackendProtocol
	if backendProtocol == "" {
		backendProtocol = "HTTP"
	}
	backendPort := cmd.BackendPort
	if backendPort == 0 {
		backendPort = 80
	}

	label := cmd.Label

	method := cmd.Method
	if method == "" {
		method = "ROUNDROBIN"
	}

	if !cmd.Force && !cmd.Verify {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
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

	connections := cmd.Connections
	if cmd.Connections != 0 {
		protocol.MaxConn = &connections
	}

	var sessionType string
	if strings.ToLower(cmd.Sticky) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocol.SessionType = &sessionType
	} else if strings.ToLower(cmd.Sticky) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocol.SessionType = &sessionType
	} else if cmd.Sticky != "" {
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
	if cmd.Verify {
		orderReceipt, err := cmd.LoadBalancerManager.CreateLoadBalancerVerify(dataCenter, name, lbTypeRequest, label, protocols, subnet, cmd.UsePublicSubnet)
		if err != nil {
			return errors.NewAPIError(T("Failed to verify load balancer with name {{.Name}} on {{.Location}}.\n",
				map[string]interface{}{"Name": name, "Location": dataCenter}), err.Error(), 2)
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
	orderReceipt, err := cmd.LoadBalancerManager.CreateLoadBalancer(dataCenter, name, lbTypeRequest, label, protocols, subnet, cmd.UsePublicSubnet)
	if err != nil {
		return errors.NewAPIError(T("Failed to create load balancer with name {{.Name}} on {{.Location}}.\n",
			map[string]interface{}{"Name": name, "Location": dataCenter}), err.Error(), 2)
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
