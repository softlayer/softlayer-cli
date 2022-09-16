package hardware

import (
	"sort"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Cpu             int
	Domain          string
	Hostname        string
	Datacenter      string
	Memory          int
	Network         int
	Tag             []string
	PublicIp        string
	PrivateIp       string
	Order           int
	Owner           string
	Sortby          string
	Column          []string
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List hardware servers"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVarP(&thisCmd.Cpu, "cpu", "c", 0, T("Filter by number of CPU cores"))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Filter by domain"))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Filter by hostname"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter"))
	cobraCmd.Flags().IntVarP(&thisCmd.Memory, "memory", "m", 0, T("Filter by memory in gigabytes"))
	cobraCmd.Flags().IntVarP(&thisCmd.Network, "network", "n", 0, T("Filter by network port speed in Mbps"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "g", []string{}, T("Filter by tags, multiple occurrence allowed"))
	cobraCmd.Flags().StringVarP(&thisCmd.PublicIp, "public-ip", "p", "", T("Filter by public IP address"))
	cobraCmd.Flags().StringVarP(&thisCmd.PrivateIp, "private-ip", "v", "", T("Filter by private IP address"))
	cobraCmd.Flags().IntVarP(&thisCmd.Order, "order", "o", 0, T("Filter by ID of the order which purchased hardware server"))
	cobraCmd.Flags().StringVar(&thisCmd.Owner, "owner", "", T("Filter by ID of the owner"))
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by, default:hostname, option:id,guid,hostname,domain,public_ip,private_ip,cpu,memory,os,datacenter,status,ipmi_ip,created,created_by"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Column, "column", []string{}, T("Column to display,  options are: id,hostname,domain,public_ip,private_ip,datacenter,status,guid,cpu,memory,os,ipmi_ip,created,created_by,tags. This option can be specified multiple times"))

	thisCmd.Command = cobraCmd
	return thisCmd
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

func (cmd *ListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	sortby := cmd.Sortby
	if sortby == "" {
		sortby = "hostname"
	}

	columns := cmd.Column

	defaultColumns := []string{"id", "hostname", "domain", "public_ip", "private_ip", "datacenter", "status"}
	optionalColumns := []string{"guid", "cpu", "memory", "os", "ipmi_ip", "created", "created_by", "tags"}
	sortColumns := []string{"id", "guid", "hostname", "domain", "public_ip", "private_ip", "cpu", "memory", "os", "datacenter", "status", "ipmi_ip", "created", "created_by"}

	showColumns, err := utils.ValidateColumns2(sortby, columns, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	hws, err := cmd.HardwareManager.ListHardware(cmd.Tag, cmd.Cpu, cmd.Memory, cmd.Hostname, cmd.Domain, cmd.Datacenter, cmd.Network, cmd.PublicIp, cmd.PrivateIp, cmd.Owner, cmd.Order, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware servers on your account.\n"), err.Error(), 2)
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
