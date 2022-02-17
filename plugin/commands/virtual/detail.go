package virtual

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
	VirtualServerManager managers.VirtualServerManager
}

func NewDetailCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("This command requires one argument.")
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, managers.INSTANCE_DETAIL_MASK)
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, virtualGuest)
	}

	var host datatypes.Virtual_DedicatedHost
	if virtualGuest.DedicatedHost != nil && virtualGuest.DedicatedHost.Id != nil {
		hostId := *virtualGuest.DedicatedHost.Id
		host, err = cmd.VirtualServerManager.GetDedicatedHost(hostId)
		if err != nil {
			return cli.NewExitError(T("Failed to get virtual server {{.VsID}} dedicated host: {{.HostID}}.\n",
				map[string]interface{}{"VsID": vsID, "HostID": hostId})+err.Error(), 2)
		}
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(virtualGuest.Id))
	table.Add(T("guid"), utils.FormatStringPointer(virtualGuest.GlobalIdentifier))
	table.Add(T("hostname"), utils.FormatStringPointer(virtualGuest.Hostname))
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
	table.Add(T("public ip"), utils.FormatStringPointer(virtualGuest.PrimaryIpAddress))
	table.Add(T("private ip"), utils.FormatStringPointer(virtualGuest.PrimaryBackendIpAddress))
	table.Add(T("private network"), utils.FormatBoolPointer(virtualGuest.PrivateNetworkOnlyFlag))
	table.Add(T("private cpu"), utils.FormatBoolPointer(virtualGuest.DedicatedAccountHostOnlyFlag))
	table.Add(T("created"), utils.FormatSLTimePointer(virtualGuest.CreateDate))
	table.Add(T("updated"), utils.FormatSLTimePointer(virtualGuest.ModifyDate))

	if virtualGuest.BillingItem != nil &&
		virtualGuest.BillingItem.OrderItem != nil &&
		virtualGuest.BillingItem.OrderItem.Order != nil &&
		virtualGuest.BillingItem.OrderItem.Order.UserRecord != nil {
		table.Add(T("owner"), utils.FormatStringPointer(virtualGuest.BillingItem.OrderItem.Order.UserRecord.Username))
	}

	if virtualGuest.Notes != nil && *virtualGuest.Notes != "" {
		table.Add(T("note"), utils.FormatStringPointer(virtualGuest.Notes))
	}

	if tags := virtualGuest.TagReferences; len(tags) > 0 {
		table.Add(T("tag"), utils.TagRefsToString(tags))
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

	if c.IsSet("passwords") {
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

	if c.IsSet("price") {
		var sum datatypes.Float64
		if virtualGuest.BillingItem != nil && virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount != nil {
			sum = *virtualGuest.BillingItem.NextInvoiceTotalRecurringAmount
		} else {
			sum = 0.0
		}
		for _, item := range virtualGuest.BillingItem.Children {
			if item.NextInvoiceTotalRecurringAmount != nil {
				sum += *item.NextInvoiceTotalRecurringAmount
			}
		}
		table.Add(T("price rate"), fmt.Sprintf("%.2f", sum))
	}

	table.Print()
	return nil
}

func VSDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "detail",
		Description: T("Get details for a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs detail IDENTIFIER [OPTIONS] 
	
EXAMPLE:
   ${COMMAND_NAME} sl vs details 12345678
   This command lists detailed information about virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "passwords",
				Usage: T("Show passwords (check over your shoulder!)"),
			},
			cli.BoolFlag{
				Name:  "price",
				Usage: T("Show associated prices"),
			},
			metadata.OutputFlag(),
		},
	}
}
