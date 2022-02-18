package loadbal

import (
	"bytes"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OptionsCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
	NetworkManager      managers.NetworkManager
}

func NewOptionsCommand(ui terminal.UI, lbManager managers.LoadBalancerManager, networkManager managers.NetworkManager) (cmd *OptionsCommand) {
	return &OptionsCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
		NetworkManager:      networkManager,
	}
}

type Price struct {
	KeyName      string
	DefaultPrice *datatypes.Float64
	RegionPrice  *datatypes.Float64
}

func (cmd *OptionsCommand) Run(c *cli.Context) error {
	pkgs, err := cmd.LoadBalancerManager.CreateLoadBalancerOptions()
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer product packages.")+err.Error(), 2)
	}
	cmd.UI.Say("")
	iterNumber := 0
	if c.IsSet("d") {
		for _, region := range pkgs[len(pkgs)-1].Regions {

			var dcName string
			if region.Location != nil && region.Location.Location != nil && region.Location.Location.Name != nil {
				dcName = *region.Location.Location.Name
			}
			if c.IsSet("d") && dcName != c.String("d") {
				continue
			}
			iterNumber = iterNumber + 1
			if iterNumber > 1 {
				cmd.UI.Say("-----------------------------\n")
			}
			var locationGroup []int
			if region.Location != nil && region.Location.Location != nil {
				for _, group := range region.Location.Location.Groups {
					if group.Id != nil {
						locationGroup = append(locationGroup, *group.Id)
					}
				}
			}

			table := cmd.UI.Table([]string{T("Prices:"), T("Private Subnets")})
			bufPrice := new(bytes.Buffer)
			tblPrice := terminal.NewTable(bufPrice, []string{T("Key Name"), T("Cost")})
			var prices []Price
			for _, item := range pkgs[len(pkgs)-1].Items {
				var iPrice Price
				if item.KeyName != nil {
					iPrice = Price{
						KeyName: *item.KeyName,
					}
				}
				for _, price := range item.Prices {
					if price.LocationGroupId == nil {
						iPrice.DefaultPrice = price.HourlyRecurringFee
					} else if findItemInList(price.LocationGroupId, locationGroup) {
						iPrice.RegionPrice = price.HourlyRecurringFee
					}
				}
				prices = append(prices, iPrice)
			}
			for _, price := range prices {
				if price.RegionPrice != nil {
					tblPrice.Add(price.KeyName, utils.FormatSLFloatPointerToFloat(price.RegionPrice))
				} else {
					tblPrice.Add(price.KeyName, utils.FormatSLFloatPointerToFloat(price.DefaultPrice))
				}
			}
			tblPrice.Print()

			subnets, err := cmd.NetworkManager.ListSubnets("", dcName, 0, "", "PRIVATE", 0, "networkVlan,podName,addressSpace")
			if err != nil {
				table.Add(T("Private Subnets"), T("Failed to get subnets.")+err.Error())
			} else {
				if len(subnets) > 0 {
					bufSubnet := new(bytes.Buffer)
					tblSubnet := terminal.NewTable(bufSubnet, []string{T("ID"), T("Subnet"), T("Vlan")})
					for _, subnet := range subnets {
						if subnet.SubnetType != nil && *subnet.SubnetType != "PRIMARY" && *subnet.SubnetType != "ADDITIONAL_PRIMARY" {
							continue
						}
						space := fmt.Sprintf("%s/%s", utils.FormatStringPointer(subnet.NetworkIdentifier), utils.FormatIntPointer(subnet.Cidr))
						var vlanNumber string
						if subnet.NetworkVlan != nil {
							vlanNumber = utils.FormatIntPointer(subnet.NetworkVlan.VlanNumber)
						}
						vlan := fmt.Sprintf("%s.%s", utils.FormatStringPointer(subnet.PodName), vlanNumber)
						tblSubnet.Add(utils.FormatIntPointer(subnet.Id), space, vlan)
					}
					tblSubnet.Print()
					table.Add(bufPrice.String(), bufSubnet.String())
				} else {
					table.Add(bufPrice.String(), T("Not Found"))
				}
			}
			table.Print()
		}
	} else {
		table := cmd.UI.Table([]string{T("Datacenter"), T("keyName")})
		for _, region := range pkgs[len(pkgs)-1].Regions {
			table.Add(fmt.Sprint(*region.Keyname), fmt.Sprint(*region.Location.Location.Name))
		}
		table.Print()
		fmt.Println("Use `ibmcloud sl loadbal order-options --datacenter <DC>` to find pricing information and private subnets for that specific site.")
	}

	return nil
}

func findItemInList(item *int, list []int) bool {
	if item == nil {
		return false
	}
	for _, i := range list {
		if *item == i {
			return true
		}
	}
	return false
}

func LoadbalOrderOptionsMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "order-options",
		Description: T("List options for order a load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal order-options [-d, --datacenter DATACENTER]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Show only selected datacenter, use shortname (dal13) format"),
			},
		},
	}
}
