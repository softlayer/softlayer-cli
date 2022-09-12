package virtual

import (
	"sort"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Domain               string
	Hostname             string
	Datacenter           string
	PublicIp             string
	PrivateIp            string
	Owner                string
	Sortby               string
	Cpu                  int
	Memory               int
	Network              int
	Order                int
	Hourly               bool
	Monthly              bool
	Tag                  []string
	UserColumns          []string
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List virtual server instances on your account"),
		Long: T(`${COMMAND_NAME} sl vs list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs list --domain ibm.com --hourly --sortby memory
   This command lists all hourly-billing virtual server instances on current account filtering domain equals to "ibm.com" and sort them by memory.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Filter by domain portion of the FQDN"))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Filter by host portion of the FQDN"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter shortname"))
	cobraCmd.Flags().StringVarP(&thisCmd.PublicIp, "public-ip", "P", "", T("Filter by public IP address"))
	cobraCmd.Flags().StringVarP(&thisCmd.PrivateIp, "private-ip", "p", "", T("Filter by private IP address"))
	cobraCmd.Flags().StringVar(&thisCmd.Owner, "owner", "", T("Filtered by Id of user who owns the instances"))
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "hostname", T("Column to sort by, default is:hostname, options are:id,hostname,domain,datacenter,cpu,memory,public_ip,private_ip"))

	cobraCmd.Flags().IntVarP(&thisCmd.Cpu, "cpu", "c", 0, T("Filter by number of CPU cores"))
	cobraCmd.Flags().IntVarP(&thisCmd.Memory, "memory", "m", 0, T("Filter by memory in megabytes"))
	cobraCmd.Flags().IntVarP(&thisCmd.Network, "network", "n", 0, T("Filter by network port speed in Mbps"))
	cobraCmd.Flags().IntVarP(&thisCmd.Order, "order", "o", 0, T("Filter by ID of the order which purchased this instance"))

	cobraCmd.Flags().BoolVar(&thisCmd.Hourly, "hourly", false, T("Show only hourly instances"))
	cobraCmd.Flags().BoolVar(&thisCmd.Monthly, "monthly", false, T("Show only monthly instances"))

	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "g", []string{}, T("Filter by tags (multiple occurrence permitted)"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.UserColumns, "column", []string{}, T("Column to display. Options are: id,hostname,domain,cpu,memory,public_ip,private_ip,datacenter,action,guid,power_state,created_by,tags. This option can be specified multiple times"))

	return thisCmd
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

func (cmd *ListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	sortby := cmd.Sortby

	defaultColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter", "action"}
	optionalColumns := []string{"guid", "power_state", "created_by", "tags"}
	sortColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter"}

	showColumns, err := utils.ValidateColumns2(sortby, cmd.UserColumns, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, sortby)

	if cmd.Hourly && cmd.Monthly {
		return slErrors.NewExclusiveFlagsError("[--hourly]", "[--monthly]")
	}

	vms, err := cmd.VirtualServerManager.ListInstances(
		cmd.Hourly, cmd.Monthly, cmd.Domain, cmd.Hostname, cmd.Datacenter, cmd.PublicIp, cmd.PrivateIp, cmd.Owner,
		cmd.Cpu, cmd.Memory, cmd.Network, cmd.Order, cmd.Tag, mask,
	)

	if err != nil {
		return slErrors.NewAPIError(T("Failed to list virtual server instances on your account.\n"), err.Error(), 2)
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
		return slErrors.NewInvalidUsageError(T("--sortby '{{.Column}}' is not supported.", map[string]interface{}{"Column": sortby}))
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
