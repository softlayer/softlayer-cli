package virtual

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
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewMonitoringListCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *MonitoringListCommand) {
	return &MonitoringListCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *MonitoringListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	virtualId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[monitoringServiceComponent,networkMonitors[queryType,lastResult,responseAction],datacenter]"
	virtual, err := cmd.VirtualServerManager.GetInstance(virtualId, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server: {{.ID}}.\n", map[string]interface{}{"ID": virtualId})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, virtual)
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
	table.Print()
	return nil
}

func VSMonitoringListMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "monitoring-list",
		Description: T("Get details for a vsi monitors device."),
		Usage:       "${COMMAND_NAME} sl virtual monitoring-list IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
