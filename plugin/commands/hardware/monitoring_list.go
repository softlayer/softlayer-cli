package hardware

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type MonitoringListCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewMonitoringListCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *MonitoringListCommand) {
	return &MonitoringListCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *MonitoringListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[monitoringServiceComponent,networkMonitors[queryType,lastResult,responseAction],datacenter]"
	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error(), 2)
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

func HardwareMonitoringListMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "monitoring-list",
		Description: T("Get details for a hardware monitors device."),
		Usage:       "${COMMAND_NAME} sl hardware monitoring-list IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
