package dedicatedhost

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListGuestsCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
	Cpu                  int
	Tag                  []string
	Domain               string
	Hostname             string
	Memory               int
	Sortby               string
	Column              []string
}

func NewListGuestsCommand(sl *metadata.SoftlayerCommand) *ListGuestsCommand {
	thisCmd := &ListGuestsCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list-guests",
		Short: T("List Dedicated Host Guests."),
		Long: T(`${COMMAND_NAME} sl dedicatedhost list-guests IDENTIFIER[OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl dedicatedhost list-guests -d dal09 --sortby hostname 1234567
	This command list all Dedicated Host guests in the Account.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVarP(&thisCmd.Cpu, "cpu", "c", 0, T("Filter by the number of CPU cores"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "t", []string{}, T("Filter by tags."))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "d", "", T("Filter by domain portion of the FQDN."))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Filter by host portion of the FQDN."))
	cobraCmd.Flags().IntVarP(&thisCmd.Memory, "memory", "m", 0, T("Filter by Memory capacity in megabytes."))
	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "id", T("Column to sort by"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Column, "column", []string{}, T("Column to display. [Options are: guid, cpu, memory, datacenter, primary_ip, backend_ip, created_by, power_state, tags] [default: id,hostname,domain,primary_ip,backend_ip,power_state]."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

var maskListMap = map[string]string{
	"id":                "id",
	"hostname":          "hostname",
	"domain":            "domain",
	"guid":              "globalIdentifier",
	"private_ip":        "primaryBackendIpAddress",
	"public_ip":         "primaryIpAddress",
	"hourlyBillingFlag": "hourlyBillingFlag",
	"cpu":               "maxCpu",
	"memory":            "maxMemory",
	"datacenter":        "datacenter.name",
	"status":            "status",
	"power_state":       "powerState.name",
	"created_by":        "billingItem.orderItem.order.userRecord.username",
	"tags":              "tagReferences",
	"action":            "activeTransaction.transactionStatus.name",
}

func (cmd *ListGuestsCommand) Run(args []string) error {
	hostId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Host ID")
	}
	sortby := cmd.Sortby

	defaultColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter", "action"}
	optionalColumns := []string{"guid", "power_state", "created_by", "tags"}
	sortColumns := []string{"id", "hostname", "domain", "cpu", "memory", "public_ip", "private_ip", "datacenter"}

	showColumns, err := utils.ValidateColumns2(sortby, cmd.Column, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskListMap, showColumns, sortby)

	outputFormat := cmd.GetOutputFlag()

	guests, err := cmd.DedicatedHostManager.ListGuests(hostId, cmd.Cpu, cmd.Domain, cmd.Hostname, cmd.Memory, cmd.Tag, mask)
	if err != nil {
		return slErr.NewAPIError(T("Failed to list the host guest on your account."), err.Error(), 2)
	}

	if sortby == "" || sortby == "hostname" {
		sort.Sort(utils.VirtualGuestByHostname(guests))
	} else if sortby == "id" {
		sort.Sort(utils.VirtualGuestById(guests))
	} else if sortby == "domain" {
		sort.Sort(utils.VirtualGuestByDomain(guests))
	} else if sortby == "datacenter" {
		sort.Sort(utils.VirtualGuestByDatacenter(guests))
	} else if sortby == "cpu" {
		sort.Sort(utils.VirtualGuestByCPU(guests))
	} else if sortby == "memory" {
		sort.Sort(utils.VirtualGuestByMemory(guests))
	} else if sortby == "public_ip" {
		sort.Sort(utils.VirtualGuestByPrimaryIp(guests))
	} else if sortby == "private_ip" {
		sort.Sort(utils.VirtualGuestByBackendIp(guests))
	} else {
		return slErr.NewInvalidUsageError(T("--sortby '{{.Column}}' is not supported.", map[string]interface{}{"Column": sortby}))
	}

	table := cmd.UI.Table(utils.GetColumnHeader(showColumns))
	for _, vm := range guests {
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
	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}
