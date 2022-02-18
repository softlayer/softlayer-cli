package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ProtocolEditCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewProtocolEditCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *ProtocolEditCommand) {
	return &ProtocolEditCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *ProtocolEditCommand) Run(c *cli.Context) error {
	protocolConfiguration := datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{}

	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	protoUUID := c.String("protocol-uuid")
	if protoUUID == "" {
		return errors.NewMissingInputError("--protocol-uuid")
	}
	protocolConfiguration.ListenerUuid = &protoUUID

	if c.IsSet("front-protocol") {
		frontProtocol := c.String("front-protocol")
		protocolConfiguration.FrontendProtocol = &frontProtocol
	}

	if c.IsSet("back-protocol") {
		backProtocol := c.String("back-protocol")
		protocolConfiguration.BackendProtocol = &backProtocol
	}

	if c.IsSet("front-port") {
		frontPort := c.Int("front-port")
		protocolConfiguration.FrontendPort = &frontPort
	}

	if c.IsSet("back-port") {
		backPort := c.Int("back-port")
		protocolConfiguration.BackendPort = &backPort
	}

	if c.IsSet("m") {
		method := c.String("m")
		protocolConfiguration.LoadBalancingMethod = &method
	}

	if c.IsSet("client-timeout") {
		cTimeout := c.Int("client-timeout")
		protocolConfiguration.ClientTimeout = &cTimeout
	}

	if c.IsSet("server-timeout") {
		sTimeout := c.Int("server-timeout")
		protocolConfiguration.ServerTimeout = &sTimeout
	}

	var sessionType string
	if strings.ToLower(c.String("sticky")) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocolConfiguration.SessionType = &sessionType
	} else if strings.ToLower(c.String("sticky")) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocolConfiguration.SessionType = &sessionType
	} else if c.String("sticky") != "" {
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	if c.IsSet("c") {
		connections := c.Int("c")
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

func LoadbalProtocolEditMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "protocol-edit",
		Description: T("Edit load balancer protocol"),
		Usage:       "${COMMAND_NAME} sl loadbal protocol-edit (--id LOADBAL_ID) (--protocol-uuid PROTOCOL_UUID) [--front-protocol PROTOCOL] [back-protocol PROTOCOL] [--front-port PORT] [--back-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--client-timeout SECONDS] [--server-timeout SECONDS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "protocol-uuid",
				Usage: T("UUID of the protocol you want to edit."),
			},
			cli.StringFlag{
				Name:  "front-protocol",
				Usage: T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"),
			},
			cli.StringFlag{
				Name:  "back-protocol",
				Usage: T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"),
			},
			cli.IntFlag{
				Name:  "front-port",
				Usage: T("Internet side port. Default: 80"),
			},
			cli.IntFlag{
				Name:  "back-port",
				Usage: T("Private side port. Default: 80"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]. Default: ROUNDROBIN"),
			},
			cli.IntFlag{
				Name:  "c, connections",
				Usage: T("Maximum number of connections to allow"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
			cli.IntFlag{
				Name:  "client-timeout",
				Usage: T("Client side timeout setting, in seconds"),
			},
			cli.IntFlag{
				Name:  "server-timeout",
				Usage: T("Server side timeout setting, in seconds"),
			},
		},
	}
}
