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

type ProtocolAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewProtocolAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *ProtocolAddCommand) {
	return &ProtocolAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *ProtocolAddCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	frontProtocol := c.String("front-protocol")
	if frontProtocol == "" {
		frontProtocol = "HTTP"
	}

	backProtocol := c.String("back-protocol")
	if backProtocol == "" {
		backProtocol = frontProtocol
	}

	frontPort := c.Int("front-port")
	if frontPort == 0 {
		frontPort = 80
	}

	backPort := c.Int("back-port")
	if backPort == 0 {
		backPort = 80
	}

	method := c.String("m")
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
	if strings.ToLower(c.String("sticky")) == "cookie" {
		sessionType = "HTTP_COOKIE"
		protocolConfigurations.SessionType = &sessionType
	} else if strings.ToLower(c.String("sticky")) == "source-ip" {
		sessionType = "SOURCE_IP"
		protocolConfigurations.SessionType = &sessionType
	} else if c.String("sticky") != "" {
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	if c.IsSet("c") {
		connections := c.Int("c")
		protocolConfigurations.MaxConn = &connections
	}

    if c.IsSet("client-timeout") {
        cTimeout := c.Int("client-timeout")
        protocolConfigurations.ClientTimeout = &cTimeout
    }

    if c.IsSet("server-timeout") {
        sTimeout := c.Int("server-timeout")
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
