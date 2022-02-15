package dedicatedhost

import (
	"sort"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListGuestsCommand struct {
	UI                   terminal.UI
	DedicatedHostManager managers.DedicatedHostManager
}

func NewListGuestsCommand(ui terminal.UI, dedicatedHostManager managers.DedicatedHostManager) (cmd *ListGuestsCommand) {
	return &ListGuestsCommand{
		UI:                   ui,
		DedicatedHostManager: dedicatedHostManager,
	}
}

func DedicatedhostListGuestsMetaData() cli.Command {
	return cli.Command{
		Category:    "dedicatedhost",
		Name:        "list-guests",
		Description: T("List Dedicated Host Guests."),
		Usage: T(`${COMMAND_NAME} sl dedicatedhost list-guests IDENTIFIER[OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost list-guests -d dal09 --sortby hostname 1234567
   This command list all Dedicated Host guests in the Account.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Filter by the number of CPU cores"),
			},
			cli.StringSliceFlag{
				Name:  "t,tag",
				Usage: T("Filter by tags"),
			},
			cli.StringFlag{
				Name:  "d,domain",
				Usage: T("Filter by domain portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Filter by host portion of the FQDN"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Filter by Memory capacity in megabytes"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:hostname"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. [Options are: guid, cpu, memory, datacenter, primary_ip, backend_ip, created_by, power_state, tags] [default: id,hostname,domain,primary_ip,backend_ip,power_state]"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			metadata.OutputFlag(),
		},
	}
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

func (cmd *ListGuestsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hostId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Host ID")
	}
	sortby := c.String("sortby")
	if sortby == "" {
		sortby = "id"
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

	mask := utils.GetMask(maskListMap, showColumns, sortby)

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	guests, err := cmd.DedicatedHostManager.ListGuests(hostId, c.Int("cpu"), c.String("domain"), c.String("hostname"), c.Int("memory"), c.StringSlice("tag"), mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list the host guest on your account.\n")+err.Error(), 2)
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
		return bmxErr.NewInvalidUsageError(T("--sortby '{{.Column}}' is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, guests)
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
	table.Print()

	return nil
}
