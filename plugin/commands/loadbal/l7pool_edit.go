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

type L7PoolEditCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PoolUuid            string
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

func NewL7PoolEditCommand(sl *metadata.SoftlayerCommand) *L7PoolEditCommand {
	thisCmd := &L7PoolEditCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7pool-edit",
		Short: T("Edit a L7 pool"),
		Long:  T("${COMMAND_NAME} sl loadbal l7pool-edit (--pool-uuid L7POOL_UUID) [-m, --method METHOD] [-s, --server BACKEND_IP:PORT] [-p, --protocol PROTOCOL] [--health-path PATH] [--health-interval INTERVAL] [--health-retry RETRY] [--health-timeout TIMEOUT] [--sticky cookie | source-ip]"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.PoolUuid, "pool-uuid", "", T("UUID for the load balancer pool [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name of the load balancer L7 pool"))
	cobraCmd.Flags().StringVarP(&thisCmd.Method, "method", "m", "", T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Protocol, "protocol", "p", "", T("Protocol type to use for incoming connections"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Server, "server", "s", []string{}, T("Backend servers that are part of this pool. Format: BACKEND_IP:PORT. eg. 10.0.0.1:80 (multiple occurrence permitted)"))
	cobraCmd.Flags().StringVar(&thisCmd.HealthPath, "health-path", "", T("Health check path"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthInterval, "health-interval", 0, T("Health check interval between checks"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthRetry, "health-retry", 0, T("Health check number of times before marking as DOWN"))
	cobraCmd.Flags().IntVar(&thisCmd.HealthTimeout, "health-timeout", 0, T("Health check timeout"))
	cobraCmd.Flags().StringVar(&thisCmd.Sticky, "sticky", "", T("Use 'cookie' or 'source-ip' to stick"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PoolEditCommand) Run(args []string) error {
	poolUUID := cmd.PoolUuid
	if poolUUID == "" {
		return errors.NewMissingInputError("--pool-uuid")
	}
	if cmd.Name == "" && cmd.Method == "" && cmd.Protocol == "" && len(cmd.Server) == 0 && cmd.HealthPath == "" && cmd.HealthInterval == 0 && cmd.HealthRetry == 0 && cmd.HealthTimeout == 0 && cmd.Sticky == "" {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	l7Pool := datatypes.Network_LBaaS_L7Pool{}
	if cmd.Name != "" {
		name := cmd.Name
		l7Pool.Name = &name
	}
	if cmd.Method != "" {
		method := cmd.Method
		l7Pool.LoadBalancingAlgorithm = &method
	}
	if cmd.Protocol != "" {
		protocol := cmd.Protocol
		l7Pool.Protocol = &protocol
	}

	var members []datatypes.Network_LBaaS_L7Member
	var err error
	if len(cmd.Server) != 0 {
		servers := cmd.Server
		members, err = parseServer(servers)
		if err != nil {
			return err
		}
		l7Pool.L7Members = members
	}

	l7health := datatypes.Network_LBaaS_L7HealthMonitor{}
	if cmd.HealthPath != "" {
		healthPath := cmd.HealthPath
		l7health.UrlPath = &healthPath
	}
	if cmd.HealthInterval != 0 {
		healthInterval := cmd.HealthInterval
		l7health.Interval = &healthInterval
	}
	if cmd.HealthRetry != 0 {
		healthRetry := cmd.HealthRetry
		l7health.MaxRetries = &healthRetry
	}
	if cmd.HealthTimeout != 0 {
		healthTimeout := cmd.HealthTimeout
		l7health.Timeout = &healthTimeout
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
		return errors.NewInvalidUsageError(T("Value of option '--sticky' should be cookie or source-ip"))
	}

	_, err = cmd.LoadBalancerManager.UpdateLoadBalancerL7Pool(&poolUUID, &l7Pool, &l7health, sessionAffinity)
	if err != nil {
		return errors.New(T("Failed to update l7 pool: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("L7 pool updated"))
	return nil
}
