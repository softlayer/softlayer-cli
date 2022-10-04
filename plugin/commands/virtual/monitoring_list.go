package virtual

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type MonitoringListCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewMonitoringListCommand(sl *metadata.SoftlayerCommand) (cmd *MonitoringListCommand) {
	thisCmd := &MonitoringListCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "monitoring-list " + T("IDENTIFIER"),
		Short: T("Get details for a vsi monitors device."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *MonitoringListCommand) Run(args []string) error {

	virtualId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[monitoringServiceComponent,networkMonitors[queryType,lastResult,responseAction],datacenter]"
	virtual, err := cmd.VirtualServerManager.GetInstance(virtualId, mask)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get virtual server: {{.ID}}.\n", map[string]interface{}{"ID": virtualId}), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("Domain"), utils.FormatStringPointer(virtual.Domain))
	table.Add(T("Public IP"), utils.FormatStringPointer(virtual.PrimaryIpAddress))
	table.Add(T("Private IP"), utils.FormatStringPointer(virtual.PrimaryBackendIpAddress))
	if virtual.Datacenter != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(virtual.Datacenter.LongName))
	}

	if monitors := virtual.NetworkMonitors; len(monitors) > 0 {
		buf := new(bytes.Buffer)
		monitorTable := terminal.NewTable(buf, []string{T("Id"), T("IpAddress"), T("Status"), T("Type"), T("Notify")})
		for _, monitor := range monitors {
			monitorTable.Add(
				utils.FormatIntPointer(monitor.Id),
				utils.FormatStringPointer(monitor.IpAddress),
				utils.FormatStringPointer(monitor.Status),
				utils.FormatStringPointer(monitor.QueryType.Name),
				utils.FormatStringPointer(monitor.ResponseAction.ActionDescription),
			)
		}
		monitorTable.Print()
		table.Add("Monitors", buf.String())
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
