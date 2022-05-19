package managers

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	RULE_MASK             = "orderValue,action,destinationIpAddress,destinationIpSubnetMask,protocol,destinationPortRangeStart,destinationPortRangeEnd,sourceIpAddress,sourceIpSubnetMask,version,notes"
	FIREWALL_DEFAULT_MASK = "firewallNetworkComponents,networkVlanFirewall,dedicatedFirewallFlag,firewallGuestNetworkComponents,firewallInterfaces,firewallRules,highAvailabilityFirewallFlag"
)

type FirewallManager interface {
	AddVlanFirewall(vlanId int, HAenabled bool) (datatypes.Container_Product_Order_Receipt, error)
	AddStandardFirewall(serverId int, isVirtual bool) (datatypes.Container_Product_Order_Receipt, error)
	GetFirewalls() ([]datatypes.Network_Vlan, error)
	HasFirewall(vlan datatypes.Network_Vlan) bool
	GetFirewallBillingItem(fwId int, dedicated bool) (datatypes.Billing_Item, error)
	GetStandardFirewallRules(fwId int) ([]datatypes.Network_Component_Firewall_Rule, error)
	GetDedicatedFirewallRules(fwId int) ([]datatypes.Network_Vlan_Firewall_Rule, error)
	CancelFirewall(fwId int, dedicated bool) error
	GetStandardPackage(serverId int, isVirtual bool) ([]datatypes.Product_Item, error)
	GetDedicatedPackage(HAEnabled bool) ([]datatypes.Product_Item, error)
	GetFirewallPortSpeed(serverId int, isVirtual bool) (int, error)
	ParseFirewallID(inputString string) (string, int, error)
	EditDedicatedFirewallRules(firewallId int, rules []datatypes.Network_Vlan_Firewall_Rule) (datatypes.Network_Firewall_Update_Request, error)
	EditStandardFirewallRules(firewallId int, rules []datatypes.Network_Component_Firewall_Rule) (datatypes.Network_Firewall_Update_Request, error)
	GetMultiVlanFirewalls(mask string) ([]datatypes.Network_Gateway, error)
}

type firewallManager struct {
	VlanFirewallService      services.Network_Vlan_Firewall
	ComponentFirewallService services.Network_Component_Firewall
	AccountService           services.Account
	PackageService           services.Product_Package
	UpdateService            services.Network_Firewall_Update_Request
	BillingService           services.Billing_Item
	VirtualGuestService      services.Virtual_Guest
	HardwareService          services.Hardware_Server
	OrderService             services.Product_Order
}

func NewFirewallManager(session *session.Session) *firewallManager {
	return &firewallManager{
		services.GetNetworkVlanFirewallService(session),
		services.GetNetworkComponentFirewallService(session),
		services.GetAccountService(session),
		services.GetProductPackageService(session),
		services.GetNetworkFirewallUpdateRequestService(session),
		services.GetBillingItemService(session),
		services.GetVirtualGuestService(session),
		services.GetHardwareServerService(session),
		services.GetProductOrderService(session),
	}
}

//Creates a firewall for the specified vlan.
func (fw firewallManager) AddVlanFirewall(vlanId int, HAenabled bool) (datatypes.Container_Product_Order_Receipt, error) {
	packages, err := fw.GetDedicatedPackage(HAenabled)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if len(packages) == 0 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Failed to find product package for this firewall."))
	}
	firewallOrder := datatypes.Container_Product_Order_Network_Protection_Firewall_Dedicated{
		VlanId: sl.Int(vlanId),
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Protection_Firewall_Dedicated"),
			Quantity:    sl.Int(1),
			PackageId:   sl.Int(0),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: packages[0].Prices[0].Id,
				},
			},
		},
	}
	return fw.OrderService.PlaceOrder(&firewallOrder, sl.Bool(false))
}

// Creates a firewall for the specified virtual/hardware server.
func (fw firewallManager) AddStandardFirewall(serverId int, isVirtual bool) (datatypes.Container_Product_Order_Receipt, error) {
	packages, err := fw.GetStandardPackage(serverId, isVirtual)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if len(packages) == 0 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Failed to find product package for this firewall."))
	}
	firewallOrder := datatypes.Container_Product_Order_Network_Protection_Firewall{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Protection_Firewall"),
			Quantity:    sl.Int(1),
			PackageId:   sl.Int(0),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: packages[0].Prices[0].Id,
				},
			},
		},
	}
	if isVirtual {
		firewallOrder.VirtualGuests = []datatypes.Virtual_Guest{
			datatypes.Virtual_Guest{
				Id: sl.Int(serverId),
			},
		}
	} else {
		firewallOrder.Hardware = []datatypes.Hardware{
			datatypes.Hardware{
				Id: sl.Int(serverId),
			},
		}
	}
	return fw.OrderService.PlaceOrder(&firewallOrder, sl.Bool(false))
}

//Returns a list of all firewalls on the account.
func (fw firewallManager) GetFirewalls() ([]datatypes.Network_Vlan, error) {
	firewalls := []datatypes.Network_Vlan{}
	vlans, err := fw.AccountService.Mask(FIREWALL_DEFAULT_MASK).GetNetworkVlans()
	if err != nil {
		return nil, err
	}
	for _, vlan := range vlans {
		if fw.HasFirewall(vlan) {
			firewalls = append(firewalls, vlan)
		}
	}
	return firewalls, nil
}

//Returns a list of multi vlan firewalls on the account.
func (fw firewallManager) GetMultiVlanFirewalls(mask string) ([]datatypes.Network_Gateway, error) {
	if mask == "" {
		mask = "mask[networkFirewall[firewallType],insideVlans]"
	}
	firewalls := []datatypes.Network_Gateway{}
	firewallResponse, err := fw.AccountService.Mask(mask).GetNetworkGateways()
	if err != nil {
		return nil, err
	}
	for _, firewall := range firewallResponse {
		if firewall.NetworkFirewall != nil && firewall.NetworkFirewall.Id != nil {
			firewalls = append(firewalls, firewall)
		}
	}
	return firewalls, nil
}

func (fw firewallManager) HasFirewall(vlan datatypes.Network_Vlan) bool {
	return utils.IntPointertoInt(vlan.DedicatedFirewallFlag) > 0 ||
		utils.BoolPointertoBool(vlan.HighAvailabilityFirewallFlag) == true ||
		len(vlan.FirewallInterfaces) > 0 ||
		len(vlan.FirewallNetworkComponents) > 0 ||
		len(vlan.FirewallGuestNetworkComponents) > 0
}

//Retrieves the billing item of the firewall.
func (fw firewallManager) GetFirewallBillingItem(fwId int, dedicated bool) (datatypes.Billing_Item, error) {
	mask := "id,billingItem.id"
	if dedicated {
		vlanFW, err := fw.VlanFirewallService.Id(fwId).Mask(mask).GetObject()
		if err != nil {
			return datatypes.Billing_Item{}, err
		}
		if vlanFW.BillingItem != nil {
			return *vlanFW.BillingItem, nil
		}
		return datatypes.Billing_Item{}, errors.New(T("Billing item not found"))
	}
	compFW, err := fw.ComponentFirewallService.Id(fwId).Mask(mask).GetObject()
	if err != nil {
		return datatypes.Billing_Item{}, err
	}
	if compFW.BillingItem != nil {
		return *compFW.BillingItem, nil
	}
	return datatypes.Billing_Item{}, errors.New(T("Billing item not found"))
}

//Get the rules of a standard firewall.
func (fw firewallManager) GetStandardFirewallRules(fwId int) ([]datatypes.Network_Component_Firewall_Rule, error) {
	return fw.ComponentFirewallService.Id(fwId).Mask(RULE_MASK).GetRules()
}

//Get the rules of a dedicated firewall.
func (fw firewallManager) GetDedicatedFirewallRules(fwId int) ([]datatypes.Network_Vlan_Firewall_Rule, error) {
	return fw.VlanFirewallService.Id(fwId).Mask(RULE_MASK).GetRules()
}

//Cancels the specified firewall.
func (fw firewallManager) CancelFirewall(fwId int, dedicated bool) error {
	firewallBillingItem, err := fw.GetFirewallBillingItem(fwId, dedicated)
	if err != nil {
		return err
	}
	if firewallBillingItem.Id == nil {
		return errors.New(T("Failed to find billing item for firewall: {{.ID}}.", map[string]interface{}{"ID": fwId}))
	}
	var billingId int
	if firewallBillingItem.Id != nil {
		billingId = *firewallBillingItem.Id
	}
	_, err = fw.BillingService.Id(billingId).CancelService()
	return err
}

//Retrieves the standard firewall package for the virtual server or hardware server.
func (fw firewallManager) GetStandardPackage(serverId int, isVirtual bool) ([]datatypes.Product_Item, error) {
	firewallPortSpeed, err := fw.GetFirewallPortSpeed(serverId, isVirtual)
	if err != nil {
		return []datatypes.Product_Item{}, err
	}
	value := fmt.Sprintf("%dMbps Hardware Firewall", firewallPortSpeed)
	filters := filter.New(filter.Path("items.description").Eq(value))
	return fw.PackageService.Id(0).Filter(filters.Build()).GetItems()
}

//Retrieves the dedicated firewall package.
func (fw firewallManager) GetDedicatedPackage(HAEnabled bool) ([]datatypes.Product_Item, error) {
	fwFilter := "Hardware Firewall (Dedicated)"
	hafwFilter := "Hardware Firewall (High Availability)"
	filters := filter.New()
	if HAEnabled {
		filters = append(filters, filter.Path("items.description").Eq(hafwFilter))
	} else {
		filters = append(filters, filter.Path("items.description").Eq(fwFilter))
	}
	return fw.PackageService.Id(0).Filter(filters.Build()).GetItems()
}

//Determines the appropriate speed for a firewall.
func (fw firewallManager) GetFirewallPortSpeed(serverId int, isVirtual bool) (int, error) {
	firewallPortSpeed := 0
	if isVirtual {
		mask := "primaryNetworkComponent.maxSpeed"
		vs, err := fw.VirtualGuestService.Id(serverId).Mask(mask).GetObject()
		if err != nil {
			return 0, err
		}
		if vs.PrimaryNetworkComponent != nil && vs.PrimaryNetworkComponent.MaxSpeed != nil {
			firewallPortSpeed = *vs.PrimaryNetworkComponent.MaxSpeed
		}
	} else {
		mask := "id,maxSpeed,networkComponentGroup.networkComponents"
		networkComponents, err := fw.HardwareService.Id(serverId).Mask(mask).GetFrontendNetworkComponents()
		if err != nil {
			return 0, err
		}
		var ungrouped, grouped []int
		for _, comp := range networkComponents {
			if comp.NetworkComponentGroup == nil {
				if comp.MaxSpeed != nil {
					ungrouped = append(ungrouped, *comp.MaxSpeed)
				}
			} else {
				groupSum := 0
				//For each group, sum the maxSpeeds of each compoment in the group. Put the sum for each in a new list
				if comp.NetworkComponentGroup.NetworkComponents != nil {
					for _, groupComp := range comp.NetworkComponentGroup.NetworkComponents {
						if groupComp.MaxSpeed != nil {
							groupSum += *groupComp.MaxSpeed
						}
					}
				}
				grouped = append(grouped, groupSum)
			}
		}
		//find the maxspeed from ungrouped components
		maxUngrouped := 0
		for _, speed := range ungrouped {
			maxUngrouped = int(math.Max(float64(maxUngrouped), float64(speed)))
		}
		//find the maxspeed from grouped components
		maxGrouped := 0
		for _, speed := range grouped {
			maxGrouped = int(math.Max(float64(maxGrouped), float64(speed)))
		}
		//find the maxspeed from both ungrouped and grouped
		firewallPortSpeed = int(math.Max(float64(maxUngrouped), float64(maxGrouped)))
	}
	if firewallPortSpeed == 0 {
		return 0, errors.New(T("Unable to find network port speed for this device."))
	}
	return firewallPortSpeed, nil
}

func (fw firewallManager) ParseFirewallID(inputString string) (string, int, error) {
	keyvalue := strings.Split(inputString, ":")
	if len(keyvalue) != 2 {
		return "", 0, errors.New(T("Invalid ID {{.ID}}: ID should be of the form xxx:yyy, xxx is the type of the firewall, yyy is the positive integer ID.", map[string]interface{}{"ID": inputString}))
	}
	firewallType := keyvalue[0]
	if firewallType != "vs" && firewallType != "vlan" && firewallType != "server" {
		return "", 0, errors.New(T("Invalid firewall type {{.Type}}: firewall type should be either vlan, vs or server.", map[string]interface{}{"Type": firewallType}))
	}
	firewallID, err := strconv.Atoi(keyvalue[1])
	if err != nil {
		return "", 0, errors.New(T("Invalid ID {{.ID}}: ID should be of the form xxx:yyy, xxx is the type of the firewall, yyy is the positive integer ID.", map[string]interface{}{"ID": inputString}))
	}
	return firewallType, firewallID, nil
}

//Edit the rules for dedicated firewall.
func (fw firewallManager) EditDedicatedFirewallRules(firewallId int, rules []datatypes.Network_Vlan_Firewall_Rule) (datatypes.Network_Firewall_Update_Request, error) {
	mask := "networkVlan.firewallInterfaces.firewallContextAccessControlLists"
	firewall, err := fw.VlanFirewallService.Id(firewallId).Mask(mask).GetObject()
	if err != nil {
		return datatypes.Network_Firewall_Update_Request{}, err
	}
	networkVlan := firewall.NetworkVlan
	var fwlCtxAclId int
	for _, fwl := range networkVlan.FirewallInterfaces {
		if fwl.Name != nil && *fwl.Name == "inside" {
			continue
		}
		for _, controlList := range fwl.FirewallContextAccessControlLists {
			if controlList.Direction != nil && *controlList.Direction == "out" {
				continue
			}
			fwlCtxAclId = *controlList.Id
		}
	}
	template := datatypes.Network_Firewall_Update_Request{
		FirewallContextAccessControlListId: sl.Int(fwlCtxAclId),
		Rules:                              VlanRulesToUpdateRequestRules(rules),
	}
	return fw.UpdateService.CreateObject(&template)
}

//Edit the rules for standard firewall.
func (fw firewallManager) EditStandardFirewallRules(firewallId int, rules []datatypes.Network_Component_Firewall_Rule) (datatypes.Network_Firewall_Update_Request, error) {
	template := datatypes.Network_Firewall_Update_Request{
		NetworkComponentFirewallId: sl.Int(firewallId),
		Rules:                      ComponentRulesToUpdateRequestRules(rules),
	}
	return fw.UpdateService.CreateObject(&template)
}

func ComponentRulesToUpdateRequestRules(rules []datatypes.Network_Component_Firewall_Rule) []datatypes.Network_Firewall_Update_Request_Rule {
	updatedRules := []datatypes.Network_Firewall_Update_Request_Rule{}
	for _, rule := range rules {
		updatedRule := datatypes.Network_Firewall_Update_Request_Rule{
			Action:                    rule.Action,
			DestinationIpAddress:      rule.DestinationIpAddress,
			DestinationIpCidr:         rule.DestinationIpCidr,
			DestinationIpSubnetMask:   rule.DestinationIpSubnetMask,
			DestinationPortRangeEnd:   rule.DestinationPortRangeEnd,
			DestinationPortRangeStart: rule.DestinationPortRangeStart,
			Notes:                     rule.Notes,
			OrderValue:                rule.OrderValue,
			Protocol:                  rule.Protocol,
			SourceIpAddress:           rule.SourceIpAddress,
			SourceIpCidr:              rule.SourceIpCidr,
			SourceIpSubnetMask:        rule.SourceIpSubnetMask,
			Version:                   rule.Version,
		}
		updatedRules = append(updatedRules, updatedRule)
	}
	return updatedRules
}

func VlanRulesToUpdateRequestRules(rules []datatypes.Network_Vlan_Firewall_Rule) []datatypes.Network_Firewall_Update_Request_Rule {
	updatedRules := []datatypes.Network_Firewall_Update_Request_Rule{}
	for _, rule := range rules {
		updatedRule := datatypes.Network_Firewall_Update_Request_Rule{
			Action:                    rule.Action,
			DestinationIpAddress:      rule.DestinationIpAddress,
			DestinationIpCidr:         rule.DestinationIpCidr,
			DestinationIpSubnetMask:   rule.DestinationIpSubnetMask,
			DestinationPortRangeEnd:   rule.DestinationPortRangeEnd,
			DestinationPortRangeStart: rule.DestinationPortRangeStart,
			Notes:                     rule.Notes,
			OrderValue:                rule.OrderValue,
			Protocol:                  rule.Protocol,
			SourceIpAddress:           rule.SourceIpAddress,
			SourceIpCidr:              rule.SourceIpCidr,
			SourceIpSubnetMask:        rule.SourceIpSubnetMask,
			Version:                   rule.Version,
		}
		updatedRules = append(updatedRules, updatedRule)
	}
	return updatedRules
}
