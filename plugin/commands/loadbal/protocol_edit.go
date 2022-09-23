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

type ProtocolEditCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Id                  int
	ProtocolUuid        string
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

func NewProtocolEditCommand(sl *metadata.SoftlayerCommand) *ProtocolEditCommand {
	thisCmd := &ProtocolEditCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "protocol-edit",
		Short: T("Edit load balancer protocol"),
		Long:  T("${COMMAND_NAME} sl loadbal protocol-edit (--id LOADBAL_ID) (--protocol-uuid PROTOCOL_UUID) [--front-protocol PROTOCOL] [back-protocol PROTOCOL] [--front-port PORT] [--back-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--client-timeout SECONDS] [--server-timeout SECONDS]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.ProtocolUuid, "protocol-uuid", "", T("UUID of the protocol you want to edit."))
	cobraCmd.Flags().StringVar(&thisCmd.FrontProtocol, "front-protocol", "HTTP", T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"))
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

func (cmd *ProtocolEditCommand) Run(args []string) error {
	protocolConfiguration := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{}

	loadbalID := cmd.Id
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	protoUUID := cmd.ProtocolUuid
	if protoUUID == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}
	protocolConfiguration.ListenerUuid = &protoUUID

	if cmd.FrontProtocol != "" {
		frontProtocol := cmd.FrontProtocol
		protocolConfiguration.FrontendProtocol = &frontProtocol
	}

	if cmd.BackProtocol != "" {
		backProtocol := cmd.BackProtocol
		protocolConfiguration.BackendProtocol = &backProtocol
	}

	if cmd.FrontPort != 0 {
		frontPort := cmd.FrontPort
		protocolConfiguration.FrontendPort = &frontPort
	}

	if cmd.BackPort != 0 {
		backPort := cmd.BackPort
		protocolConfiguration.BackendPort = &backPort
	}

	if cmd.Method != "" {
		method := cmd.Method
		protocolConfiguration.LoadBalancingMethod = &method
	}

	if cmd.ClientTimeout != 0 {
		cTimeout := cmd.ClientTimeout
		protocolConfiguration.ClientTimeout = &cTimeout
	}

	if cmd.ServerTimeout != 0 {
		sTimeout := cmd.ServerTimeout
		protocolConfiguration.ServerTimeout = &sTimeout
	}

	var sessionType string
	if strings.ToLower(cmd.Sticky) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocolConfiguration.SessionType = &sessionType
	} else if strings.ToLower(cmd.Sticky) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocolConfiguration.SessionType = &sessionType
	} else if cmd.Sticky != "" {
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	if cmd.Connections != 0 {
		connections := cmd.Connections
		protocolConfiguration.MaxConn = &connections
	}

	_, err = cmd.LoadBalancerManager.AddLoadBalancerListener(&loadbalancerUUID, []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{protocolConfiguration})
	if err != nil {
		return cli.NewExitError(T("Failed to edit protocol: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Protocol edited"))
	return nil
}
