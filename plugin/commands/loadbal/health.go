package loadbal

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type HealthChecksCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewHealthChecksCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *HealthChecksCommand) {
	return &HealthChecksCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *HealthChecksCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("lb-id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--lb-id")
	}

	healthUUID := c.String("health-uuid")
	if healthUUID == "" {
		return errors.NewMissingInputError("--health-uuid")
	}

	if !c.IsSet("i") && !c.IsSet("r") && !c.IsSet("t") && !c.IsSet("u") {
		return errors.NewInvalidUsageError(fmt.Sprintf("%s :%s, %s, %s, %s,", T("At least one of these flags is required"), "-i, --interval", "-r, --retry", "-t, --timeout", " -u, --url"))
	}

	loadbalancer, err := cmd.LoadBalancerManager.GetLoadBalancer(loadbalID, "uuid,healthMonitors,listeners[uuid,defaultPool[healthMonitor]]")
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	var healthCheck datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration
	var find bool
	for _, listener := range loadbalancer.Listeners {
		if listener.DefaultPool != nil && listener.DefaultPool.HealthMonitor != nil && listener.DefaultPool.HealthMonitor.Uuid != nil && *listener.DefaultPool.HealthMonitor.Uuid == healthUUID {
			find = true
			var backendPort, interval, maxRetries, timeout *int
			var backendProtocol, urlPath *string

			backendProtocol = listener.DefaultPool.Protocol
			backendPort = listener.DefaultPool.ProtocolPort
			interval = listener.DefaultPool.HealthMonitor.Interval
			maxRetries = listener.DefaultPool.HealthMonitor.MaxRetries
			timeout = listener.DefaultPool.HealthMonitor.Timeout
			urlPath = listener.DefaultPool.HealthMonitor.UrlPath

			healthCheck = datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration{
				BackendProtocol:   backendProtocol,
				BackendPort:       backendPort,
				HealthMonitorUuid: &healthUUID,
				Interval:          interval,
				MaxRetries:        maxRetries,
				Timeout:           timeout,
				UrlPath:           urlPath,
			}
		}
	}

	if find == false {
		return cli.NewExitError(T("Unable to find health check with UUID of '{{.UUID}}' in load balancer {{.ID}}.", map[string]interface{}{"UUID": healthUUID, "ID": loadbalID}), 2)
	}

	if c.IsSet("u") && healthCheck.BackendProtocol != nil && *healthCheck.BackendProtocol == "TCP" {
		return cli.NewExitError(T("--url cannot be used with TCP checks."), 2)
	}

	interval := c.Int("i")
	if c.IsSet("i") {
		healthCheck.Interval = &interval
	}

	retry := c.Int("r")
	if c.IsSet("r") {
		healthCheck.MaxRetries = &retry
	}

	timeout := c.Int("t")
	if c.IsSet("t") {
		healthCheck.Timeout = &timeout
	}

	url := c.String("u")
	if c.IsSet("u") {
		healthCheck.UrlPath = &url
	}

	updatedLoadbalancer, err := cmd.LoadBalancerManager.UpdateLBHealthMonitors(loadbalancer.Uuid, []datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration{healthCheck})
	if err != nil {
		return cli.NewExitError(T("Failed to update health check: ")+err.Error(), 2)
	}
	PrintLoadbalancer(updatedLoadbalancer, cmd.UI)
	return nil
}

func LoadbalHealthMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "health-edit",
		Description: T("Edit load balancer health check"),
		Usage:       "${COMMAND_NAME} sl loadbal health-edit (--lb-id LOADBAL_ID)  (--health-uuid HEALTH_CHECK_UUID) [-i, --interval INTERVAL] [-r, --retry RETRY] [-t, --timeout TIMEOUT] [-u, --url URL]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "lb-id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "health-uuid",
				Usage: T("Health check UUID to modify [required]"),
			},
			cli.IntFlag{
				Name:  "i,interval",
				Usage: T("Seconds between checks. [2-60]"),
			},
			cli.IntFlag{
				Name:  "r,retry",
				Usage: T("Number of times before marking as DOWN. [1-10]"),
			},
			cli.IntFlag{
				Name:  "t,timeout",
				Usage: T("Seconds to wait for a connection. [1-59]"),
			},
			cli.StringFlag{
				Name:  "u,url",
				Usage: T("Url path for HTTP/HTTPS checks"),
			},
		},
	}
}
