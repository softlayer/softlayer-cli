package loadbal

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ProtocolAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Id                  int
	FrontProtocol       string
	BackProtocol        string
	FrontPort           int
	BackPort            int
	Method              string
	Connections         int
	Sticky              string
	ClientTimeout       int
	ServerTimeout       int
}

func NewProtocolAddCommand(sl *metadata.SoftlayerCommand) *ProtocolAddCommand {
	thisCmd := &ProtocolAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "protocol-add",
		Short: T("Add a new load balancer protocol."),
		Long:  T("${COMMAND_NAME} sl loadbal protocol-add (--id LOADBAL_ID) [--front-protocol PROTOCOL] [back-protocol PROTOCOL] [--front-port PORT] [--back-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--client-timeout SECONDS] [--server-timeout SECONDS]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.FrontProtocol, "front-protocol", "HTTP", T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]"))
	cobraCmd.Flags().StringVar(&thisCmd.BackProtocol, "back-protocol", "", T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"))
	cobraCmd.Flags().IntVar(&thisCmd.FrontPort, "front-port", 80, T("Internet side port"))
	cobraCmd.Flags().IntVar(&thisCmd.BackPort, "back-port", 80, T("Private side port"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "ROUNDROBIN", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Connections, "connections", "c", 0, T("Maximum number of connections to allow"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	cobraCmd.Flags().IntVar(&thisCmd.ClientTimeout, "client-timeout", 0, T("Client side timeout setting, in seconds"))
	cobraCmd.Flags().IntVar(&thisCmd.ServerTimeout, "server-timeout", 0, T("Server side timeout setting, in seconds"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ProtocolAddCommand) Run(args []string) error {
	loadbalID := cmd.Id
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	frontProtocol := cmd.FrontProtocol
	if frontProtocol == "" {
		frontProtocol = "HTTP"
	}

	backProtocol := cmd.BackProtocol
	if backProtocol == "" {
		backProtocol = frontProtocol
	}

	frontPort := cmd.FrontPort
	if frontPort == 0 {
		frontPort = 80
	}

	backPort := cmd.BackPort
	if backPort == 0 {
		backPort = 80
	}

	method := cmd.Method
	if method == "" {
		method = "ROUNDROBIN"
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	protocolConfigurations := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
		BackendPort:         &backPort,
		BackendProtocol:     &backProtocol,
		FrontendPort:        &frontPort,
		FrontendProtocol:    &frontProtocol,
		LoadBalancingMethod: &method,
	}

	var sessionType string
	if strings.ToLower(cmd.Sticky) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocolConfigurations.SessionType = &sessionType
	} else if strings.ToLower(cmd.Sticky) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocolConfigurations.SessionType = &sessionType
	} else if cmd.Sticky != "" {
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	if cmd.Connections != 0 {
		connections := cmd.Connections
		protocolConfigurations.MaxConn = &connections
	}

	if cmd.ClientTimeout != 0 {
		cTimeout := cmd.ClientTimeout
		protocolConfigurations.ClientTimeout = &cTimeout
	}

	if cmd.ServerTimeout != 0 {
		sTimeout := cmd.ServerTimeout
		protocolConfigurations.ServerTimeout = &sTimeout
	}

	_, err = cmd.LoadBalancerManager.AddLoadBalancerListener(&loadbalancerUUID, []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{protocolConfigurations})
	if err != nil {
		return cli.NewExitError(T("Failed to add protocol: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Protocol added"))
	return nil
}
