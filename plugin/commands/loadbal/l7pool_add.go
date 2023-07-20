package loadbal

import (
	"errors"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type L7PoolAddCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	Id                  int
	Name                string
	Method              string
	Protocol            string
	Server              []string
	HealthPath          string
	HealthInterval      int
	HealthRetry         int
	HealthTimeout       int
	Sticky              string
}

func NewL7PoolAddCommand(sl *metadata.SoftlayerCommand) *L7PoolAddCommand {
	thisCmd := &L7PoolAddCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7pool-add",
		Short: T("Add a new L7 pool"),
		Long:  T("${COMMAND_NAME} sl loadbal l7pool-add (--id LOADBAL_ID) (-n, --name NAME) [-m, --method METHOD] [-s, --server BACKEND_IP:PORT] [-p, --protocol PROTOCOL] [--health-path PATH] [--health-interval INTERVAL] [--health-retry RETRY] [--health-timeout TIMEOUT] [--sticky cookie | source-ip]\n\n Adds a new l7 pool \n\n -s is in colon deliminated format to make grouping IP:port:weight a bit easier."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the load balancer [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name for this L7 pool. [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "ROUNDROBIN", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Protocol, "protocol", "p", "HTTP", T("Protocol type to use for incoming connections"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Server, "server", "s", []string{}, T("Backend servers that are part of this pool. Format: BACKEND_IP:PORT. eg. 10.0.0.1:80 (multiple occurrence permitted)"))
	cobraCmd.Flags().StringVar(&thisCmd.HealthPath, "health-path", "/", T("Health check path"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthInterval, "health-interval", 5, T("Health check interval between checks"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthRetry, "health-retry", 2, T("Health check number of times before marking as DOWN"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthTimeout, "health-timeout", 2, T("Health check timeout"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PoolAddCommand) Run(args []string) error {
	loadbalID := cmd.Id
	if loadbalID == 0 {
		return bxErr.NewMissingInputError("--id")
	}

	name := cmd.Name
	if name == "" {
		return bxErr.NewMissingInputError("-n, --name")
	}

	method := cmd.Method
	if method == "" {
		method = "ROUNDROBIN"
	}

	protocol := cmd.Protocol
	if protocol == "" {
		protocol = "HTTP"
	}

	healthPath := cmd.HealthPath
	if healthPath == "" {
		healthPath = "/"
	}

	healthInterval := cmd.HealthInterval
	if healthInterval == 0 {
		healthInterval = 6
	}

	healthRetry := cmd.HealthRetry
	if healthRetry == 0 {
		healthRetry = 2
	}

	healthTimeout := cmd.HealthTimeout
	if healthTimeout == 0 {
		healthTimeout = 2
	}

	members := []datatypes.Network_LBaaS_L7Member{}
	var err error
	if len(cmd.Server) != 0 {
		servers := cmd.Server
		members, err = parseServer(servers)
		if err != nil {
			return err
		}
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return errors.New(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}))
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
	if strings.ToLower(cmd.Sticky) == "cookie" {
		sessionAffinityType := "HTTP_COOKIE"
		sessionAffinity = &datatypes.Network_LBaaS_L7SessionAffinity{
			Type: &sessionAffinityType,
		}
	} else if strings.ToLower(cmd.Sticky) == "source-ip" {
		sessionAffinityType := "SOURCE_IP"
		sessionAffinity = &datatypes.Network_LBaaS_L7SessionAffinity{
			Type: &sessionAffinityType,
		}
	} else if cmd.Sticky != "" {
		return bxErr.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	_, err = cmd.LoadBalancerManager.AddLoadBalancerL7Pool(&loadbalancerUUID, &l7Pool, members, &l7health, sessionAffinity)
	if err != nil {
		return errors.New(T("Failed to add load balancer l7 pool: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
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
