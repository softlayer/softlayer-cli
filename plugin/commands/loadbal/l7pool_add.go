package loadbal

import (
	"errors"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	bxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

type L7PoolAddCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewL7PoolAddCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *L7PoolAddCommand) {
	return &L7PoolAddCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *L7PoolAddCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return bxErr.NewMissingInputError("--id")
	}

	name := c.String("n")
	if name == "" {
		return bxErr.NewMissingInputError("-n, --name")
	}

	method := c.String("m")
	if method == "" {
		method = "ROUNDROBIN"
	}

	protocol := c.String("p")
	if protocol == "" {
		protocol = "HTTP"
	}

	healthPath := c.String("health-path")
	if healthPath == "" {
		healthPath = "/"
	}

	healthInterval := c.Int("health-interval")
	if healthInterval == 0 {
		healthInterval = 6
	}

	healthRetry := c.Int("health-retry")
	if healthRetry == 0 {
		healthRetry = 2
	}

	healthTimeout := c.Int("health-timeout")
	if healthTimeout == 0 {
		healthTimeout = 2
	}

	members := []datatypes.Network_LBaaS_L7Member{}
	var err error
	if c.IsSet("s") {
		servers := c.StringSlice("s")
		members, err = parseServer(servers)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	l7Pool := datatypes.Network_LBaaS_L7Pool{
		Name:                   &name,
		LoadBalancingAlgorithm: &method,
		Protocol:               &protocol,
	}

	l7health := datatypes.Network_LBaaS_L7HealthMonitor{
		Interval:   &healthInterval,
		Timeout:    &healthTimeout,
		MaxRetries: &healthRetry,
		UrlPath:    &healthPath,
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

	_, err = cmd.LoadBalancerManager.AddLoadBalancerL7Pool(&loadbalancerUUID, &l7Pool, members, &l7health, sessionAffinity)
	if err != nil {
		return cli.NewExitError(T("Failed to add load balancer l7 pool: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 pool added"))
	return nil
}

func parseServer(servers []string) ([]datatypes.Network_LBaaS_L7Member, error) {
	var members []datatypes.Network_LBaaS_L7Member
	for _, server := range servers {
		splitOutput := strings.Split(server, ":")
		if len(splitOutput) != 2 {
			return nil, errors.New(T("--server needs a port. {{.Server}} improperly formatted", map[string]interface{}{"Server": server}))
		}
		port, err := strconv.Atoi(splitOutput[1])
		if err != nil {
			return nil, errors.New(T("The port has to be a positive integer."))
		}
		member := datatypes.Network_LBaaS_L7Member{
			Address: &splitOutput[0],
			Port:    &port,
		}
		members = append(members, member)
	}
	return members, nil
}
