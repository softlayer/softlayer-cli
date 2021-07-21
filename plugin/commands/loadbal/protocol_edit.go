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
        return cli.NewExitError(T("Failed to add protocol: {{.Error}}.\n",
            map[string]interface{}{"Error": err.Error()}), 2)
    }
    cmd.UI.Ok()
    cmd.UI.Say(T("Protocol added"))
    return nil
}
