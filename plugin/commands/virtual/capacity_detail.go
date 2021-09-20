package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
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
	id, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Reserved Capacity Group Virtual server ID")
	}
	capacity, err := cmd.VirtualServerManager.GetCapacityDetail(id)

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

	defaultColumns := []string{"name","id", "hostname", "domain", "primary_id", "backend_id"}
	optionalColumns := []string{"name","id", "hostname", "domain", "primary_id", "backend_id"}
	sortColumns := []string{"name","id", "hostname", "domain", "primary_id", "backend_id"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	values := make(map[string]string)
	values["name"] = utils.FormatStringPointer(capacity.Name)
	values["id"] = utils.FormatIntPointer(capacity.Id)
	values["hostname"] = utils.FormatStringPointer(capacity.Instances[0].Guest.Hostname)
	values["domain"] = utils.FormatStringPointer(capacity.Instances[0].Guest.Domain)
	values["primary_id"] = utils.FormatStringPointer(capacity.Instances[0].Guest.PrimaryIpAddress)
	values["backend_id"] = utils.FormatStringPointer(capacity.Instances[0].Guest.PrimaryBackendIpAddress)

	row := make([]string, len(showColumns))
	for i, col := range showColumns {
		row[i] = values[col]
	}
	table.Add(row...)
	table.Print()
	return nil
}
