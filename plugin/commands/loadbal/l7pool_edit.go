package loadbal

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type L7PoolEditCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PoolEditCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PoolEditCommand) {
	return &L7PoolEditCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PoolEditCommand) Run(c *cli.Context) error {
	poolUUID := c.String("pool-uuid")
	if poolUUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}

	if c.NumFlags() <= 1 {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	l7Pool := datatypes.Network_LBaaS_L7Pool{}
	if c.IsSet("n") {
		name := c.String("n")
		l7Pool.Name = &name
	}
	if c.IsSet("m") {
		method := c.String("m")
		l7Pool.LoadBalancingAlgorithm = &method
	}
	if c.IsSet("p") {
		protocol := c.String("p")
		l7Pool.Protocol = &protocol
	}

	var members []datatypes.Network_LBaaS_L7Member
	var err error
	if c.IsSet("s") {
		servers := c.StringSlice("s")
		members, err = parseServer(servers)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		l7Pool.L7Members = members
	}

	l7health := datatypes.Network_LBaaS_L7HealthMonitor{}
	if c.IsSet("health-path") {
		healthPath := c.String("health-path")
		l7health.UrlPath = &healthPath
	}
	if c.IsSet("health-interval") {
		healthInterval := c.Int("health-interval")
		l7health.Interval = &healthInterval
	}
	if c.IsSet("health-retry") {
		healthRetry := c.Int("health-retry")
		l7health.MaxRetries = &healthRetry
	}
	if c.IsSet("health-timeout") {
		healthTimeout := c.Int("health-timeout")
		l7health.Timeout = &healthTimeout
	}

	var sessionAffinity *datatypes.Network_LBaaS_L7SessionAffinity
	if strings.ToLower(c.String("sticky")) == "cookie" {
		sessionAffinityType := "HTTP_COOKIE"
		sessionAffinity = &datatypes.Network_LBaaS_L7SessionAffinity{
			Type: &sessionAffinityType,
		}
	} else if strings.ToLower(c.String("sticky")) == "source-ip" {
		sessionAffinityType := "SOURCE_IP"
		sessionAffinity = &datatypes.Network_LBaaS_L7SessionAffinity{
			Type: &sessionAffinityType,
		}
	} else if c.String("sticky") != "" {
		return bxErr.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	_, err = cmd.LoadBalancerManager.UpdateLoadBalancerL7Pool(&poolUUID, &l7Pool, &l7health, sessionAffinity)
	if err != nil {
		return cli.NewExitError(T("Failed to update l7 pool: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 pool updated"))
	return nil
}

func LoadbalL7PoolEditMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "l7pool-edit",
		Description: T("Edit a L7 pool"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-edit (--pool-uuid L7POOL_UUID) [-m, --method METHOD] [-s, --server BACKEND_IP:PORT] [-p, --protocol PROTOCOL] [--health-path PATH] [--health-interval INTERVAL] [--health-retry RETRY] [--health-timeout TIMEOUT] [--sticky cookie | source-ip]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "pool-uuid",
				Usage: T("UUID for the load balancer pool [required]"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"),
			},
			cli.StringFlag{
				Name:  "p, protocol",
				Usage: T("Protocol type to use for incoming connections"),
			},
			cli.StringSliceFlag{
				Name:  "s, server",
				Usage: T("Backend servers that are part of this pool. Format: BACKEND_IP:PORT. eg. 10.0.0.1:80 (multiple occurrence permitted)"),
			},
			cli.StringFlag{
				Name:  "health-path",
				Usage: T("Health check path"),
			},
			cli.IntFlag{
				Name:  "health-interval",
				Usage: T("Health check interval between checks"),
			},
			cli.IntFlag{
				Name:  "health-retry",
				Usage: T("Health check number of times before marking as DOWN"),
			},
			cli.IntFlag{
				Name:  "health-timeout",
				Usage: T("Health check timeout"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
		},
	}
}
