package managers

import (
	"errors"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

const (
	CONTEXT_DETAULT_MASK = "id,name,friendlyName,internalPeerIpAddress,customerPeerIpAddress,createDate"
)

// This provides helpers to manage IPSEC contexts, private and remote subnets, and NAT translations.
type IPSECManager interface {
	AddInternalSubnet(contextId, subnetId int) error
	AddRemoteSubnet(contextId, subnetId int) error
	AddServiceSubnet(contextId, subnetId int) error
	ApplyConfiguration(contextId int) error
	CreateRemoteSubnet(accountId int, networkId string, cidr int) (datatypes.Network_Customer_Subnet, error)
	CreateTranslation(contextId int, staticIp, remoteIp, note string) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error)
	DeleteRemoteSubnet(contextId, subnetId int) error
	GetTunnelContexts(order int, mask string) ([]datatypes.Network_Tunnel_Module_Context, error)
	GetTunnelContext(contextId int, mask string) (datatypes.Network_Tunnel_Module_Context, error)
	GetTranslations(contextId int) ([]datatypes.Network_Tunnel_Module_Context_Address_Translation, error)
	GetTranslation(contextId, translationId int) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error)
	RemoveInternalSubnet(contextId, subnetId int) error
	RemoveRemoteSubnet(contextId, subnetId int) error
	RemoveServiceSubnet(contextId, subnetId int) error
	RemoveTranslation(contextId, translationId int) error
	UpdateTranslation(contextId, translationId int, staticIp, remoteIp, notes string) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error)
	UpdateTunnelContext(contextId int, name, remotePeer, presharedKey, phase1Auth, phase1Crypto string, phase1Dh, phase1KeyTtl int, phase2Auth, phase2Crypto string, phase2Dh, phase2ForwareSecrecy, phase2KeyTtl int) (datatypes.Network_Tunnel_Module_Context, error)
	OrderTunnelContext(location string) (datatypes.Container_Product_Order_Receipt, error)
	CancelTunnelContext(contextId int, immediate bool, reason string) error
}

type ipsecManager struct {
	AccountService      services.Account
	ContextService      services.Network_Tunnel_Module_Context
	RemoteSubnetService services.Network_Customer_Subnet
	PackageService      services.Product_Package
	LocationService     services.Location_Datacenter
	OrderService        services.Product_Order
	BillingService      services.Billing_Item
}

func NewIPSECManager(session *session.Session) *ipsecManager {
	return &ipsecManager{
		services.GetAccountService(session),
		services.GetNetworkTunnelModuleContextService(session),
		services.GetNetworkCustomerSubnetService(session),
		services.GetProductPackageService(session),
		services.GetLocationDatacenterService(session),
		services.GetProductOrderService(session),
		services.GetBillingItemService(session),
	}
}

//Add an internal subnet to a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the internal subnet
func (i ipsecManager) AddInternalSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).AddPrivateSubnetToNetworkTunnel(&subnetId)
	return err
}

//Add a remote subnet to a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the internal subnet
func (i ipsecManager) AddRemoteSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).AddCustomerSubnetToNetworkTunnel(&subnetId)
	return err
}

//Add a service subnet to a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the internal subnet
func (i ipsecManager) AddServiceSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).AddServiceSubnetToNetworkTunnel(&subnetId)
	return err
}

//Requests network configuration for a tunnel context.
//contextId: The id-value representing the context instance
func (i ipsecManager) ApplyConfiguration(contextId int) error {
	_, err := i.ContextService.Id(contextId).ApplyConfigurationsToDevice()
	return err
}

//Creates a remote subnet on the given account
//accountId: The account identifier
//networkId: The network identifier of the remote subnet
//cidr: The CIDR value of the remote subnet
func (i ipsecManager) CreateRemoteSubnet(accountId int, networkId string, cidr int) (datatypes.Network_Customer_Subnet, error) {
	remoteSubnet := datatypes.Network_Customer_Subnet{
		AccountId:         sl.Int(accountId),
		NetworkIdentifier: sl.String(networkId),
		Cidr:              sl.Int(cidr),
	}
	return i.RemoteSubnetService.CreateObject(&remoteSubnet)
}

//Creates an address translation on a tunnel context
//staticIp: The IP address value representing the internal side of the translation entry
//remoteIp: The IP address value representing the remote side of the translation entry
//notes: The notes to supply with the translation entry
func (i ipsecManager) CreateTranslation(contextId int, staticIp, remoteIp, note string) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error) {
	translation := datatypes.Network_Tunnel_Module_Context_Address_Translation{
		CustomerIpAddress: sl.String(remoteIp),
		InternalIpAddress: sl.String(staticIp),
		Notes:             sl.String(note),
	}
	return i.ContextService.Id(contextId).CreateAddressTranslation(&translation)
}

//Deletes a remote subnet from the current account
//contextId: The id representing the tunnel context
//subnetId: The id-value representing the remote subnet
func (i ipsecManager) DeleteRemoteSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).RemoveCustomerSubnetFromNetworkTunnel(&subnetId)
	return err
}

//Retrieves the network tunnel context instance
//contextId: The id-value representing the context instance
func (i ipsecManager) GetTunnelContext(contextId int, mask string) (datatypes.Network_Tunnel_Module_Context, error) {
	filters := filter.New(filter.Path("networkTunnelContexts.id").Eq(contextId))
	if mask == "" {
		mask = CONTEXT_DETAULT_MASK
	}
	contexts, err := i.AccountService.Mask(mask).Filter(filters.Build()).GetNetworkTunnelContexts()
	if err != nil {
		return datatypes.Network_Tunnel_Module_Context{}, err
	}
	if len(contexts) == 0 {
		return datatypes.Network_Tunnel_Module_Context{}, errors.New(T("SoftLayer_Exception_ObjectNotFound: Unable to find object with id {{.ContextID}}", map[string]interface{}{"ContextID": contextId}))
	}
	return contexts[0], nil
}

//Retrieves the network tunnel contexts
func (i ipsecManager) GetTunnelContexts(orderId int, mask string) ([]datatypes.Network_Tunnel_Module_Context, error) {
	if mask == "" {
		mask = CONTEXT_DETAULT_MASK
	}
	if orderId != 0 {
		filters := filter.New()
		filters = append(filters, filter.Path("networkTunnelContexts.billingItem.orderItem.order.id").Eq(orderId))
		return i.AccountService.Mask(mask).Filter(filters.Build()).GetNetworkTunnelContexts()
	}
	return i.AccountService.Mask(mask).GetNetworkTunnelContexts()
}

//Retrieves all translation entries for a tunnel context
//contextId: The id-value representing the context instance
func (i ipsecManager) GetTranslations(contextId int) ([]datatypes.Network_Tunnel_Module_Context_Address_Translation, error) {
	mask := "addressTranslations[customerIpAddressRecord,internalIpAddressRecord]"
	context, err := i.GetTunnelContext(contextId, mask)
	if err != nil {
		return nil, err
	}
	for _, t := range context.AddressTranslations {
		remoteIp := t.CustomerIpAddressRecord
		internalIp := t.InternalIpAddressRecord
		t.CustomerIpAddress = remoteIp.IpAddress
		t.InternalIpAddress = internalIp.IpAddress
		t.CustomerIpAddressRecord = nil
		t.InternalIpAddressRecord = nil
	}
	return context.AddressTranslations, nil
}

//Retrieves a translation entry for the given id values
//contextId: The id-value representing the context instance
//translationId: The id-value representing the translation instance
func (i ipsecManager) GetTranslation(contextId, translationId int) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error) {
	translations, err := i.GetTranslations(contextId)
	if err != nil {
		return datatypes.Network_Tunnel_Module_Context_Address_Translation{}, err
	}
	for _, t := range translations {
		if t.Id != nil && *t.Id == translationId {
			return t, nil
		}
	}
	return datatypes.Network_Tunnel_Module_Context_Address_Translation{}, errors.New(T("SoftLayer_Exception_ObjectNotFound: Unable to find object with id {{.TranslationID}}", map[string]interface{}{"TranslationID": translationId}))
}

//Remove an internal subnet from a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the internal subnet
func (i ipsecManager) RemoveInternalSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).RemovePrivateSubnetFromNetworkTunnel(&subnetId)
	return err
}

//Removes a remote subnet from a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the remote subnet
func (i ipsecManager) RemoveRemoteSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).RemoveCustomerSubnetFromNetworkTunnel(&subnetId)
	return err
}

//Removes a service subnet from a tunnel context
//contextId: The id-value representing the context instance
//subnetId: The id-value representing the service subnet
func (i ipsecManager) RemoveServiceSubnet(contextId, subnetId int) error {
	_, err := i.ContextService.Id(contextId).RemoveServiceSubnetFromNetworkTunnel(&subnetId)
	return err
}

//Removes a translation entry from a tunnel context
//contextId: The id-value representing the context instance
//translationId:  The id-value representing the translation
func (i ipsecManager) RemoveTranslation(contextId, translationId int) error {
	_, err := i.ContextService.Id(contextId).DeleteAddressTranslation(&translationId)
	return err
}

//Updates an address translation entry using the given values
//contextId: The id-value representing the context instance.
//staticIp: The static IP address value to update.
//remoteIp: The remote IP address value to update.
//notes: The notes value to update.
func (i ipsecManager) UpdateTranslation(contextId, translationId int, staticIp, remoteIp, notes string) (datatypes.Network_Tunnel_Module_Context_Address_Translation, error) {
	translation, err := i.GetTranslation(contextId, translationId)
	if err != nil {
		return datatypes.Network_Tunnel_Module_Context_Address_Translation{}, err
	}
	if staticIp != "" {
		translation.InternalIpAddress = &staticIp
		translation.InternalIpAddressId = nil
	}
	if remoteIp != "" {
		translation.CustomerIpAddress = &remoteIp
		translation.CustomerIpAddressId = nil
	}
	if notes != "" {
		translation.Notes = &notes
	}
	return i.ContextService.Id(contextId).EditAddressTranslation(&translation)
}

//Updates a tunnel context using the given values
//contextId: The id-value representing the context.
//name: The friendly name value to update.
//remotePeer: The remote peer IP address value to update.
//presharedKey: The preshared key value to update.
//phase1Auth: The phase 1 authentication value to update.
//phase1Crypto: The phase 1 encryption value to update.
//phase1Dh: The phase 1 diffie hellman group value to update.
//phase1KeyTtl: The phase 1 key life value to update.
//phase2Auth: The phase 2 authentication value to update.
//phase2Crypto: The phase 2 encryption value to update.
//phase2Dh: The phase 2 diffie hellman group value to update.
//phase2ForwardSecrecy: The phase 2 perfect forward secrecy value to update.
//phase2KeyTtl: The phase 2 key life value to update.
func (i ipsecManager) UpdateTunnelContext(contextId int, name, remotePeer, presharedKey, phase1Auth, phase1Crypto string, phase1Dh, phase1KeyTtl int, phase2Auth, phase2Crypto string, phase2Dh, phase2ForwareSecrecy, phase2KeyTtl int) (datatypes.Network_Tunnel_Module_Context, error) {
	context, err := i.GetTunnelContext(contextId, "")
	if err != nil {
		return datatypes.Network_Tunnel_Module_Context{}, err
	}
	if name != "" {
		context.FriendlyName = &name
	}
	if remotePeer != "" {
		context.CustomerPeerIpAddress = &remotePeer
	}
	if presharedKey != "" {
		context.PresharedKey = &presharedKey
	}
	if phase1Auth != "" {
		context.PhaseOneAuthentication = &phase1Auth
	}
	if phase1Crypto != "" {
		context.PhaseOneEncryption = &phase1Crypto
	}
	if phase1Dh != 0 {
		context.PhaseOneDiffieHellmanGroup = &phase1Dh
	}
	if phase1KeyTtl != 0 {
		context.PhaseOneKeylife = &phase1KeyTtl
	}
	if phase2Auth != "" {
		context.PhaseTwoAuthentication = &phase2Auth
	}
	if phase2Crypto != "" {
		context.PhaseTwoEncryption = &phase2Crypto
	}
	if phase2Dh != 0 {
		context.PhaseTwoDiffieHellmanGroup = &phase2Dh
	}
	if phase2ForwareSecrecy != 0 {
		context.PhaseTwoPerfectForwardSecrecy = &phase2ForwareSecrecy
	}
	if phase2KeyTtl != 0 {
		context.PhaseTwoKeylife = &phase2KeyTtl
	}
	_, err = i.ContextService.Id(contextId).EditObject(&context)
	if err != nil {
		return datatypes.Network_Tunnel_Module_Context{}, err
	}
	return context, nil
}

func (i ipsecManager) OrderTunnelContext(location string) (datatypes.Container_Product_Order_Receipt, error) {
	categoryCode := "network_tunnel"
	filters := filter.New()
	filters = append(filters, filter.Path("items.itemCategory.categoryCode").Eq(categoryCode))
	packageItems, err := i.PackageService.Id(0).Filter(filters.Build()).GetItems()
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if len(packageItems) == 0 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Failed to find product package for this IPSEC."))
	}
	locationId, err := GetLocationId(i.LocationService, location)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	order := datatypes.Container_Product_Order_Network_Tunnel_Ipsec{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("Container_Product_Order_Network_Tunnel_Ipsec"),
			PackageId:   sl.Int(0),
			Location:    sl.String(strconv.Itoa(locationId)),
			Quantity:    sl.Int(1),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: packageItems[0].Prices[0].Id,
				},
			},
		},
	}
	return i.OrderService.PlaceOrder(&order, sl.Bool(false))
}

func (i ipsecManager) CancelTunnelContext(contextId int, immediate bool, reason string) error {
	context, err := i.GetTunnelContext(contextId, "id,billingItem.id")
	if err != nil {
		return err
	}
	if context.BillingItem == nil || context.BillingItem.Id == nil {
		return errors.New(T("No billing item is found to cancel."))
	}
	billitemId := *context.BillingItem.Id
	_, err = i.BillingService.Id(billitemId).CancelItem(sl.Bool(immediate), sl.Bool(true), sl.String(reason), sl.String(""))
	return err
}
