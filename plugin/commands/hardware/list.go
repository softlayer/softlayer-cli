package hardware

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewListCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

var maskMap = map[string]string{
	"id":         "id",
	"hostname":   "hostname",
	"domain":     "domain",
	"public_ip":  "primaryIpAddress",
	"private_ip": "primaryBackendIpAddress",
	"datacenter": "datacenter",
	"status":     "hardwareStatus.status",
	"guid":       "globalIdentifier",
	"cpu":        "processorPhysicalCoreAmount",
	"memory":     "memoryCapacity",
	"os":         "operatingSystem.softwareLicense.softwareDescription.name",
	"ipmi_ip":    "networkManagementIpAddress",
	"created":    "provisionDate",
	"created_by": "billingItem.orderItem.order.userRecord.username",
	"tags":       "tagReferences",
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

	defaultColumns := []string{"id", "hostname", "domain", "public_ip", "private_ip", "datacenter", "status"}
	optionalColumns := []string{"guid", "cpu", "memory", "os", "ipmi_ip", "created", "created_by", "tags"}
	sortColumns := []string{"id", "guid", "hostname", "domain", "public_ip", "private_ip", "cpu", "memory", "os", "datacenter", "status", "ipmi_ip", "created", "created_by"}

	showColumns, err := utils.ValidateColumns(sortby, columns, defaultColumns, optionalColumns, sortColumns, c)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	hws, err := cmd.HardwareManager.ListHardware(c.StringSlice("g"), c.Int("c"), c.Int("m"), c.String("H"), c.String("D"), c.String("d"), c.Int("n"), c.String("p"), c.String("v"), c.String("owner"), c.Int("o"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware servers on your account.\n")+err.Error(), 2)
	}

	if sortby == "" || sortby == "hostname" {
		sort.Sort(utils.HardwareByHostname(hws))
	} else if sortby == "id" {
		sort.Sort(utils.HardwareById(hws))
	} else if sortby == "guid" {
		sort.Sort(utils.HardwareByGuid(hws))
	} else if sortby == "domain" {
		sort.Sort(utils.HardwareByDomain(hws))
	} else if sortby == "datacenter" {
		sort.Sort(utils.HardwareByLocation(hws))
	} else if sortby == "cpu" {
		sort.Sort(utils.HardwareByCPU(hws))
	} else if sortby == "memory" {
		sort.Sort(utils.HardwareByMemory(hws))
	} else if sortby == "os" {
		sort.Sort(utils.HardwareByOS(hws))
	} else if sortby == "public_ip" {
		sort.Sort(utils.HardwareByPublicIP(hws))
	} else if sortby == "private_ip" {
		sort.Sort(utils.HardwareByPrivateIP(hws))
	} else if sortby == "ipmi_ip" {
		sort.Sort(utils.HardwareByRemoteIP(hws))
	} else if sortby == "created" {
		sort.Sort(utils.HardwareByCreated(hws))
	} else if sortby == "created_by" {
		sort.Sort(utils.HardwareByCreatedBy(hws))
	} else if sortby == "status" {
		sort.Sort(utils.HardwareByStatus(hws))
	} else {
		return bmxErr.NewInvalidUsageError(T("--sortby '{{.Column}}' is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hws)
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, hw := range hws {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(hw.Id)
		values["guid"] = utils.FormatStringPointer(hw.GlobalIdentifier)
		values["hostname"] = utils.FormatStringPointer(hw.Hostname)
		values["domain"] = utils.FormatStringPointer(hw.Domain)
		values["cpu"] = utils.FormatUIntPointer(hw.ProcessorPhysicalCoreAmount)
		values["memory"] = utils.FormatUIntPointer(hw.MemoryCapacity)
		if hw.OperatingSystem != nil && hw.OperatingSystem.SoftwareLicense != nil && hw.OperatingSystem.SoftwareLicense.SoftwareDescription != nil {
			values["os"] = utils.FormatStringPointer(hw.OperatingSystem.SoftwareLicense.SoftwareDescription.Name)
		}
		if hw.HardwareStatus != nil {
			values["status"] = utils.FormatStringPointer(hw.HardwareStatus.Status)
		}
		values["public_ip"] = utils.FormatStringPointer(hw.PrimaryIpAddress)
		values["private_ip"] = utils.FormatStringPointer(hw.PrimaryBackendIpAddress)
		values["ipmi_ip"] = utils.FormatStringPointer(hw.NetworkManagementIpAddress)
		values["created"] = utils.FormatSLTimePointer(hw.ProvisionDate)
		if hw.Datacenter != nil {
			values["datacenter"] = utils.FormatStringPointer(hw.Datacenter.Name)
		}
		if hw.BillingItem != nil && hw.BillingItem.OrderItem != nil && hw.BillingItem.OrderItem.Order != nil && hw.BillingItem.OrderItem.Order.UserRecord != nil {
			values["created_by"] = utils.FormatStringPointer(hw.BillingItem.OrderItem.Order.UserRecord.Username)
		}
		values["tags"] = utils.TagRefsToString(hw.TagReferences)

		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()
	return nil
}
