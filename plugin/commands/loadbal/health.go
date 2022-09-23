package loadbal

import (
	"fmt"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type HealthChecksCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	HealthUuid          string
	Interval            int
	Retry               int
	Timeout             int
	Url                 string
}

func NewHealthChecksCommand(sl *metadata.SoftlayerCommand) *HealthChecksCommand {
	thisCmd := &HealthChecksCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "health-edit " + T("IDENTIFIER"),
		Short: T("Edit load balancer health check"),
		Long:  T("${COMMAND_NAME} sl loadbal health-edit (--lb-id LOADBAL_ID)  (--health-uuid HEALTH_CHECK_UUID) [-i, --interval INTERVAL] [-r, --retry RETRY] [-t, --timeout TIMEOUT] [-u, --url URL]"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.HealthUuid, "health-uuid", "", T("Health check UUID to modify [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Interval, "interval", "i", 0, T("Seconds between checks. [2-60]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Retry, "retry", "r", 0, T("Number of times before marking as DOWN. [1-10]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Timeout, "timeout", "t", 0, T("Seconds to wait for a connection. [1-59]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Url, "url", "u", "", T("Url path for HTTP/HTTPS checks"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *HealthChecksCommand) Run(args []string) error {
	loadbalID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("LoadBalancer ID")
	}
	healthUUID := cmd.HealthUuid
	if healthUUID == "" {
		return errors.NewMissingInputError("--health-uuid")
	}

	if cmd.Interval == 0 && cmd.Retry == 0 && cmd.Timeout == 0 && cmd.Url == "" {
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

	if cmd.Url != "" && healthCheck.BackendProtocol != nil && *healthCheck.BackendProtocol == "TCP" {
		return cli.NewExitError(T("--url cannot be used with TCP checks."), 2)
	}

	interval := cmd.Interval
	if cmd.Interval != 0 {
		healthCheck.Interval = &interval
	}

	retry := cmd.Retry
	if cmd.Retry != 0 {
		healthCheck.MaxRetries = &retry
	}

	timeout := cmd.Timeout
	if cmd.Timeout != 0 {
		healthCheck.Timeout = &timeout
	}

	url := cmd.Url
	if cmd.Url != "" {
		healthCheck.UrlPath = &url
	}

	updatedLoadbalancer, err := cmd.LoadBalancerManager.UpdateLBHealthMonitors(loadbalancer.Uuid, []datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration{healthCheck})
	if err != nil {
		return errors.NewAPIError(T("Failed to update health check: "), err.Error(), 2)
	}
	PrintLoadbalancer(updatedLoadbalancer, cmd.UI)
	return nil
}
