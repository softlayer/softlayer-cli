package hardware

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type MonitoringListCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewMonitoringListCommand(sl *metadata.SoftlayerCommand) (cmd *MonitoringListCommand) {
	thisCmd := &MonitoringListCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "monitoring-list " + T("IDENTIFIER"),
		Short: T("Get details for a hardware monitors device."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *MonitoringListCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[monitoringServiceComponent,networkMonitors[queryType,lastResult,responseAction],datacenter]"
	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("Domain"), utils.FormatStringPointer(hardware.Domain))
	table.Add(T("Public IP"), utils.FormatStringPointer(hardware.PrimaryIpAddress))
	table.Add(T("Private IP"), utils.FormatStringPointer(hardware.PrimaryBackendIpAddress))
	if hardware.Datacenter != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(hardware.Datacenter.LongName))
	}

	if monitors := hardware.NetworkMonitors; len(monitors) > 0 {
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

	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}
	return nil
}
