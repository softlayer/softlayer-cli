package loadbal

import (
	"bytes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type L7PoolDetailCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
	PoolId              int
}

func NewL7PoolDetailCommand(sl *metadata.SoftlayerCommand) *L7PoolDetailCommand {
	thisCmd := &L7PoolDetailCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "l7pool-detail",
		Short: T("Show L7 pool details."),
		Long:  T("${COMMAND_NAME} sl loadbal l7pool-detail (--pool-id L7POOL_ID)"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.PoolId, "pool-id", 0, T("ID for the load balancer pool [required]"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *L7PoolDetailCommand) Run(args []string) error {
	l7PoolID := cmd.PoolId
	if l7PoolID == 0 {
		return errors.NewMissingInputError("--pool-id")
	}

	l7pool, err := cmd.LoadBalancerManager.GetLoadBalancerL7Pool(l7PoolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get L7 Pool {{.L7PoolID}}: {{.Error}}.\n",
			map[string]interface{}{"L7PoolID": l7PoolID, "Error": err.Error()}), 2)
	}

	l7SessionAffinity, err := cmd.LoadBalancerManager.GetL7SessionAffinity(l7PoolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get L7 Pool Session Affinity: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	l7HealthMonitor, err := cmd.LoadBalancerManager.GetL7HealthMonitor(l7PoolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get L7 Health Monitor: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	l7Members, err := cmd.LoadBalancerManager.ListL7Members(l7PoolID)
	if err != nil {
		return cli.NewExitError(T("Failed to get L7 Members: {{.Error}}.\n",
			map[string]interface{}{"Error": err.Error()}), 2)
	}

	printL7Pool(l7pool, l7Members, l7HealthMonitor, l7SessionAffinity, cmd.UI)
	return nil
}
func printL7Pool(l7pool datatypes.Network_LBaaS_L7Pool, l7Members []datatypes.Network_LBaaS_L7Member, l7HealthMonitor datatypes.Network_LBaaS_L7HealthMonitor, l7SessionAffinity datatypes.Network_LBaaS_L7SessionAffinity, ui terminal.UI) {
	ui.Ok()
	table := ui.Table([]string{T("Name"), T("Value")})
	table.Add(T("Name"), utils.FormatStringPointer(l7pool.Name))
	table.Add(T("ID"), utils.FormatIntPointer(l7pool.Id))
	table.Add(T("UUID"), utils.FormatStringPointer(l7pool.Uuid))
	table.Add(T("Method"), utils.FormatStringPointer(l7pool.LoadBalancingAlgorithm))
	table.Add(T("Protocol"), utils.FormatStringPointer(l7pool.Protocol))
	if l7SessionAffinity.Type != nil {
		table.Add(T("Session Stickiness"), utils.FormatStringPointer(l7SessionAffinity.Type))
	}

	bufHealth := new(bytes.Buffer)
	tblHealth := terminal.NewTable(bufHealth, []string{
		"Interval",
		"Retries",
		"Type",
		"Timeout",
		"URL",
		"Modify",
		"Active",
	})
	tblHealth.Add(
		utils.FormatIntPointer(l7HealthMonitor.Interval),
		utils.FormatIntPointer(l7HealthMonitor.MaxRetries),
		utils.FormatStringPointer(l7HealthMonitor.MonitorType),
		utils.FormatIntPointer(l7HealthMonitor.Timeout),
		utils.FormatStringPointer(l7HealthMonitor.UrlPath),
		utils.FormatSLTimePointer(l7HealthMonitor.ModifyDate),
		utils.FormatStringPointer(l7HealthMonitor.ProvisioningStatus),
	)
	tblHealth.Print()
	table.Add("Health Check:", bufHealth.String())

	bufMember := new(bytes.Buffer)
	memCol := []string{
		"ID",
		"UUID",
		"Address",
		"Weight",
		"Modify",
		"Active",
	}
	tblMember := terminal.NewTable(bufMember, memCol)
	for _, member := range l7Members {
		row := []string{
			utils.FormatIntPointer(member.Id),
			utils.FormatStringPointer(member.Uuid),
			utils.FormatStringPointer(member.Address),
			utils.FormatIntPointer(member.Weight),
			utils.FormatSLTimePointer(member.ModifyDate),
			utils.FormatStringPointer(member.ProvisioningStatus),
		}
		tblMember.Add(row...)

	}
	tblMember.Print()
	table.Add("Members:", bufMember.String())

	table.Print()
}
