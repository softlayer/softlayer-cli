package loadbal

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewDetailCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--id")
	}

	loadbal, err := cmd.LoadBalancerManager.GetLoadBalancer(loadbalID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer with ID {{.LoadbalID}}.\n",
			map[string]interface{}{"LoadbalID": loadbalID})+err.Error(), 2)
	}
	PrintLoadbalancer(loadbal, cmd.UI)
	return nil
}

func PrintLoadbalancer(loadbal datatypes.Network_LBaaS_LoadBalancer, ui terminal.UI) {
	table := ui.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(loadbal.Id))
	table.Add(T("UUID"), utils.FormatStringPointer(loadbal.Uuid))
	table.Add(T("Name"), utils.FormatStringPointer(loadbal.Name))
	table.Add(T("Address"), utils.FormatStringPointer(loadbal.Address))

	var lbType string
	if utils.FormatIntPointer(loadbal.Type) == "1" {
		lbType = "Public to Private"
	} else if utils.FormatIntPointer(loadbal.Type) == "0" {
		lbType = "Private to Private"
	} else if utils.FormatIntPointer(loadbal.Type) == "2" {
		lbType = "Public to Public"
	}
	table.Add(T("Type"), lbType)

	if loadbal.Datacenter != nil {
		table.Add(T("Location"), utils.FormatStringPointer(loadbal.Datacenter.LongName))
	}
	table.Add(T("Description"), utils.FormatStringPointer(loadbal.Description))
	table.Add(T("Status"), fmt.Sprintf("%s/%s", utils.FormatStringPointer(loadbal.ProvisioningStatus), utils.FormatStringPointer(loadbal.OperatingStatus)))

	if len(loadbal.Listeners) > 0 {
		pools := make(map[string]string)
		bufListener := new(bytes.Buffer)
		tblListener := terminal.NewTable(bufListener, []string{
			"ID",
			"UUID",
			"Mapping",
			"Method",
			"Max Connection",
			"Modify",
			"Active",
		})
		for _, listener := range loadbal.Listeners {
			var pool *datatypes.Network_LBaaS_Pool
			if listener.DefaultPool != nil {
				pool = listener.DefaultPool
			}
			privMap := fmt.Sprintf("%s:%s", utils.FormatStringPointer(pool.Protocol), utils.FormatIntPointer(pool.ProtocolPort))
			if pool.Uuid != nil {
				pools[*pool.Uuid] = privMap
			}

			mapping := fmt.Sprintf("%s:%s -> %s", utils.FormatStringPointer(listener.Protocol), utils.FormatIntPointer(listener.ProtocolPort), privMap)
			tblListener.Add(
				utils.FormatIntPointer(listener.Id),
				utils.FormatStringPointer(listener.Uuid),
				mapping,
				utils.FormatStringPointer(pool.LoadBalancingAlgorithm),
				utils.FormatIntPointer(listener.ConnectionLimit),
				utils.FormatSLTimePointer(listener.ModifyDate),
				utils.FormatStringPointer(listener.ProvisioningStatus),
			)
		}
		tblListener.Print()
		table.Add("Protocols:", bufListener.String())
	} else {
		table.Add("Protocols:", T("Not Found"))
	}

	if len(loadbal.Members) > 0 {
		bufMember := new(bytes.Buffer)
		memCol := []string{
			"ID",
			"UUID",
			"Address",
			"Modify",
			"Active",
		}
		tblMember := terminal.NewTable(bufMember, memCol)
		for _, member := range loadbal.Members {
			row := []string{
				utils.FormatIntPointer(member.Id),
				utils.FormatStringPointer(member.Uuid),
				utils.FormatStringPointer(member.Address),
				utils.FormatSLTimePointer(member.ModifyDate),
				utils.FormatStringPointer(member.ProvisioningStatus),
			}
			tblMember.Add(row...)

		}
		tblMember.Print()
		table.Add("Members:", bufMember.String())
	} else {
		table.Add("Members:", T("Not Found"))
	}

	if len(loadbal.HealthMonitors) > 0 {
		bufHealth := new(bytes.Buffer)
		tblHealth := terminal.NewTable(bufHealth, []string{
			"ID",
			"UUID",
			"Protocol",
			"Interval",
			"Retries",
			"Timeout",
			"URL",
			"Modify",
			"Active",
		})
		for _, healthMonitor := range loadbal.HealthMonitors {
			tblHealth.Add(
				utils.FormatIntPointer(healthMonitor.Id),
				utils.FormatStringPointer(healthMonitor.Uuid),
				utils.FormatStringPointer(healthMonitor.MonitorType),
				utils.FormatIntPointer(healthMonitor.Interval),
				utils.FormatIntPointer(healthMonitor.MaxRetries),
				utils.FormatIntPointer(healthMonitor.Timeout),
				utils.FormatStringPointer(healthMonitor.UrlPath),
				utils.FormatSLTimePointer(healthMonitor.ModifyDate),
				utils.FormatStringPointer(healthMonitor.ProvisioningStatus),
			)
		}
		tblHealth.Print()
		table.Add("Health Check:", bufHealth.String())
	} else {
		table.Add("Health Check:", T("Not Found"))
	}

	if len(loadbal.L7Pools) > 0 {
		bufL7 := new(bytes.Buffer)
		tblL7 := terminal.NewTable(bufL7, []string{
			"ID",
			"UUID",
			"Name",
			"Protocol",
			"Method",
			"Modify Date",
			"ProvisioningStatus",
		})
		for _, l7Pool := range loadbal.L7Pools {
			tblL7.Add(
				utils.FormatIntPointer(l7Pool.Id),
				utils.FormatStringPointer(l7Pool.Uuid),
				utils.FormatStringPointer(l7Pool.Name),
				utils.FormatStringPointer(l7Pool.Protocol),
				utils.FormatStringPointer(l7Pool.LoadBalancingAlgorithm),
				utils.FormatSLTimePointer(l7Pool.ModifyDate),
				utils.FormatStringPointer(l7Pool.ProvisioningStatus),
			)
		}
		tblL7.Print()
		table.Add("L7 Pools:", bufL7.String())
	} else {
		table.Add("L7 Pools:", T("Not Found"))
	}

	table.Print()
}
