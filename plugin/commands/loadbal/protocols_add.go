package loadbal

import (
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

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
	SslId               int
}

func NewProtocolAddCommand(sl *metadata.SoftlayerCommand) *ProtocolAddCommand {
	thisCmd := &ProtocolAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "protocol-add",
		Short: T("Add a new load balancer protocol"),
		Long: T(`Creates a new mapping between incoming traffic to the loadbalancer and the backend servers.
Use '{COMMAND_NAME}  sl security cert-list' to get IDs for the --ssl-id option.
See: https://cloud.ibm.com/docs/loadbalancer-service?topic=loadbalancer-service-about-ibm-cloud-load-balancer for more details

Example:
	${COMMAND_NAME} sl loadbal protocol-add --id 1115129 --front-port 443 --front-protocol HTTPS --back-port 80 --back-protocol HTTP --ssl-id 335659 --client-timeout 60 --connections 100
	Creates a new protocol on Load Balancer 1115129 that terminates SSL on port 443, mapping to a backend port 80 HTTP. Using SSL cert 335659
`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.FrontProtocol, "front-protocol", "HTTP", T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"))
	cobraCmd.Flags().StringVar(&thisCmd.BackProtocol, "back-protocol", "HTTP", T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"))
	cobraCmd.Flags().IntVar(&thisCmd.FrontPort, "front-port", 80, T("Internet side port"))
	cobraCmd.Flags().IntVar(&thisCmd.BackPort, "back-port", 80, T("Private side port"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "ROUNDROBIN", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Connections, "connections", "c", 0, T("Maximum number of connections to allow"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	cobraCmd.Flags().IntVar(&thisCmd.ClientTimeout, "client-timeout", 0, T("Client side timeout setting, in seconds"))
	cobraCmd.Flags().IntVar(&thisCmd.ServerTimeout, "server-timeout", 0, T("Server side timeout setting, in seconds"))
	cobraCmd.Flags().IntVar(&thisCmd.SslId, "ssl-id", 0, T("Identifier of the SSL certificate to attach to this protocol. Only valid for HTTPS."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ProtocolAddCommand) Run(args []string) error {
	loadbalID := cmd.Id
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return errors.New(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}))
	}

	// Sets up all the required parameters
	protocolConfigurations := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
		BackendPort:         &cmd.BackPort,
		BackendProtocol:     &cmd.BackProtocol,
		FrontendPort:        &cmd.FrontPort,
		FrontendProtocol:    &cmd.FrontProtocol,
		LoadBalancingMethod: &cmd.Method,
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
		protocolConfigurations.MaxConn = &cmd.Connections
	}

	if cmd.ClientTimeout != 0 {
		protocolConfigurations.ClientTimeout = &cmd.ClientTimeout
	}

	if cmd.ServerTimeout != 0 {
		protocolConfigurations.ServerTimeout = &cmd.ServerTimeout
	}

	if cmd.SslId != 0 {
		protocolConfigurations.TlsCertificateId = &cmd.SslId
	}
	_, err = cmd.LoadBalancerManager.AddLoadBalancerListener(
		&loadbalancerUUID, []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{protocolConfigurations},
	)
	if err != nil {
		return errors.New(T("Failed to add protocol: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Protocol added"))
	return nil
}
