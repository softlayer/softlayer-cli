package managers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

const (
	DEFAULT_IP_MASK            = "hardware,virtualGuest,subnet[id,networkIdentifier,cidr,netmask,gateway,subnetType]"
	DEFAULT_GLOBALIP_MASK      = "id,ipAddress,destinationIpAddress[ipAddress,virtualGuest.fullyQualifiedDomainName,hardware.fullyQualifiedDomainName]"
	DEFAULT_SUBNET_MASK        = "id,datacenter.name,hardware.id,ipAddresses.id,networkIdentifier,networkVlan[id,networkSpace],subnetType,virtualGuests.id"
	DEFAULT_SUBNET_DETAIL_MASK = "id,broadcastAddress,cidr,datacenter.name,gateway,hardware[id,hostname,domain,primaryIpAddress,primaryBackendIpAddress],ipAddresses.id,networkIdentifier,networkVlan[id,networkSpace],subnetType,virtualGuests[id,hostname,domain,primaryIpAddress,primaryBackendIpAddress]"
	DEFAULT_VLAN_MASK          = "id,vlanNumber,name,firewallInterfaces,primaryRouter[id,fullyQualifiedDomainName,datacenter.name,hostname],hardwareCount,networkSpace,subnetCount,virtualGuestCount,totalPrimaryIpAddressCount"
	DEFAULT_VLAN_DETAIL_MASK   = "id,vlanNumber,primaryRouter[datacenterName,fullyQualifiedDomainName],firewallInterfaces," +
		"subnets[id,networkIdentifier,netmask,gateway,subnetType,usableIpAddressCount]," +
		"virtualGuests[hostname,domain,primaryIpAddress,primaryBackendIpAddress]," +
		"hardware[hostname,domain,primaryIpAddress,primaryBackendIpAddress]"
	DEFAULT_SECURITYGROUP_MASK = "id,name,description,rules[id,remoteIp,remoteGroupId,direction,ethertype,portRangeMin,portRangeMax,protocol]," +
		"networkComponentBindings[networkComponent[id,port,guest[id,hostname,primaryBackendIpAddress,primaryIpAddress]]]"
)

//Manage SoftLayer network objects: VLANs, subnets and GlobalIPs.
//See product information here: http://www.softlayer.com/networking
type NetworkManager interface {
	AddVlan(vlanType string, datacenter string, router string, name string) (datatypes.Container_Product_Order_Receipt, error)
	AddGlobalIP(version int, test bool) (datatypes.Container_Product_Order_Receipt, error)
	AddSubnet(subnetType string, quantity int, vlanID int, version int, test bool) (datatypes.Container_Product_Order_Receipt, error)
	AssignGlobalIP(globalIPID int, targetIPAddress string) (bool, error)
	UnassignGlobalIP(globalIPID int) (bool, error)
	CancelVLAN(vlanID int) error
	CancelGlobalIP(globalIPID int) error
	CancelSubnet(subnetId int) error
	EditVlan(vlanId int, name string) error
	GetSubnet(subnetId int, mask string) (datatypes.Network_Subnet, error)
	GetVlan(vlanId int, mask string) (datatypes.Network_Vlan, error)
	IPLookup(ipAddress string) (datatypes.Network_Subnet_IpAddress, error)
	ListSubnets(identifier string, datacenter string, version int, subnetType string, networkSpace string, order int, mask string) ([]datatypes.Network_Subnet, error)
	ListGlobalIPs(version int, order int) ([]datatypes.Network_Subnet_IpAddress_Global, error)
	ListVlans(datacenter string, vlanNum int, name string, order int, mask string) ([]datatypes.Network_Vlan, error)
	ListDatacenters() (map[int]string, error)
	ListRouters(dataceterId int, mask string) ([]string, error)

	ListSecurityGroups() ([]datatypes.Network_SecurityGroup, error)
	CreateSecurityGroup(name, description string) (datatypes.Network_SecurityGroup, error)
	GetSecurityGroup(groupId int, mask string) (datatypes.Network_SecurityGroup, error)
	DeleteSecurityGroup(groupId int) error
	EditSecurityGroup(groupId int, name, description string) error

	AttachSecurityGroupComponent(groupId, componentId int) error
	AttachSecurityGroupComponents(groupId int, componentIds []int) error
	DetachSecurityGroupComponent(groupId, componentId int) error
	DetachSecurityGroupComponents(groupId int, componentIds []int) error

	AddSecurityGroupRule(groupId int, remoteIp string, remoteGroup int, direction string, etherType string, portMax, portMin int, protocol string) (datatypes.Network_SecurityGroup_Rule, error)
	AddSecurityGroupRules(groupId int, rules []datatypes.Network_SecurityGroup_Rule) ([]datatypes.Network_SecurityGroup_Rule, error)
	EditSecurityGroupRule(groupId, ruleId int, remoteIp string, remoteGroup int, direction string, etherType string, portMax, portMin int, protocol string) error
	EditSecurityGroupRules(groupId int, rules []datatypes.Network_SecurityGroup_Rule) error
	ListSecurityGroupRules(grouId int) ([]datatypes.Network_SecurityGroup_Rule, error)
	RemoveSecurityGroupRule(groupId, ruleId int) error
	RemoveSecurityGroupRules(groupId int, ruleIds []int) error
	GetCancelFailureReasons(vlanId int) []string
	Route(subnetId int, typeRoute string, typeId string) (bool, error)
	ClearRoute(subnetId int) (bool, error)
	SetSubnetTags(subnetId int, tags string) (bool, error)
	SetSubnetNote(subnetId int, note string) (bool, error)
	GetIpByAddress(ipAddress string) (datatypes.Network_Subnet_IpAddress, error)
	EditSubnetIpAddress(subnetIpAddressId int, subnetIpAddressTemplate datatypes.Network_Subnet_IpAddress) (bool, error)
}

type networkManager struct {
	SubnetService        services.Network_Subnet
	VlanService          services.Network_Vlan
	IPService            services.Network_Subnet_IpAddress
	GlobalIPService      services.Network_Subnet_IpAddress_Global
	PackageService       services.Product_Package
	OrderService         services.Product_Order
	AccountService       services.Account
	BillingService       services.Billing_Item
	LocationService      services.Location_Datacenter
	SecurityGroupService services.Network_SecurityGroup
}

func NewNetworkManager(session *session.Session) *networkManager {
	return &networkManager{
		services.GetNetworkSubnetService(session),
		services.GetNetworkVlanService(session),
		services.GetNetworkSubnetIpAddressService(session),
		services.GetNetworkSubnetIpAddressGlobalService(session),
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetAccountService(session),
		services.GetBillingItemService(session),
		services.GetLocationDatacenterService(session),
		services.GetNetworkSecurityGroupService(session),
	}
}

//Add a vlan to the account
//vlanType: type of a vlan, public or private
//datacenter: short name of datacenter
//router: full qualified domain name of the router a new vlan is placed to, fcr.XXX means a public vlan router, bcr.xxx means a private vlan router
//name: vlan name
func (n networkManager) AddVlan(vlanType string, datacenter string, router string, name string) (datatypes.Container_Product_Order_Receipt, error) {
	var routerHostname string
	if router != "" {
		routerHostname = strings.Split(router, ".")[0]
		datacenter = strings.Split(router, ".")[1]
		if strings.HasPrefix(routerHostname, "fcr") {
			vlanType = "public"
		}
		if strings.HasPrefix(routerHostname, "bcr") {
			vlanType = "private"
		}
	}
	locationId, err := n.GetLocationId(datacenter)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	var routerId *int
	if router != "" {
		routers, err := n.LocationService.Id(locationId).Mask("mask[hostname,id]").GetHardwareRouters()
		if err != nil {
			return datatypes.Container_Product_Order_Receipt{}, err
		}

		for _, r := range routers {
			if r.Hostname != nil && *r.Hostname == router {
				routerId = r.Id
			}
		}

		if routerId == nil {
			return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Can't get router ID with hostname {{.HostName}}", map[string]interface{}{"HostName": router}))
		}
	}

	var vlanPriceId int
	items, err := n.PackageService.Id(0).Mask("mask[itemCategory]").GetItems()
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	for _, item := range items {
		var keyName string
		if item.KeyName != nil {
			keyName = *item.KeyName
		}
		if keyName == strings.ToUpper(vlanType)+"_NETWORK_VLAN" {
			for _, price := range item.Prices {
				if price.LocationGroupId == nil {
					vlanPriceId = *price.Id
					break
				}
			}
		}
	}
	if vlanPriceId == 0 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Failed to find price for this type of vlan."))
	}

	vlanOrder := datatypes.Container_Product_Order_Network_Vlan{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Vlan"),
			Quantity:    sl.Int(1),
			Location:    sl.String(strconv.Itoa(locationId)),
			PackageId:   sl.Int(0),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: sl.Int(vlanPriceId),
				},
			},
		},
		Name:     sl.String(name),
		RouterId: routerId,
	}
	return n.OrderService.PlaceOrder(&vlanOrder, sl.Bool(false))
}

//Edit vlan's name
//vlanId: ID of vlan
//name: name of vlan
func (n networkManager) EditVlan(vlanId int, name string) error {
	vlan := datatypes.Network_Vlan{
		Name: sl.String(name),
	}
	_, err := n.VlanService.Id(vlanId).EditObject(&vlan)
	return err
}

//Add a global IP address to the account
//version: Specifies whether this is IPv4 or IPv6
//test: If true, this will only verify the order.
func (n networkManager) AddGlobalIP(version int, test bool) (datatypes.Container_Product_Order_Receipt, error) {
	return n.AddSubnet("global", 0, 0, version, test)
}

//Order a new subnet
//subnetType: Type of subnet to add: private, public, global
//quantity: Number of IPs in the subnet
//version: 4 for IPv4, 6 for IPv6
//test: If true, this will only verify the order.
func (n networkManager) AddSubnet(subnetType string, quantity int, vlanID int, version int, test bool) (datatypes.Container_Product_Order_Receipt, error) {
	category := "sov_sec_ip_addresses_priv"
	desc := ""
	if version == 4 {
		if subnetType == "global" {
			quantity = 0
			category = "global_ipv4"
		} else if subnetType == "public" {
			category = "sov_sec_ip_addresses_pub"
		}
	} else {
		category = "static_ipv6_addresses"
		if subnetType == "global" {
			quantity = 0
			category = "global_ipv6"
			desc = "Global"
		} else if subnetType == "public" {
			desc = "Portable"
		}
	}
	//In the API, every non-server item is contained within package ID 0.
	//This means that we need to get all of the items and loop through them
	//looking for the items we need based upon the category, quantity, and
	//item description.
	var priceID int
	quantityStr := strconv.Itoa(quantity)
	items, err := n.PackageService.Id(0).Mask("itemCategory").GetItems()
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	for _, item := range items {
		var categoryCode string
		if item.ItemCategory != nil && item.ItemCategory.CategoryCode != nil {
			categoryCode = *item.ItemCategory.CategoryCode
		}
		if item.Capacity != nil {
			capacity := float64(*item.Capacity)
			if category == categoryCode && strconv.FormatFloat(capacity, 'f', 0, 64) == quantityStr &&
				(version == 4 || (version == 6 && item.Description != nil && strings.Contains(*item.Description, desc))) {
				if len(item.Prices) > 0 && item.Prices[0].Id != nil {
					priceID = *item.Prices[0].Id
					break
				}
			}
		}
	}
	if priceID == 0 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid combination specified for ordering a subnet."))
	}
	subnetOrder := datatypes.Container_Product_Order_Network_Subnet{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId: sl.Int(0),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: sl.Int(priceID),
				},
			},
			Quantity:    sl.Int(1),
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Network_Subnet"),
		},
	}
	if subnetType != "global" {
		subnetOrder.EndPointVlanId = sl.Int(vlanID)
	}
	if test {
		resp, err := n.OrderService.VerifyOrder(&subnetOrder)
		return datatypes.Container_Product_Order_Receipt{OrderDetails: &resp}, err
	}
	return n.OrderService.PlaceOrder(&subnetOrder, sl.Bool(false))
}

//Assign a global IP address to a specified target
//globalIPID: The ID of the global IP being assigned
//targetIPAddress: The IP address to assign
func (n networkManager) AssignGlobalIP(globalIPID int, targetIPAddress string) (bool, error) {
	return n.GlobalIPService.Id(globalIPID).Route(&targetIPAddress)
}

//Unassign a global IP address from a target
//globalIPID: The ID of the global IP to be cancelled.
func (n networkManager) UnassignGlobalIP(globalIPID int) (bool, error) {
	return n.GlobalIPService.Id(globalIPID).Unroute()
}

//Cancel the specifeid vlan.
//vlanID: ID of vlan to be cancelled
func (n networkManager) CancelVLAN(vlanID int) error {
	vlan, err := n.VlanService.Id(vlanID).Mask("id, billingItem.id").GetObject()
	if err != nil {
		return err
	}
	if vlan.BillingItem == nil || vlan.BillingItem.Id == nil {
		message := T("{{.TYPE}} {{.ID}} is automatically assigned and free of charge. It will automatically be removed from your account when it is empty",
			map[string]interface{}{"TYPE": "Vlan", "ID": vlanID})
		return errors.New(message)
	}
	billingID := *vlan.BillingItem.Id
	_, err = n.BillingService.Id(billingID).CancelService()
	return err
}

//Cancel the specifeid global IP address.
//globalIPID: ID of globalip to be cancelled
func (n networkManager) CancelGlobalIP(globalIPID int) error {
	IPAddress, err := n.GlobalIPService.Id(globalIPID).Mask("id, billingItem.id").GetObject()
	if err != nil {
		return err
	}
	if IPAddress.BillingItem == nil || IPAddress.BillingItem.Id == nil {
		return errors.New(T("IP address {{.ID}} can not be cancelled.", map[string]interface{}{"ID": globalIPID}))
	}
	billingID := *IPAddress.BillingItem.Id
	_, err = n.BillingService.Id(billingID).CancelService()
	return err
}

//Cancel the specifeid subnet
//subnetId: ID of subnet to be cancelled
func (n networkManager) CancelSubnet(subnetId int) error {
	subnet, err := n.GetSubnet(subnetId, "id, billingItem.id")
	if err != nil {
		return err
	}
	if subnet.BillingItem == nil || subnet.BillingItem.Id == nil {
		message := T("{{.TYPE}} {{.ID}} is automatically assigned and free of charge. It will "+
			"automatically be removed from your account when it is empty",
			map[string]interface{}{"TYPE": "Subnet", "ID": subnetId})
		return errors.New(message)
	}
	billingID := *subnet.BillingItem.Id
	_, err = n.BillingService.Id(billingID).CancelService()
	return err
}

//Returns information about a signle subnet.
//subnetId: ID of subnet
//mask: mask of properties
func (n networkManager) GetSubnet(subnetId int, mask string) (datatypes.Network_Subnet, error) {
	if mask == "" {
		mask = DEFAULT_SUBNET_DETAIL_MASK
	}
	return n.SubnetService.Id(subnetId).Mask(mask).GetObject()
}

//Returns information about a single VLAN
//vlanId: ID of vlan
//mask: mask of properties
func (n networkManager) GetVlan(vlanId int, mask string) (datatypes.Network_Vlan, error) {
	if mask == "" {
		mask = DEFAULT_VLAN_DETAIL_MASK
	}
	return n.VlanService.Id(vlanId).Mask(mask).GetObject()
}

//Looks up an IP address and returns network information about it.
//ipAddress: the ip address to be looked up
func (n networkManager) IPLookup(ipAddress string) (datatypes.Network_Subnet_IpAddress, error) {
	return n.IPService.Mask(DEFAULT_IP_MASK).GetByIpAddress(&ipAddress)
}

//List all subnets on the account
//identifier: subnet identifier to be filtered
//datacenter: datacenter shortname to be filtered
//version: v4 or v6 to be filtered
//subnetType: type of subnet to be filtered
//networkSpace: vlan space (public or private) to be filtered
//orderID: ID of order to be filtered
func (n networkManager) ListSubnets(identifier string, datacenter string, version int, subnetType string, networkSpace string, orderId int, mask string) ([]datatypes.Network_Subnet, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("subnets.id").OrderBy("ASC"))
	if identifier != "" {
		filters = append(filters, filter.Path("subnets.networkIdentifier").Eq(identifier))
	}
	if datacenter != "" {
		filters = append(filters, filter.Path("subnets.datacenter.name").Eq(datacenter))
	}
	if version != 0 {
		filters = append(filters, filter.Path("subnets.version").Eq(version))
	}
	if subnetType != "" {
		filters = append(filters, filter.Path("subnets.subnetType").Eq(subnetType))
	} else {
		filters = append(filters, filter.Path("subnets.subnetType").NotEq("GLOBAL_IP"))
	}
	if networkSpace != "" {
		filters = append(filters, filter.Path("subnets.networkVlan.networkSpace").Eq(networkSpace))
	}
	if orderId != 0 {
		filters = append(filters, filter.Path("subnets.billingItem.orderItem.order.id").Eq(orderId))
	}
	var subnetMask string
	if mask == "" {
		subnetMask = DEFAULT_SUBNET_MASK
	} else {
		subnetMask = mask
	}

	i := 0
	resourceList := []datatypes.Network_Subnet{}
	for {
		resp, err := n.AccountService.Mask(subnetMask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetSubnets()
		i++
		if err != nil {
			return []datatypes.Network_Subnet{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

//List all global ip records on current account
//version:  4 for IPv4, 6 for IPv6 to be filtered
//orderId: ID of order to be filtered
func (n networkManager) ListGlobalIPs(version int, orderId int) ([]datatypes.Network_Subnet_IpAddress_Global, error) {
	ips := []datatypes.Network_Subnet_IpAddress_Global{}
	var err error
	if version == 0 {
		if orderId == 0 {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).GetGlobalIpRecords()
		} else {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).Filter(filter.Path("globalIpRecords.billingItem.orderItem.order.id").Eq(orderId).Build()).GetGlobalIpRecords()
		}

	} else if version == 4 {
		if orderId == 0 {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).GetGlobalIpv4Records()
		} else {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).Filter(filter.Path("globalIpv4Records.billingItem.orderItem.order.id").Eq(orderId).Build()).GetGlobalIpv4Records()
		}
	} else if version == 6 {
		if orderId == 0 {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).GetGlobalIpv6Records()
		} else {
			ips, err = n.AccountService.Mask(DEFAULT_GLOBALIP_MASK).Filter(filter.Path("globalIpv6Records.billingItem.orderItem.order.id").Eq(orderId).Build()).GetGlobalIpv6Records()
		}
	}
	return ips, err
}

//List all VLANs on the account, filter by datacenter name, vlan number, and vlan name
//datacenter: datacenter shortname to be filtered
//vlanNum: number of vlan to be filtered
//name: name of vlan to be filtered
//orderId: ID of order to be filtered
func (n networkManager) ListVlans(datacenter string, vlanNum int, name string, orderId int, mask string) ([]datatypes.Network_Vlan, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("networkVlans.id").OrderBy("ASC"))
	if datacenter != "" {
		filters = append(filters, filter.Path("networkVlans.primaryRouter.datacenter.name").Eq(datacenter))
	}
	if vlanNum != 0 {
		filters = append(filters, filter.Path("networkVlans.vlanNumber").Eq(strconv.Itoa(vlanNum)))
	}
	if name != "" {
		filters = append(filters, filter.Path("networkVlans.name").Eq(name))
	}
	if orderId != 0 {
		filters = append(filters, filter.Path("networkVlans.billingItem.orderItem.order.id").Eq(orderId))
	}
	if mask == "" {
		mask = DEFAULT_VLAN_MASK
	}

	i := 0
	resourceList := []datatypes.Network_Vlan{}
	for {
		resp, err := n.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetNetworkVlans()
		i++
		if err != nil {
			return []datatypes.Network_Vlan{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

//List all datacenters, key of the map is datacenter ID, value of the map is datacenter short name
func (n networkManager) ListDatacenters() (map[int]string, error) {
	datacenters, err := n.LocationService.GetDatacenters()
	if err != nil {
		return nil, err
	}
	result := make(map[int]string)
	for _, d := range datacenters {
		if d.Id != nil && d.Name != nil {
			result[*d.Id] = *d.Name
		}
	}
	return result, nil
}

//List all routers in given datacenter
//datacenterId: ID of datacenter
func (n networkManager) ListRouters(datacenterId int, mask string) ([]string, error) {
	var routers []string
	rs, err := n.LocationService.Id(datacenterId).Mask(mask).GetHardwareRouters()
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		if r.Hostname != nil {
			routers = append(routers, *r.Hostname)
		}
	}
	return routers, nil
}

//Returns location id of datacenter for ProductOrder::placeOrder().
//location: shortname of datacenter
func (n networkManager) GetLocationId(location string) (int, error) {
	filters := filter.New(filter.Path("name").Eq(location))
	datacenters, err := n.LocationService.Mask("longName,id,name").Filter(filters.Build()).GetDatacenters()
	if err != nil {
		return 0, err
	}
	for _, datacenter := range datacenters {
		if datacenter.Name != nil && *datacenter.Name == location {
			return *datacenter.Id, nil
		}
	}
	return 0, errors.New(T("Invalid datacenter name specified."))
}

//List security groups
func (n networkManager) ListSecurityGroups() ([]datatypes.Network_SecurityGroup, error) {
	return n.SecurityGroupService.GetAllObjects()
}

//Create a security group
//name: The name of the security group
//description: The description of the security group
func (n networkManager) CreateSecurityGroup(name, description string) (datatypes.Network_SecurityGroup, error) {
	created := datatypes.Network_SecurityGroup{
		Name:        sl.String(name),
		Description: sl.String(description),
	}
	return n.SecurityGroupService.CreateObject(&created)
}

//Return the information about the given security group
//groupId: The ID of the security group
//mask: the properties to be returned
func (n networkManager) GetSecurityGroup(groupId int, mask string) (datatypes.Network_SecurityGroup, error) {
	if mask == "" {
		mask = DEFAULT_SECURITYGROUP_MASK
	}
	return n.SecurityGroupService.Id(groupId).Mask(mask).GetObject()
}

//Delete th specified security group
//groupId: The ID of the security group
func (n networkManager) DeleteSecurityGroup(groupId int) error {
	_, err := n.SecurityGroupService.Id(groupId).DeleteObject()
	return err
}

//Edit security group details
//groupId: the ID of the security group
//name: the name of the security group
//description: the description of the security group
func (n networkManager) EditSecurityGroup(groupId int, name, description string) error {
	updated := datatypes.Network_SecurityGroup{}
	if name != "" {
		updated.Name = sl.String(name)
	}
	if description != "" {
		updated.Description = sl.String(description)
	}
	_, err := n.SecurityGroupService.Id(groupId).EditObject(&updated)
	return err
}

//Attach a network component to a security group
//groupId: The ID of the security group
//componentId: the ID of the network component to attach
func (n networkManager) AttachSecurityGroupComponent(groupId, componentId int) error {
	return n.AttachSecurityGroupComponents(groupId, []int{componentId})
}

//Attach network components to a security group
//groupId: The ID of the security group
//componentIds: the IDs of the network components to attach
func (n networkManager) AttachSecurityGroupComponents(groupId int, componentIds []int) error {
	_, err := n.SecurityGroupService.Id(groupId).AttachNetworkComponents(componentIds)
	return err
}

//Detach a network component to a security group
//groupId: The ID of the security group
//componentId: the ID of the network component to detach
func (n networkManager) DetachSecurityGroupComponent(groupId, componentId int) error {
	return n.DetachSecurityGroupComponents(groupId, []int{componentId})
}

//Detach network components from a security group
//groupId: The ID of the security group
//componentIds: the IDs of the network components to detach
func (n networkManager) DetachSecurityGroupComponents(groupId int, componentIds []int) error {
	_, err := n.SecurityGroupService.Id(groupId).DetachNetworkComponents(componentIds)
	return err
}

//Add a rule to a security group
//groupId: The ID of the security group to add this rule to
//remoteIp: the remote IP or CIDR to enforce the rule on
//remoteGroup: the remote security group ID to enfore the rule on
//direction: the direction to enforce (egress or ingress)
//etherType: the etherType to enforce (IPv4 or IPv6)
//portMax: the upper port bound to enforce
//portMin: the lower port bound to enforce
//protocol: the portocol to enforce (icmp, udp, tcp)
func (n networkManager) AddSecurityGroupRule(groupId int, remoteIp string, remoteGroup int, direction string, etherType string, portMax, portMin int, protocol string) (datatypes.Network_SecurityGroup_Rule, error) {
	rule := datatypes.Network_SecurityGroup_Rule{Direction: sl.String(direction)}
	if etherType != "" {
		rule.Ethertype = sl.String(etherType)
	}
	if portMax != 0 {
		rule.PortRangeMax = sl.Int(portMax)
	}
	if portMin != 0 {
		rule.PortRangeMin = sl.Int(portMin)
	}
	if protocol != "" {
		rule.Protocol = sl.String(protocol)
	}
	if remoteIp != "" {
		rule.RemoteIp = sl.String(remoteIp)
	}
	if remoteGroup != 0 {
		rule.RemoteGroupId = sl.Int(remoteGroup)
	}
	_, err := n.AddSecurityGroupRules(groupId, []datatypes.Network_SecurityGroup_Rule{rule})
	if err != nil {
		return datatypes.Network_SecurityGroup_Rule{}, err
	}
	return rule, nil
}

//Add rules to a security group
//groupId: The ID of the security group to add rules to
//rules: rules to be added to the security group
func (n networkManager) AddSecurityGroupRules(groupId int, rules []datatypes.Network_SecurityGroup_Rule) ([]datatypes.Network_SecurityGroup_Rule, error) {
	_, err := n.SecurityGroupService.Id(groupId).AddRules(rules)
	if err != nil {
		return nil, err
	}
	return rules, err
}

//Edit a security group rule
//groupId: The ID of the security group the rule belongs to
//ruleId: The ID of the rule to edit
//remoteIp: the remote IP or CIDR to enforce the rule on
//remoteGroup: the remote security group ID to enfore the rule on
//direction: the direction to enforce (egress or ingress)
//etherType: the etherType to enforce (IPv4 or IPv6)
//portMax: the upper port bound to enforce
//portMin: the lower port bound to enforce
//protocol: the portocol to enforce (icmp, udp, tcp)
func (n networkManager) EditSecurityGroupRule(groupId, ruleId int, remoteIp string, remoteGroup int, direction string, etherType string, portMax, portMin int, protocol string) error {
	rule := datatypes.Network_SecurityGroup_Rule{Id: sl.Int(ruleId)}
	if remoteIp != "" {
		rule.RemoteIp = sl.String(remoteIp)
	}
	if remoteGroup != 0 {
		rule.RemoteGroupId = sl.Int(remoteGroup)
	}
	if direction != "" {
		rule.Direction = sl.String(direction)
	}
	if etherType != "" {
		rule.Ethertype = sl.String(etherType)
	}
	if portMax != 0 {
		rule.PortRangeMax = sl.Int(portMax)
	}
	if portMin != 0 {
		rule.PortRangeMin = sl.Int(portMin)
	}
	if protocol != "" {
		rule.Protocol = sl.String(protocol)
	}
	return n.EditSecurityGroupRules(groupId, []datatypes.Network_SecurityGroup_Rule{rule})
}

//Edit a security group rule
//groupId: The ID of the security group the rule belongs to
//rules: The rules to edit
func (n networkManager) EditSecurityGroupRules(groupId int, rules []datatypes.Network_SecurityGroup_Rule) error {
	_, err := n.SecurityGroupService.Id(groupId).EditRules(rules)
	return err
}

//List security group rules associated with a security group
//groupId: the ID of the security group to list rules for
func (n networkManager) ListSecurityGroupRules(groupId int) ([]datatypes.Network_SecurityGroup_Rule, error) {
	return n.SecurityGroupService.Id(groupId).GetRules()
}

//Remove a rule from a security group
//groupId: the ID of the security group
//ruleId: the ID of the rule to remove
func (n networkManager) RemoveSecurityGroupRule(groupId, ruleId int) error {
	return n.RemoveSecurityGroupRules(groupId, []int{ruleId})
}

//Remove rules from a security group
//groupId: the ID of the security group
//ruleIds: the list of Ids of the rules to remove
func (n networkManager) RemoveSecurityGroupRules(groupId int, ruleIds []int) error {
	_, err := n.SecurityGroupService.Id(groupId).RemoveRules(ruleIds)
	return err
}

//Calls SoftLayer_Network_Vlan::getCancelFailureReasons()
//vlanId Id for the vlan
//returns a list of strings for why a vlan can not be cancelled.
func (n networkManager) GetCancelFailureReasons(vlanId int) []string {
	reasons, err := n.VlanService.Id(vlanId).GetCancelFailureReasons()
	if err != nil {
		reasons = []string{err.Error()}
	}
	return reasons
}

//This interface allows you to change the route of your Account Owned subnets.
//subnetId int: The subnet identifier.
//typeRoute string: type value in static routing: e.g. SoftLayer_Network_Subnet_IpAddress.
//typeId string: The type identifier.
func (n networkManager) Route(subnetId int, typeRoute string, typeId string) (bool, error) {
	return n.SubnetService.Id(subnetId).Route(&typeRoute, &typeId)
}

//This interface allows you to remove the route of your Account Owned subnets.
//subnetId int: The subnet identifier.
func (n networkManager) ClearRoute(subnetId int) (bool, error) {
	return n.SubnetService.Id(subnetId).ClearRoute()
}

// Set tags of a subnet.
// subnetId int: The subnet identifier.
// tags string: Tags to be set.
func (n networkManager) SetSubnetTags(subnetId int, tags string) (bool, error) {
	return n.SubnetService.Id(subnetId).SetTags(&tags)
}

// Set note of a subnet.
// subnetId int: The subnet identifier.
// note string: Note to be set.
func (n networkManager) SetSubnetNote(subnetId int, note string) (bool, error) {
	return n.SubnetService.Id(subnetId).EditNote(&note)
}

// Get ip object by address.
// ipAddress string: ip address to find.
func (n networkManager) GetIpByAddress(ipAddress string) (datatypes.Network_Subnet_IpAddress, error) {
	return n.IPService.GetByIpAddress(&ipAddress)
}

// Set note of a subnet ip.
// ipId int: The ip identifier.
// subnetIpAddressTemplate datatypes.Network_Subnet_IpAddress: New subnet ip address templatet.
func (n networkManager) EditSubnetIpAddress(subnetIpAddressId int, subnetIpAddressTemplate datatypes.Network_Subnet_IpAddress) (bool, error) {
	return n.IPService.Id(subnetIpAddressId).EditObject(&subnetIpAddressTemplate)
}
