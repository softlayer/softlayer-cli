package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

const (
	DEFAULT_LOADBAL_MASK = "loadBalancerHardware[datacenter],ipAddress,virtualServers[serviceGroups[routingMethod,routingType,services[healthChecks[type],groupReferences,ipAddress]]]"
)

type LoadBalancerManager interface {
	GetADCs() ([]datatypes.Network_Application_Delivery_Controller, error)
	GetADC(identifier int) (datatypes.Network_Application_Delivery_Controller, error)

	//loadbalancer
	CreateLoadBalancer(datacenter, name string, lbtype int, desc string, protocols []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, subnet_id int, publicSubnetIP bool) (datatypes.Container_Product_Order_Receipt, error)
	CreateLoadBalancerVerify(datacenter, name string, lbtype int, desc string, protocols []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, subnet_id int, publicSubnetIP bool) (datatypes.Container_Product_Order, error)
	CreateLoadBalancerOptions() ([]datatypes.Product_Package, error)
	CancelLoadBalancer(uuid *string) (bool, error)
	GetLoadBalancers() ([]datatypes.Network_LBaaS_LoadBalancer, error)
	GetLoadBalancer(identifier int, mask string) (datatypes.Network_LBaaS_LoadBalancer, error)
	GetLoadBalancerUUID(id int) (string, error)
	UpdateLBHealthMonitors(loadBalancerUuid *string, healthMonitorConfigurations []datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration) (datatypes.Network_LBaaS_LoadBalancer, error)

	//protocol
	AddLoadBalancerListener(loadBalancerUuid *string, protocolConfigurations []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration) (datatypes.Network_LBaaS_LoadBalancer, error)
	DeleteLoadBalancerListener(loadBalancerUuid *string, listenerUuids []string) (resp datatypes.Network_LBaaS_LoadBalancer, err error)

	//member
	DeleteLoadBalancerMember(loadBalancerUuid *string, memberUuids []string) (datatypes.Network_LBaaS_LoadBalancer, error)
	AddLoadBalancerMember(loadBalancerUuid *string, instanceInfo []datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo) (datatypes.Network_LBaaS_LoadBalancer, error)

	//l7 pool
	AddLoadBalancerL7Pool(loadBalancerUuid *string, l7Pool *datatypes.Network_LBaaS_L7Pool, l7Members []datatypes.Network_LBaaS_L7Member, l7HealthMonitor *datatypes.Network_LBaaS_L7HealthMonitor, l7SessionAffinity *datatypes.Network_LBaaS_L7SessionAffinity) (datatypes.Network_LBaaS_LoadBalancer, error)
	DeleteLoadBalancerL7Pool(identifier int) (datatypes.Network_LBaaS_LoadBalancer, error)
	GetLoadBalancerL7Pool(indentifier int) (datatypes.Network_LBaaS_L7Pool, error)
	UpdateLoadBalancerL7Pool(l7pooluuid *string, l7Pool *datatypes.Network_LBaaS_L7Pool, l7HealthMonitor *datatypes.Network_LBaaS_L7HealthMonitor, l7SessionAffinity *datatypes.Network_LBaaS_L7SessionAffinity) (datatypes.Network_LBaaS_LoadBalancer, error)

	//l7 member
	AddL7Member(l7PoolUuid *string, memberInstances []datatypes.Network_LBaaS_L7Member) (resp datatypes.Network_LBaaS_LoadBalancer, err error)
	DeleteL7Member(l7PoolUuid *string, memberUuids string) (resp datatypes.Network_LBaaS_LoadBalancer, err error)
	ListL7Members(l7PoolId int) ([]datatypes.Network_LBaaS_L7Member, error)

	//L7SessionAffinity
	GetL7SessionAffinity(l7PoolId int) (datatypes.Network_LBaaS_L7SessionAffinity, error)

	//L7pool health
	GetL7HealthMonitor(l7PoolId int) (datatypes.Network_LBaaS_L7HealthMonitor, error)

	//policies
	AddL7Policy(listenerUuid *string, policiesRules datatypes.Network_LBaaS_PolicyRule) (datatypes.Network_LBaaS_LoadBalancer, error)
	GetL7Policies(protocolId int) ([]datatypes.Network_LBaaS_L7Policy, error)
	GetL7Policy(policyId int) (datatypes.Network_LBaaS_L7Policy, error)
	DeleteL7Policy(policy int) (datatypes.Network_LBaaS_LoadBalancer, error)
	EditL7Policy(policy int, templateObject *datatypes.Network_LBaaS_L7Policy) (datatypes.Network_LBaaS_LoadBalancer, error)

	//rule
	AddL7Rule(policyUuid *string, rule datatypes.Network_LBaaS_L7Rule) (resp datatypes.Network_LBaaS_LoadBalancer, err error)
	DeleteL7Rule(policyUuid *string, ruleUuids string) (resp datatypes.Network_LBaaS_LoadBalancer, err error)
	ListL7Rule(policyid int) ([]datatypes.Network_LBaaS_L7Rule, error)
}

type loadBalancerManager struct {
	AccountService                   services.Account
	ProductService                   services.Product_Package
	DeliveryController               services.Network_Application_Delivery_Controller
	LoadBalancerService              services.Network_LBaaS_LoadBalancer
	LoadBalancerServiceHealthMonitor services.Network_LBaaS_HealthMonitor
	LoadBalancerMemberService        services.Network_LBaaS_Member
	LoadBalancerListenerService      services.Network_LBaaS_Listener
	LoadBalancerL7PoolService        services.Network_LBaaS_L7Pool
	LoadBalancerL7Member             services.Network_LBaaS_L7Member
	LoadBalancerL7PolicyService      services.Network_LBaaS_L7Policy
	LoadBalancerL7Rule               services.Network_LBaaS_L7Rule
	ProductOder                      services.Product_Order
	PackageName                      string
	OrderManager                     *orderManager
}

func NewLoadBalancerManager(session *session.Session) *loadBalancerManager {
	return &loadBalancerManager{
		services.GetAccountService(session),
		services.GetProductPackageService(session),
		services.GetNetworkApplicationDeliveryControllerService(session),
		services.GetNetworkLBaaSLoadBalancerService(session),
		services.GetNetworkLBaaSHealthMonitorService(session),
		services.GetNetworkLBaaSMemberService(session),
		services.GetNetworkLBaaSListenerService(session),
		services.GetNetworkLBaaSL7PoolService(session),
		services.GetNetworkLBaaSL7MemberService(session),
		services.GetNetworkLBaaSL7PolicyService(session),
		services.GetNetworkLBaaSL7RuleService(session),
		services.GetProductOrderService(session),
		"LBAAS",
		NewOrderManager(session),
	}
}

func (l loadBalancerManager) GetADCs() ([]datatypes.Network_Application_Delivery_Controller, error) {
	return l.AccountService.Mask("managementIpAddress,outboundPublicBandwidthUsage,primaryIpAddress,datacenter").GetApplicationDeliveryControllers()
}

func (l loadBalancerManager) GetADC(identifier int) (datatypes.Network_Application_Delivery_Controller, error) {
	return l.DeliveryController.Mask("networkVlans,password,managementIpAddress,primaryIpAddress,subnets,tagReferences,licenseExpirationDate,datacenter").Id(identifier).GetObject()
}

func (l loadBalancerManager) GetLoadBalancers() ([]datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerService.Mask("name,id,address,type,createDate,datacenter,operatingStatus,provisioningStatus").GetAllObjects()
}

func (l loadBalancerManager) GetLoadBalancer(identifier int, mask string) (datatypes.Network_LBaaS_LoadBalancer, error) {
	if mask == "" {
		mask = "healthMonitors,l7Pools,members,sslCiphers,listeners[defaultPool[healthMonitor, members, sessionAffinity],l7Policies]"
	}
	return l.LoadBalancerService.Mask(mask).Id(identifier).GetObject()
}

func (l loadBalancerManager) UpdateLBHealthMonitors(loadBalancerUuid *string, healthMonitorConfigurations []datatypes.Network_LBaaS_LoadBalancerHealthMonitorConfiguration) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerServiceHealthMonitor.UpdateLoadBalancerHealthMonitors(loadBalancerUuid, healthMonitorConfigurations)
}

func (l loadBalancerManager) DeleteLoadBalancerMember(loadBalancerUuid *string, memberUuids []string) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerMemberService.DeleteLoadBalancerMembers(loadBalancerUuid, memberUuids)
}
func (l loadBalancerManager) AddLoadBalancerMember(loadBalancerUuid *string, instanceInfo []datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerMemberService.AddLoadBalancerMembers(loadBalancerUuid, instanceInfo)
}

func (l loadBalancerManager) AddLoadBalancerListener(loadBalancerUuid *string, protocolConfigurations []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerListenerService.UpdateLoadBalancerProtocols(loadBalancerUuid, protocolConfigurations)
}

func (l loadBalancerManager) DeleteLoadBalancerListener(loadBalancerUuid *string, listenerUuids []string) (resp datatypes.Network_LBaaS_LoadBalancer, err error) {
	return l.LoadBalancerListenerService.DeleteLoadBalancerProtocols(loadBalancerUuid, listenerUuids)
}

func (l loadBalancerManager) AddLoadBalancerL7Pool(loadBalancerUuid *string, l7Pool *datatypes.Network_LBaaS_L7Pool, l7Members []datatypes.Network_LBaaS_L7Member, l7HealthMonitor *datatypes.Network_LBaaS_L7HealthMonitor, l7SessionAffinity *datatypes.Network_LBaaS_L7SessionAffinity) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PoolService.CreateL7Pool(loadBalancerUuid, l7Pool, l7Members, l7HealthMonitor, l7SessionAffinity)
}

func (l loadBalancerManager) DeleteLoadBalancerL7Pool(identifier int) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PoolService.Id(identifier).DeleteObject()
}

func (l loadBalancerManager) GetLoadBalancerL7Pool(indentifier int) (datatypes.Network_LBaaS_L7Pool, error) {
	return l.LoadBalancerL7PoolService.Id(indentifier).GetObject()
}

func (l loadBalancerManager) UpdateLoadBalancerL7Pool(l7pooluuid *string, l7Pool *datatypes.Network_LBaaS_L7Pool, l7HealthMonitor *datatypes.Network_LBaaS_L7HealthMonitor, l7SessionAffinity *datatypes.Network_LBaaS_L7SessionAffinity) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PoolService.UpdateL7Pool(l7pooluuid, l7Pool, l7HealthMonitor, l7SessionAffinity)
}

func (l loadBalancerManager) AddL7Member(l7PoolUuid *string, memberInstances []datatypes.Network_LBaaS_L7Member) (resp datatypes.Network_LBaaS_LoadBalancer, err error) {
	return l.LoadBalancerL7Member.AddL7PoolMembers(l7PoolUuid, memberInstances)
}
func (l loadBalancerManager) DeleteL7Member(l7PoolUuid *string, memberUuids string) (resp datatypes.Network_LBaaS_LoadBalancer, err error) {
	return l.LoadBalancerL7Member.DeleteL7PoolMembers(l7PoolUuid, []string{memberUuids})
}
func (l loadBalancerManager) ListL7Members(l7PoolId int) ([]datatypes.Network_LBaaS_L7Member, error) {
	return l.LoadBalancerL7PoolService.Id(l7PoolId).GetL7Members()
}

func (l loadBalancerManager) GetL7SessionAffinity(l7PoolId int) (datatypes.Network_LBaaS_L7SessionAffinity, error) {
	return l.LoadBalancerL7PoolService.Id(l7PoolId).GetL7SessionAffinity()
}

func (l loadBalancerManager) GetL7HealthMonitor(l7PoolId int) (datatypes.Network_LBaaS_L7HealthMonitor, error) {
	return l.LoadBalancerL7PoolService.Id(l7PoolId).GetL7HealthMonitor()
}

func (l loadBalancerManager) AddL7Policy(listenerUuid *string, policiesRules datatypes.Network_LBaaS_PolicyRule) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PolicyService.AddL7Policies(listenerUuid, []datatypes.Network_LBaaS_PolicyRule{policiesRules})
}

func (l loadBalancerManager) GetL7Policies(protocolId int) ([]datatypes.Network_LBaaS_L7Policy, error) {
	return l.LoadBalancerListenerService.Id(protocolId).GetL7Policies()
}

func (l loadBalancerManager) GetL7Policy(policyId int) (datatypes.Network_LBaaS_L7Policy, error) {
	return l.LoadBalancerL7PolicyService.Id(policyId).GetObject()
}

func (l loadBalancerManager) DeleteL7Policy(policy int) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PolicyService.Id(policy).DeleteObject()
}

func (l loadBalancerManager) EditL7Policy(policy int, templateObject *datatypes.Network_LBaaS_L7Policy) (datatypes.Network_LBaaS_LoadBalancer, error) {
	return l.LoadBalancerL7PolicyService.Id(policy).EditObject(templateObject)
}

func (l loadBalancerManager) AddL7Rule(policyUuid *string, rule datatypes.Network_LBaaS_L7Rule) (resp datatypes.Network_LBaaS_LoadBalancer, err error) {
	return l.LoadBalancerL7Rule.AddL7Rules(policyUuid, []datatypes.Network_LBaaS_L7Rule{rule})
}

func (l loadBalancerManager) DeleteL7Rule(policyUuid *string, ruleUuids string) (resp datatypes.Network_LBaaS_LoadBalancer, err error) {
	return l.LoadBalancerL7Rule.DeleteL7Rules(policyUuid, []string{ruleUuids})
}

func (l loadBalancerManager) ListL7Rule(policyid int) ([]datatypes.Network_LBaaS_L7Rule, error) {
	return l.LoadBalancerL7PolicyService.Id(policyid).GetL7Rules()
}

func (l loadBalancerManager) CreateLoadBalancer(datacenter, name string, lbtype int, desc string, protocols []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, subnet_id int, publicSubnetIP bool) (datatypes.Container_Product_Order_Receipt, error) {
	productPackage, err := l.OrderManager.GetPackageByKey(l.PackageName, "id,keyName,itemPrices")
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	var prices []datatypes.Product_Item_Price
	for _, price := range productPackage.ItemPrices {
		if price.LocationGroupId == nil {
			prices = append(prices, price)
		}
	}
	order := datatypes.Container_Product_Order_Network_LoadBalancer_AsAService{}
	complexType := "SoftLayer_Container_Product_Order_Network_LoadBalancer_AsAService"
	order.ComplexType = &complexType
	order.Name = &name
	order.Type = &lbtype
	order.Description = &desc
	order.Location = &datacenter
	order.PackageId = productPackage.Id
	useHourlyPricing := true
	order.UseHourlyPricing = &useHourlyPricing
	order.Prices = prices
	order.ProtocolConfigurations = protocols
	order.Subnets = []datatypes.Network_Subnet{
		datatypes.Network_Subnet{
			Id: &subnet_id,
		},
	}
	order.UseSystemPublicIpPool = &publicSubnetIP
	saveAsQuote := false
	return l.ProductOder.PlaceOrder(&order, &saveAsQuote)

}

func (l loadBalancerManager) CreateLoadBalancerVerify(datacenter, name string, lbtype int, desc string, protocols []datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, subnet_id int, publicSubnetIP bool) (datatypes.Container_Product_Order, error) {
	productPackage, err := l.OrderManager.GetPackageByKey(l.PackageName, "id,keyName,itemPrices")
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	var prices []datatypes.Product_Item_Price
	for _, price := range productPackage.ItemPrices {
		if price.LocationGroupId == nil {
			prices = append(prices, price)
		}
	}
	order := datatypes.Container_Product_Order_Network_LoadBalancer_AsAService{}
	complexType := "SoftLayer_Container_Product_Order_Network_LoadBalancer_AsAService"
	order.ComplexType = &complexType
	order.Name = &name
	order.Type = &lbtype
	order.Description = &desc
	order.Location = &datacenter
	order.PackageId = productPackage.Id
	useHourlyPricing := true
	order.UseHourlyPricing = &useHourlyPricing
	order.Prices = prices
	order.ProtocolConfigurations = protocols
	order.Subnets = []datatypes.Network_Subnet{
		datatypes.Network_Subnet{
			Id: &subnet_id,
		},
	}
	order.UseSystemPublicIpPool = &publicSubnetIP
	return l.ProductOder.VerifyOrder(&order)

}

func (l loadBalancerManager) CreateLoadBalancerOptions() ([]datatypes.Product_Package, error) {
	filters := filter.New(filter.Path("keyName").Eq(l.PackageName))
	return l.ProductService.Mask("id,keyName,name,items[prices],regions[location[location[groups]]]").Filter(filters.Build()).GetAllObjects()
}

func (l loadBalancerManager) CancelLoadBalancer(uuid *string) (bool, error) {
	return l.LoadBalancerService.CancelLoadBalancer(uuid)
}

func (l loadBalancerManager) GetLoadBalancerUUID(id int) (string, error) {
	loadBalancer, err := l.LoadBalancerService.Id(id).Mask("uuid").GetObject()
	if err != nil {
		return "", err
	}
	if loadBalancer.Uuid != nil {
		return *loadBalancer.Uuid, nil
	}
	return "", nil
}
