package virtual

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Passwords            bool
	Price                bool
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details for a virtual server instance"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVar(&thisCmd.Passwords, "passwords", false, T("Show passwords (check over your shoulder!)"))
	cobraCmd.Flags().BoolVar(&thisCmd.Price, "price", false, T("Show associated prices"))
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, managers.INSTANCE_DETAIL_MASK)
	subs := map[string]interface{}{
		"VsID":   vsID,
		"ID":     vsID,
		"HostID": 0,
	}
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}

	localDisks, err := cmd.VirtualServerManager.GetLocalDisks(vsID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get the local disks detail for the virtual server {{.ID}}.\n", subs), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, virtualGuest)
	}

	var host datatypes.Virtual_DedicatedHost
	if virtualGuest.DedicatedHost != nil && virtualGuest.DedicatedHost.Id != nil {
		hostId := *virtualGuest.DedicatedHost.Id
		host, err = cmd.VirtualServerManager.GetDedicatedHost(hostId)
		if err != nil {
			subs["HostID"] = hostId
			return slErrors.NewAPIError(T("Failed to get virtual server {{.VsID}} dedicated host: {{.HostID}}.\n", subs), err.Error(), 2)
		}
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(virtualGuest.Id))
	table.Add(T("guid"), utils.FormatStringPointer(virtualGuest.GlobalIdentifier))
	table.Add(T("Hostname"), utils.FormatStringPointer(virtualGuest.Hostname))
	table.Add(T("domain"), utils.FormatStringPointer(virtualGuest.Domain))
	table.Add(T("fqdn"), utils.FormatStringPointer(virtualGuest.FullyQualifiedDomainName))
	if virtualGuest.Status != nil {
		table.Add(T("status"), utils.FormatStringPointer(virtualGuest.Status.Name))
	}
	if virtualGuest.PowerState != nil {
		table.Add(T("state"), utils.FormatStringPointer(virtualGuest.PowerState.Name))
	}

	if virtualGuest.ActiveTransaction != nil && virtualGuest.ActiveTransaction.TransactionStatus != nil {
		table.Add(T("active transaction"), utils.FormatStringPointer(virtualGuest.ActiveTransaction.TransactionStatus.Name))
	}
	if virtualGuest.Datacenter != nil {
		table.Add(T("datacenter"), utils.FormatStringPointer(virtualGuest.Datacenter.Name))
	}
	if virtualGuest.OperatingSystem != nil &&
		virtualGuest.OperatingSystem.SoftwareLicense != nil &&
		virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription != nil {
		table.Add(T("os"), utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name))
		table.Add(T("os version"), utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Version))
	}

	table.Add(T("cpu cores"), utils.FormatIntPointer(virtualGuest.MaxCpu))
	table.Add(T("memory"), utils.FormatIntPointer(virtualGuest.MaxMemory))

	if localDisks != nil && len(localDisks) > 0 {
		buf := new(bytes.Buffer)
		drivesTable := terminal.NewTable(buf, []string{T("type"), T("name"), T("drive"), T("capacity")})
		for _, localDisk := range localDisks {
			diskType := "System"
			if localDisk.DiskImage != nil && localDisk.DiskImage.Description != nil {
				if strings.Contains(*localDisk.DiskImage.Description, "SWAP") {
					diskType = "Swap"
				}
			}
			drivesTable.Add(
				diskType,
				utils.FormatStringPointer(localDisk.MountType),
				utils.FormatStringPointer(localDisk.Device),
				fmt.Sprintf("%d %s", *localDisk.DiskImage.Capacity, *localDisk.DiskImage.Units),
			)
		}
		drivesTable.Print()
		table.Add("drives", buf.String())
	} else {
		table.Add("drives", "-")
	}

	table.Add(T("public ip"), utils.FormatStringPointer(virtualGuest.PrimaryIpAddress))
	table.Add(T("private ip"), utils.FormatStringPointer(virtualGuest.PrimaryBackendIpAddress))
	table.Add(T("private network"), utils.FormatBoolPointer(virtualGuest.PrivateNetworkOnlyFlag))
	table.Add(T("private cpu"), utils.FormatBoolPointer(virtualGuest.DedicatedAccountHostOnlyFlag))

	if virtualGuest.TransientGuestFlag != nil {
		table.Add(T("transient"), utils.FormatBoolPointer(virtualGuest.TransientGuestFlag))
	} else {
		table.Add(T("transient"), "false")
	}

	table.Add(T("created"), utils.FormatSLTimePointer(virtualGuest.CreateDate))
	table.Add(T("updated"), utils.FormatSLTimePointer(virtualGuest.ModifyDate))

	lastTransaction := "-"
	if virtualGuest.LastTransaction != nil && virtualGuest.LastTransaction.TransactionGroup != nil {
		lastTransaction = fmt.Sprintf("%s (%s)", *virtualGuest.LastTransaction.TransactionGroup.Name,
			utils.FormatSLTimePointer(virtualGuest.LastTransaction.ModifyDate))
	}
	table.Add(T("last transaction"), lastTransaction)

	billing := "Monthly"
	if virtualGuest.HourlyBillingFlag != nil && *virtualGuest.HourlyBillingFlag {
		billing = "Hourly"
	}
	table.Add(T("billing"), billing)

	if virtualGuest.BillingItem != nil &&
		virtualGuest.BillingItem.OrderItem != nil &&
		virtualGuest.BillingItem.OrderItem.Preset != nil &&
		virtualGuest.BillingItem.OrderItem.Preset.KeyName != nil {
		table.Add(T("preset"), utils.FormatStringPointer(virtualGuest.BillingItem.OrderItem.Preset.KeyName))
	} else {
		table.Add(T("preset"), "-")
	}

	if virtualGuest.BillingItem != nil &&
		virtualGuest.BillingItem.OrderItem != nil &&
		virtualGuest.BillingItem.OrderItem.Order != nil &&
		virtualGuest.BillingItem.OrderItem.Order.UserRecord != nil {
		table.Add(T("owner"), utils.FormatStringPointer(virtualGuest.BillingItem.OrderItem.Order.UserRecord.Username))
	}

	if virtualGuest.Notes != nil && *virtualGuest.Notes != "" {
		table.Add(T("notes"), utils.FormatStringPointer(virtualGuest.Notes))
	} else {
		table.Add(T("notes"), "-")
	}

	if virtualGuest.TagReferences != nil && len(virtualGuest.TagReferences) > 0 {
		table.Add(T("tags"), utils.TagRefsToString(virtualGuest.TagReferences))
	} else {
		table.Add(T("tags"), "-")
	}

	if vlans := virtualGuest.NetworkVlans; len(vlans) > 0 {
		buf := new(bytes.Buffer)
		vlanTable := terminal.NewTable(buf, []string{T("type"), T("number"), T("id")})
		for _, vlan := range vlans {
			vlanTable.Add(utils.FormatStringPointer(vlan.NetworkSpace),
				utils.FormatIntPointer(vlan.VlanNumber),
				utils.FormatIntPointer(vlan.Id))
		}
		vlanTable.Print()
		table.Add("vlans", buf.String())
	}

	hasSecGroups := false
	buf := new(bytes.Buffer)
	secGroupTable := terminal.NewTable(buf, []string{T("interface"), T("id"), T("name")})
	for _, comp := range virtualGuest.NetworkComponents {
		nicType := T("public")
		if (comp.Port != nil && *comp.Port == 0) || comp.Port == nil {
			nicType = T("private")
		}
		for _, binding := range comp.SecurityGroupBindings {
			hasSecGroups = true
			secgroup := binding.SecurityGroup
			secGroupTable.Add(nicType, utils.FormatIntPointer(secgroup.Id), utils.FormatStringPointer(secgroup.Name))
		}
	}
	if hasSecGroups {
		secGroupTable.Print()
		table.Add(T("security groups"), buf.String())
	}

	if virtualGuest.DedicatedHost != nil && virtualGuest.DedicatedHost.Id != nil {
		buf := new(bytes.Buffer)
		hostTable := terminal.NewTable(buf, []string{T("id"), T("name")})
		hostTable.Add(utils.FormatIntPointer(host.Id),
			utils.FormatStringPointer(host.Name))
		hostTable.Print()
		table.Add(T("dedicated host"), buf.String())
	}

	if cmd.Passwords {
		if virtualGuest.OperatingSystem != nil && virtualGuest.OperatingSystem.Passwords != nil {
			buf := new(bytes.Buffer)
			userTable := terminal.NewTable(buf, []string{T("software"), T("username"), T("password")})
			for _, pwd := range virtualGuest.OperatingSystem.Passwords {
				software := ""
				if virtualGuest.OperatingSystem.SoftwareLicense != nil && virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription != nil && virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name != nil {
					software = utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name)
				}
				userTable.Add(software, utils.FormatStringPointer(pwd.Username), utils.FormatStringPointer(pwd.Password))
			}
			userTable.Print()
			table.Add("users", buf.String())
		}
	}

	if cmd.Price {
		if virtualGuest.BillingItem != nil && virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			buf := new(bytes.Buffer)
			priceTable := terminal.NewTable(buf, []string{T("Item"), T("CategoryCode"), T("Recurring Price")})

			totalPrice := virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount
			priceTable.Add("Total", "-", fmt.Sprintf("%.2f", *totalPrice))
			for _, item := range virtualGuest.BillingItem.NextInvoiceChildren {
				if item.RecurringFee != nil {
					priceTable.Add(*item.Description, *item.CategoryCode, fmt.Sprintf("%.2f", *item.RecurringFee))
				}
			}
			priceTable.Print()
			table.Add("Prices", buf.String())
			table.Add(T("Price rate"), fmt.Sprintf("%.2f", *totalPrice))
		}
	}

	table.Print()
	return nil
}
