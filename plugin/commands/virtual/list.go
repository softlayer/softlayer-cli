package virtual

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewListCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

var maskMap = map[string]string{
	"id":          "id",
	"hostname":    "hostname",
	"domain":      "domain",
	"cpu":         "maxCpu",
	"memory":      "maxMemory",
	"public_ip":   "primaryIpAddress",
	"private_ip":  "primaryBackendIpAddress",
	"datacenter":  "datacenter",
	"action":      "activeTransaction.transactionStatus.name",
	"guid":        "globalIdentifier",
	"power_state": "powerState.name",
	"created_by":  "billingItem.orderItem.order.userRecord.username",
	"tags":        "tagReferences",
}

func (cmd *ListCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

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

	defaultColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter", "action"}
	optionalColumns := []string{"guid", "power_state", "created_by", "tags"}
	sortColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	if c.IsSet("hourly") && c.IsSet("monthly") {
		return bmxErr.NewExclusiveFlagsError("[--hourly]", "[--monthly]")
	}

	vms, err := cmd.VirtualServerManager.ListInstances(c.IsSet("hourly"), c.IsSet("monthly"), c.String("D"), c.String("H"), c.String("d"), c.String("public-ip"), c.String("private-ip"), c.String("owner"), c.Int("c"), c.Int("m"), c.Int("n"), c.Int("o"), c.StringSlice("tag"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list virtual server instances on your account.\n")+err.Error(), 2)
	}

	if sortby == "" || sortby == "hostname" {
		sort.Sort(utils.VirtualGuestByHostname(vms))
	} else if sortby == "id" {
		sort.Sort(utils.VirtualGuestById(vms))
	} else if sortby == "domain" {
		sort.Sort(utils.VirtualGuestByDomain(vms))
	} else if sortby == "datacenter" {
		sort.Sort(utils.VirtualGuestByDatacenter(vms))
	} else if sortby == "cpu" {
		sort.Sort(utils.VirtualGuestByCPU(vms))
	} else if sortby == "memory" {
		sort.Sort(utils.VirtualGuestByMemory(vms))
	} else if sortby == "public_ip" {
		sort.Sort(utils.VirtualGuestByPrimaryIp(vms))
	} else if sortby == "private_ip" {
		sort.Sort(utils.VirtualGuestByBackendIp(vms))
	} else {
		return bmxErr.NewInvalidUsageError(T("--sortby '{{.Column}}' is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, vms)
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, vm := range vms {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(vm.Id)
		values["guid"] = utils.FormatStringPointer(vm.GlobalIdentifier)
		values["hostname"] = utils.FormatStringPointer(vm.Hostname)
		values["domain"] = utils.FormatStringPointer(vm.Domain)
		values["cpu"] = utils.FormatIntPointer(vm.MaxCpu)
		values["memory"] = utils.FormatIntPointer(vm.MaxMemory)
		values["public_ip"] = utils.FormatStringPointer(vm.PrimaryIpAddress)
		values["private_ip"] = utils.FormatStringPointer(vm.PrimaryBackendIpAddress)
		if vm.Datacenter != nil {
			values["datacenter"] = utils.FormatStringPointer(vm.Datacenter.Name)
		}
		if vm.ActiveTransaction != nil && vm.ActiveTransaction.TransactionStatus != nil {
			values["action"] = utils.FormatStringPointer(vm.ActiveTransaction.TransactionStatus.Name)
		}
		if vm.PowerState != nil {
			values["power_state"] = utils.FormatStringPointer(vm.PowerState.Name)
		}
		if vm.BillingItem != nil && vm.BillingItem.OrderItem != nil && vm.BillingItem.OrderItem.Order != nil && vm.BillingItem.OrderItem.Order.UserRecord != nil {
			values["created_by"] = utils.FormatStringPointer(vm.BillingItem.OrderItem.Order.UserRecord.Username)
		}
		values["tags"] = utils.TagRefsToString(vm.TagReferences)

		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()

	return nil
}
