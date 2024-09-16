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
	SslId               int
}

func NewProtocolEditCommand(sl *metadata.SoftlayerCommand) *ProtocolEditCommand {
	thisCmd := &ProtocolEditCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "protocol-edit",
		Short: T("Edit load balancer protocol"),
		Long: T(`Use '${COMMAND_NAME} sl loadbal detail' to find the --protocol-uuid values for a loadbalancer
Example:
	${COMMAND_NAME} sl loadbal protocol-add --id 1115129 --protocol-uuid 8ec8911a-c32d-4678-89fe-979f182c822f --ssl-id 123
	This command changes the SSL certificate
`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", -1, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.ProtocolUuid, "protocol-uuid", "", T("UUID of the protocol you want to edit."))
	cobraCmd.Flags().StringVar(&thisCmd.FrontProtocol, "front-protocol", "", T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"))
	cobraCmd.Flags().StringVar(&thisCmd.BackProtocol, "back-protocol", "", T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"))
	cobraCmd.Flags().IntVar(&thisCmd.FrontPort, "front-port", -1, T("Internet side port"))
	cobraCmd.Flags().IntVar(&thisCmd.BackPort, "back-port", -1, T("Private side port"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Connections, "connections", "c", -1, T("Maximum number of connections to allow"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	cobraCmd.Flags().IntVar(&thisCmd.ClientTimeout, "client-timeout", -1, T("Client side timeout setting, in seconds"))
	cobraCmd.Flags().IntVar(&thisCmd.ServerTimeout, "server-timeout", -1, T("Server side timeout setting, in seconds"))
	cobraCmd.Flags().IntVar(&thisCmd.SslId, "ssl-id", -1, T("Identifier of the SSL certificate to attach to this protocol. Only valid for HTTPS."))
	cobraCmd.MarkFlagRequired("id")
	cobraCmd.MarkFlagRequired("protocol-uuid")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ProtocolEditCommand) Run(args []string) error {
	protocolConfiguration := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{}

	loadbalID := cmd.Id
	if loadbalID == -1 {
		return errors.NewMissingInputError("--id")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return errors.New(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}))
	}

	if cmd.ProtocolUuid == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}
	protocolConfiguration.ListenerUuid = &cmd.ProtocolUuid

	if cmd.FrontProtocol != "" {
		protocolConfiguration.FrontendProtocol = &cmd.FrontProtocol
	}

	if cmd.BackProtocol != "" {
		protocolConfiguration.BackendProtocol = &cmd.BackProtocol
	}

	if cmd.FrontPort != -1 {
		protocolConfiguration.FrontendPort = &cmd.FrontPort
	}

	if cmd.BackPort != -1 {
		protocolConfiguration.BackendPort = &cmd.BackPort
	}

	if cmd.Method != "" {
		protocolConfiguration.LoadBalancingMethod = &cmd.Method
	}

	if cmd.ClientTimeout != -1 {
		protocolConfiguration.ClientTimeout = &cmd.ClientTimeout
	}

	if cmd.ServerTimeout != -1 {
		protocolConfiguration.ServerTimeout = &cmd.ServerTimeout
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

	if cmd.Connections != -1 {
		protocolConfiguration.MaxConn = &cmd.Connections
	}

	if cmd.SslId != 0 {
		protocolConfiguration.TlsCertificateId = &cmd.SslId
	}

	_, err = cmd.LoadBalancerManager.AddLoadBalancerListener(
		&loadbalancerUUID, []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{protocolConfiguration},
	)
	if err != nil {
		return errors.New(T("Failed to edit protocol: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Protocol edited"))
	return nil
}
