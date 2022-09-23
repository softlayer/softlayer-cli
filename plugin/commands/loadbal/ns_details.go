package loadbal

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NetscalerDetailCommand struct {
	*metadata.SoftlayerCommand
	LoadBalancerManager managers.LoadBalancerManager
	Command             *cobra.Command
}

func NewNetscalerDetailCommand(sl *metadata.SoftlayerCommand) *NetscalerDetailCommand {
	thisCmd := &NetscalerDetailCommand{
		SoftlayerCommand:    sl,
		LoadBalancerManager: managers.NewLoadBalancerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "ns-detail " + T("IDENTIFIER"),
		Short: T("Get Netscaler details."),
		Long:  T("${COMMAND_NAME} sl  loadbal ns-detail [OPTIONS] IDENTIFIER"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *NetscalerDetailCommand) Run(args []string) error {
	netscalerID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The netscaler ID has to be a positive integer."))
	}

	outputFormat := cmd.GetOutputFlag()

	ns, err := cmd.LoadBalancerManager.GetADC(netscalerID)
	if err != nil {
		return errors.NewAPIError(T("Failed to get netscaler {{.ID}} on your account.", map[string]interface{}{"ID": netscalerID}), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add("ID", utils.FormatIntPointer(ns.Id))
	table.Add("Type", utils.FormatStringPointer(ns.Description))
	table.Add("Name", utils.FormatStringPointer(ns.Name))

	var location string
	if ns.Datacenter != nil {
		location = utils.FormatStringPointer(ns.Datacenter.LongName)
	}
	table.Add("Location", location)

	table.Add("Management IP", utils.FormatStringPointer(ns.ManagementIpAddress))

	var password string
	if ns.Password != nil {
		password = utils.FormatStringPointer(ns.Password.Password)
	}
	table.Add("Root Password", password)

	table.Add("Primary IP", utils.FormatStringPointer(ns.PrimaryIpAddress))
	table.Add("License Expiration", utils.FormatSLTimePointer(ns.LicenseExpirationDate))

	bufSubnet := new(bytes.Buffer)
	tblSubnet := terminal.NewTable(bufSubnet, []string{
		"ID",
		"Subnet",
		"Type",
		"Space",
	})
	for _, subnet := range ns.Subnets {
		tblSubnet.Add(utils.FormatIntPointer(subnet.Id), fmt.Sprintf("%s/%s", utils.FormatStringPointer(subnet.NetworkIdentifier), utils.FormatIntPointer(subnet.Cidr)), utils.FormatStringPointer(subnet.SubnetType), utils.FormatStringPointer(subnet.AddressSpace))
	}
	tblSubnet.Print()
	table.Add("Subnet", bufSubnet.String())

	bufVlan := new(bytes.Buffer)
	tblVlan := terminal.NewTable(bufVlan, []string{
		"ID",
		"Number",
	})
	for _, vlan := range ns.NetworkVlans {
		tblVlan.Add(utils.FormatIntPointer(vlan.Id), utils.FormatIntPointer(vlan.VlanNumber))
	}
	tblVlan.Print()
	table.Add("Vlans", bufVlan.String())

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
