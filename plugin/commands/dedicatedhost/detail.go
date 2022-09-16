package dedicatedhost

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
	Price                bool
	Guests               bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) *DetailCommand {
	thisCmd := &DetailCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details for a dedicated host."),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Price, "price", false, T("Show associated prices"))
	cobraCmd.Flags().BoolVar(&thisCmd.Guests, "guests", false, T("Show guests on dedicated host"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {

	hostID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Host ID")
	}

	outputFormat := cmd.GetOutputFlag()

	dedicatedhost, err := cmd.DedicatedHostManager.GetInstance(hostID, managers.DEDICATEDHOST_DETAIL_MASK)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get dedicatedhost instance: {{.HostID}}.", map[string]interface{}{"HostID": hostID}), err.Error(), 2)
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

	if cmd.Price {
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

	if cmd.Guests {
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

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
