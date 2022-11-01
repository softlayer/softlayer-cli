package virtual

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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

	var host datatypes.Virtual_DedicatedHost
	if virtualGuest.DedicatedHost != nil && virtualGuest.DedicatedHost.Id != nil {
		hostId := *virtualGuest.DedicatedHost.Id
		host, err = cmd.VirtualServerManager.GetDedicatedHost(hostId)
		if err != nil {
			subs["HostID"] = hostId
			return slErrors.NewAPIError(T("Failed to get virtual server {{.VsID}} dedicated host: {{.HostID}}.\n", subs), err.Error(), 2)
		}
	}
	var rows [][]string
	rows = append(rows, []string{T("Name"), T("Value")})
	rows = append(rows, []string{T("ID"), utils.FormatIntPointer(virtualGuest.Id)})

	rows = append(rows, []string{T("guid"), utils.FormatStringPointer(virtualGuest.GlobalIdentifier)})
	rows = append(rows, []string{T("Hostname"), utils.FormatStringPointer(virtualGuest.Hostname)})
	rows = append(rows, []string{T("domain"), utils.FormatStringPointer(virtualGuest.Domain)})
	rows = append(rows, []string{T("fqdn"), utils.FormatStringPointer(virtualGuest.FullyQualifiedDomainName)})
	if virtualGuest.Status != nil {
		rows = append(rows, []string{T("status"), utils.FormatStringPointer(virtualGuest.Status.Name)})
	}
	if virtualGuest.PowerState != nil {
		rows = append(rows, []string{T("state"), utils.FormatStringPointer(virtualGuest.PowerState.Name)})
	}

	if virtualGuest.ActiveTransaction != nil && virtualGuest.ActiveTransaction.TransactionStatus != nil {
		rows = append(rows, []string{T("active transaction"), utils.FormatStringPointer(virtualGuest.ActiveTransaction.TransactionStatus.Name)})
	}
	if virtualGuest.Datacenter != nil {
		rows = append(rows, []string{T("datacenter"), utils.FormatStringPointer(virtualGuest.Datacenter.Name)})
	}
	if virtualGuest.OperatingSystem != nil &&
		virtualGuest.OperatingSystem.SoftwareLicense != nil &&
		virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription != nil {
		rows = append(rows, []string{T("os"), utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name)})
		rows = append(rows, []string{T("os version"), utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Version)})
	}

	rows = append(rows, []string{T("cpu cores"), utils.FormatIntPointer(virtualGuest.MaxCpu)})
	rows = append(rows, []string{T("memory"), utils.FormatIntPointer(virtualGuest.MaxMemory)})

	if localDisks != nil && len(localDisks) > 0 {
		var drivesRows [][]string
		drivesRows = append(drivesRows, []string{T("type"), T("name"), T("drive"), T("capacity")})
		for _, localDisk := range localDisks {
			diskType := "System"
			if localDisk.DiskImage != nil && localDisk.DiskImage.Description != nil {
				if strings.Contains(*localDisk.DiskImage.Description, "SWAP") {
					diskType = "Swap"
				}
			}
			drivesRows = append(drivesRows, []string{
				diskType,
				utils.FormatStringPointer(localDisk.MountType),
				utils.FormatStringPointer(localDisk.Device),
				fmt.Sprintf("%d %s", *localDisk.DiskImage.Capacity, *localDisk.DiskImage.Units)},
			)
		}
		title := "drives"
		rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, drivesRows, outputFormat)})
	} else {
		rows = append(rows, []string{"drives", "-"})
	}

	rows = append(rows, []string{T("public ip"), utils.FormatStringPointer(virtualGuest.PrimaryIpAddress)})
	rows = append(rows, []string{T("private ip"), utils.FormatStringPointer(virtualGuest.PrimaryBackendIpAddress)})
	rows = append(rows, []string{T("private network"), utils.FormatBoolPointer(virtualGuest.PrivateNetworkOnlyFlag)})
	rows = append(rows, []string{T("private cpu"), utils.FormatBoolPointer(virtualGuest.DedicatedAccountHostOnlyFlag)})

	if virtualGuest.TransientGuestFlag != nil {
		rows = append(rows, []string{T("transient"), utils.FormatBoolPointer(virtualGuest.TransientGuestFlag)})
	} else {
		rows = append(rows, []string{T("transient"), "false"})
	}

	rows = append(rows, []string{T("created"), utils.FormatSLTimePointer(virtualGuest.CreateDate)})
	rows = append(rows, []string{T("updated"), utils.FormatSLTimePointer(virtualGuest.ModifyDate)})

	lastTransaction := "-"
	if virtualGuest.LastTransaction != nil && virtualGuest.LastTransaction.TransactionGroup != nil {
		lastTransaction = fmt.Sprintf("%s (%s)", *virtualGuest.LastTransaction.TransactionGroup.Name,
			utils.FormatSLTimePointer(virtualGuest.LastTransaction.ModifyDate))
	}
	rows = append(rows, []string{T("last transaction"), lastTransaction})

	billing := "Monthly"
	if virtualGuest.HourlyBillingFlag != nil && *virtualGuest.HourlyBillingFlag {
		billing = "Hourly"
	}
	rows = append(rows, []string{T("billing"), billing})

	if virtualGuest.BillingItem != nil &&
		virtualGuest.BillingItem.OrderItem != nil &&
		virtualGuest.BillingItem.OrderItem.Preset != nil &&
		virtualGuest.BillingItem.OrderItem.Preset.KeyName != nil {
		rows = append(rows, []string{T("preset"), utils.FormatStringPointer(virtualGuest.BillingItem.OrderItem.Preset.KeyName)})
	} else {
		rows = append(rows, []string{T("preset"), "-"})
	}

	if virtualGuest.BillingItem != nil &&
		virtualGuest.BillingItem.OrderItem != nil &&
		virtualGuest.BillingItem.OrderItem.Order != nil &&
		virtualGuest.BillingItem.OrderItem.Order.UserRecord != nil {
		rows = append(rows, []string{T("owner"), utils.FormatStringPointer(virtualGuest.BillingItem.OrderItem.Order.UserRecord.Username)})
	}

	if virtualGuest.Notes != nil && *virtualGuest.Notes != "" {
		rows = append(rows, []string{T("notes"), utils.FormatStringPointer(virtualGuest.Notes)})
	} else {
		rows = append(rows, []string{T("notes"), "-"})
	}

	if virtualGuest.TagReferences != nil && len(virtualGuest.TagReferences) > 0 {
		rows = append(rows, []string{T("tags"), utils.TagRefsToString(virtualGuest.TagReferences)})
	} else {
		rows = append(rows, []string{T("tags"), "-"})
	}

	if vlans := virtualGuest.NetworkVlans; len(vlans) > 0 {
		var vlanRows [][]string
		vlanRows = append(vlanRows, []string{T("type"), T("number"), T("id")})
		for _, vlan := range vlans {
			vlanRows = append(vlanRows, []string{utils.FormatStringPointer(vlan.NetworkSpace),
				utils.FormatIntPointer(vlan.VlanNumber),
				utils.FormatIntPointer(vlan.Id)})
		}
		title := "vlans"
		rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, vlanRows, outputFormat)})
	}

	hasSecGroups := false
	var secGroupRows [][]string
	secGroupRows = append(secGroupRows, []string{T("interface"), T("id"), T("name")})
	for _, comp := range virtualGuest.NetworkComponents {
		nicType := T("public")
		if (comp.Port != nil && *comp.Port == 0) || comp.Port == nil {
			nicType = T("private")
		}
		for _, binding := range comp.SecurityGroupBindings {
			hasSecGroups = true
			secgroup := binding.SecurityGroup
			secGroupRows = append(secGroupRows, []string{nicType, utils.FormatIntPointer(secgroup.Id), utils.FormatStringPointer(secgroup.Name)})
		}
	}
	if hasSecGroups {
		title := "security groups"
		rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, secGroupRows, outputFormat)})
	}

	if virtualGuest.DedicatedHost != nil && virtualGuest.DedicatedHost.Id != nil {
		var hostRows [][]string
		hostRows = append(hostRows, []string{T("id"), T("name")})
		hostRows = append(hostRows, []string{utils.FormatIntPointer(host.Id),
			utils.FormatStringPointer(host.Name)})
		title := "dedicated host"
		rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, hostRows, outputFormat)})
	}

	if cmd.Passwords {
		if virtualGuest.OperatingSystem != nil && virtualGuest.OperatingSystem.Passwords != nil {
			var userRows [][]string
			userRows = append(userRows, []string{T("software"), T("username"), T("password")})
			for _, pwd := range virtualGuest.OperatingSystem.Passwords {
				software := ""
				if virtualGuest.OperatingSystem.SoftwareLicense != nil && virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription != nil && virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name != nil {
					software = utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Name)
				}
				userRows = append(userRows, []string{software, utils.FormatStringPointer(pwd.Username), utils.FormatStringPointer(pwd.Password)})
			}
			title := "users"
			rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, userRows, outputFormat)})
		}
	}

	if cmd.Price {
		if virtualGuest.BillingItem != nil && virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			var priceRows [][]string
			priceRows = append(priceRows, []string{T("Item"), T("CategoryCode"), T("Recurring Price")})
			totalPrice := virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount
			priceRows = append(priceRows, []string{"Total", "-", fmt.Sprintf("%.2f", *totalPrice)})
			sum := *virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount
			for _, item := range virtualGuest.BillingItem.NextInvoiceChildren {
				if item.RecurringFee != nil {
					sum += *item.RecurringFee
					priceRows = append(priceRows, []string{*item.Description, *item.CategoryCode, fmt.Sprintf("%.2f", *item.RecurringFee)})
				}
			}
			title := "Prices"
			rows = append(rows, []string{title, utils.ParseNestedTable(cmd.UI, title, priceRows, outputFormat)})
			rows = append(rows, []string{T("Price rate"), fmt.Sprintf("%.2f", sum)})
		}
	}

	utils.PrintTableWithCSV(cmd.UI, rows, outputFormat)
	return nil
}
