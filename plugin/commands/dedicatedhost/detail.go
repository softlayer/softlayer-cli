package dedicatedhost

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI                   terminal.UI
	DedicatedHostManager managers.DedicatedHostManager
}

func NewDetailCommand(ui terminal.UI, dedicatedHostManager managers.DedicatedHostManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:                   ui,
		DedicatedHostManager: dedicatedHostManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("This command requires one argument.")
	}
	hostID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Host ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	dedicatedhost, err := cmd.DedicatedHostManager.GetInstance(hostID, managers.DEDICATEDHOST_DETAIL_MASK)
	if err != nil {
		return cli.NewExitError(T("Failed to get dedicatedhost instance: {{.HostID}}.\n", map[string]interface{}{"HostID": hostID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, dedicatedhost)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(dedicatedhost.Id))
	table.Add(T("Name"), utils.FormatStringPointer(dedicatedhost.Name))
	table.Add(T("Cpu Count"), utils.FormatIntPointer(dedicatedhost.CpuCount))
	table.Add(T("Memory Capacity"), utils.FormatIntPointer(dedicatedhost.MemoryCapacity))
	table.Add(T("Disk Capacity"), utils.FormatIntPointer(dedicatedhost.DiskCapacity))
	table.Add(T("Create Date"), utils.FormatSLTimePointer(dedicatedhost.CreateDate))
	table.Add(T("Modify Date"), utils.FormatSLTimePointer(dedicatedhost.ModifyDate))
	if dedicatedhost.BackendRouter != nil && dedicatedhost.BackendRouter.Id != nil {
		table.Add(T("Router Id"), utils.FormatIntPointer(dedicatedhost.BackendRouter.Id))
	}
	if dedicatedhost.BackendRouter != nil && dedicatedhost.BackendRouter.Hostname != nil {
		table.Add(T("Router Hostname"), utils.FormatStringPointer(dedicatedhost.BackendRouter.Hostname))
	}

	if dedicatedhost.BillingItem != nil &&
		dedicatedhost.BillingItem.OrderItem != nil &&
		dedicatedhost.BillingItem.OrderItem.Order != nil &&
		dedicatedhost.BillingItem.OrderItem.Order.UserRecord != nil &&
		dedicatedhost.BillingItem.OrderItem.Order.UserRecord.Username != nil {
		table.Add(T("Owner"), utils.FormatStringPointer(dedicatedhost.BillingItem.OrderItem.Order.UserRecord.Username))
	}
	if dedicatedhost.Datacenter != nil && dedicatedhost.Datacenter.Name != nil {
		table.Add(T("Datacenter"), utils.FormatStringPointer(dedicatedhost.Datacenter.Name))
	}

	if c.IsSet("price") {
		var sum datatypes.Float64
		if dedicatedhost.BillingItem != nil && dedicatedhost.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			sum = *dedicatedhost.BillingItem.NextInvoiceTotalRecurringAmount
		} else {
			sum = 0.0
		}
		for _, item := range dedicatedhost.BillingItem.Children {
			if item.NextInvoiceTotalRecurringAmount != nil {
				sum += *item.NextInvoiceTotalRecurringAmount
			}
		}
		table.Add(T("Price Rate"), fmt.Sprintf("%.2f", sum))
	}

	if c.IsSet("guests") {
		if dedicatedhost.Guests != nil {
			buf := new(bytes.Buffer)
			guestTable := terminal.NewTable(buf, []string{T("Id"), T("Hostname"), T("Domain"), T("uuid")})
			for _, guest := range dedicatedhost.Guests {
				guestTable.Add(utils.FormatIntPointer(guest.Id), utils.FormatStringPointer(guest.Hostname), utils.FormatStringPointer(guest.Domain), utils.FormatStringPointer(guest.Uuid))
			}
			guestTable.Print()
			table.Add("guests", buf.String())
		}
	}

	table.Print()
	return nil
}

func DedicatedhostDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "dedicatedhost",
		Name:        "detail",
		Description: T("Get details for a dedicated host."),
		Usage: T(`${COMMAND_NAME} sl dedicatedhost detail IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost detail 1234567
   This command gets the details for a dedicated host.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "price",
				Usage: T("Show associated prices"),
			},
			cli.BoolFlag{
				Name:  "guests",
				Usage: T("Show guests on dedicated host"),
			},
			metadata.OutputFlag(),
		},
	}
}
