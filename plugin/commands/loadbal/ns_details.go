package loadbal

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type NetscalerDetailCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewNetscalerDetailCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *NetscalerDetailCommand) {
	return &NetscalerDetailCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *NetscalerDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("Netscaler ID is required.")
	}

	netscalerID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("The netscaler ID has to be a positive integer."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	ns, err := cmd.LoadBalancerManager.GetADC(netscalerID)
	if err != nil {
		return cli.NewExitError(T("Failed to get netscaler {{.ID}} on your account.", map[string]interface{}{"ID": netscalerID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, ns)
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

	table.Print()
	return nil
}

func LoadbalNetscalerDetailMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "ns-detail",
		Description: T("Get Netscaler details."),
		Usage:       "${COMMAND_NAME} sl  loadbal ns-detail [OPTIONS] IDENTIFIER",
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
