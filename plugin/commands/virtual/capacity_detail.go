package virtual

import (
	"bytes"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"strconv"
)

type CapacityDetailCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCapacityDetailCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CapacityDetailCommand) {
	return &CapacityDetailCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CapacityDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErrors.NewInvalidUsageError(T("This command requires one argument."))
	}
	id, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Reserved Capacity Group Virtual server ID")
	}
	capacity, err := cmd.VirtualServerManager.GetCapacityDetail(id)
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Reserved Capacity Gruop Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	sortby := c.String("sortby")
	if sortby == "" {
		sortby = "hostname"
	}
	var columns []string
	if c.IsSet("column") {
		columns = c.StringSlice("column")
	} else if c.IsSet("columns") {
		columns = c.StringSlice("columns")
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, capacity)
	}

	defaultColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}
	optionalColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}
	sortColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	mainTable := cmd.UI.Table([]string{T("detail")})
	mainTable.Add(utils.FormatStringPointer(capacity.Name))
	buf := new(bytes.Buffer)
	table := terminal.NewTable(buf,utils.GetColumnHeader(showColumns))
	for _, instance := range capacity.Instances {
		values := make(map[string]string)
		if instance.Guest != nil{
			values["id"] = utils.FormatIntPointer(instance.Id)
			values["hostname"] = utils.FormatStringPointer(instance.Guest.Hostname)
			values["domain"] = utils.FormatStringPointer(instance.Guest.Domain)
			values["primary_id"] = utils.FormatStringPointer(instance.Guest.PrimaryIpAddress)
			values["backend_id"] = utils.FormatStringPointer(instance.Guest.PrimaryBackendIpAddress)
		}else{
			values["id"] = utils.EMPTY_VALUE
			values["hostname"] = utils.EMPTY_VALUE
			values["domain"] = utils.EMPTY_VALUE
			values["primary_id"] = utils.EMPTY_VALUE
			values["backend_id"] = utils.EMPTY_VALUE
		}
		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()
	mainTable.Add(buf.String())
	mainTable.Print()
	return nil
}

func VSCapacityDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "capacity-detail",
		Description: T("Get Reserved Capacity Group details."),
		Usage: T(`${COMMAND_NAME} sl vs capacity-detail IDENTIFIER [OPTIONS]
EXAMPLE:
   ${COMMAND_NAME} sl vs capacity-details 12345678
    Get Reserved Capacity Group details with ID 12345678.`),
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id, hostname, domain, primary_ip, backend_ip. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			metadata.OutputFlag(),
		}}
}