package managers

import (
	"errors"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type StatusInfo struct {
	Id     int
	Fqdn   string
	Status string
}

//Manages SoftLayer Dedicated host.
type DedicatedHostManager interface {
	ListGuests(identifier int, cpu int, domain string, hostname string, memory int, tags []string, mask string) ([]datatypes.Virtual_Guest, error)
	GenerateOrderTemplate(size, hostname, domain, datacenter string, billing string, routerId int) (datatypes.Container_Product_Order_Virtual_DedicatedHost, error)
	VerifyInstanceCreation(orderTemplate datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order, error)
	OrderInstance(orderTemplate datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order_Receipt, error)
	CancelGuests(id int) ([]StatusInfo, error)
}

type dedicatedhostManager struct {
	AccountService       services.Account
	VirtualDedicatedHost services.Virtual_DedicatedHost
	PackageService       services.Product_Package
	OrderService         services.Product_Order
	VirtualGuestService  services.Virtual_Guest
}

func NewDedicatedhostManager(session *session.Session) *dedicatedhostManager {
	return &dedicatedhostManager{
		services.GetAccountService(session),
		services.GetVirtualDedicatedHostService(session),
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetVirtualGuestService(session),
	}
}

//Cancel an instance immediately, deleting all its data.
//id: the instance ID to cancel
func (d dedicatedhostManager) CancelGuests(id int) ([]StatusInfo, error) {
	guests, err := d.VirtualDedicatedHost.Id(id).GetGuests()
	var listCaneledGuests []StatusInfo
	if err != nil {
		return listCaneledGuests, err
	}

	for _, guest := range guests {
		deletedGuest, err := d.DeleteGuest(*guest.Id)
		if err != nil {
			return listCaneledGuests, err
		}
		var statusGuest = StatusInfo{
			Id:     *guest.Id,
			Fqdn:   *guest.FullyQualifiedDomainName,
			Status: deletedGuest,
		}
		listCaneledGuests = append(listCaneledGuests, statusGuest)
	}
	return listCaneledGuests, nil
}

//Deletes a guest and returns 'Cancelled' or and Exception message
func (d dedicatedhostManager) DeleteGuest(Id int) (string, error) {
	status := "Cancelled"
	_, err := d.VirtualGuestService.Id(Id).DeleteObject()
	if err != nil {
		return "Failed", err
	}
	return status, nil
}

//Retrieve a list of all virtual servers on the dedicated host.
//integer identifier: The identifier of a dedicated host.
//integer cpus: filter based on number of CPUS.
//string domain: filter based on domain.
//string hostname: filter based on hostname.
//integer memory: filter based on amount of memory.
//list tags: filter based on list of tags.
func (d dedicatedhostManager) ListGuests(identifier int, cpu int, domain string, hostname string, memory int, tags []string, mask string) ([]datatypes.Virtual_Guest, error) {
	filters := filter.New()
	if cpu != 0 {
		filters = append(filters, filter.Path("guests.maxCpu").Eq(cpu))
	}
	if domain != "" {
		filters = append(filters, utils.QueryFilter(domain, "guests.domain"))
	}
	if hostname != "" {
		filters = append(filters, utils.QueryFilter(hostname, "guests.hostname"))
	}
	if memory != 0 {
		filters = append(filters, filter.Path("guests.maxMemory").Eq(memory))
	}
	if len(tags) > 0 {
		tagInterfaces := make([]interface{}, len(tags))
		for i, v := range tags {
			tagInterfaces[i] = v
		}
		filters = append(filters, filter.Path("guests.tagReferences.tag.name").In(tagInterfaces...))
	}

	guestList, err := d.VirtualDedicatedHost.Id(identifier).Mask(mask).Filter(filters.Build()).GetGuests()
	if err != nil {
		return []datatypes.Virtual_Guest{}, err
	}
	return guestList, nil
}

//Generate dedicated host payload.
func (d dedicatedhostManager) GenerateOrderTemplate(size, hostname, domain, datacenter string, billing string, routerId int) (datatypes.Container_Product_Order_Virtual_DedicatedHost, error) {
	mask := "items[keyName,capacity,description,attributes[id,attributeTypeKeyName],itemCategory[id,categoryCode],softwareDescription[id,referenceCode,longDescription],prices],activePresets,regions[location[location[priceGroups]]]"
	packages, err := d.PackageService.Mask(mask).Filter(filter.Path("keyName").Eq("DEDICATED_HOST").Build()).GetAllObjects()
	if err != nil {
		return datatypes.Container_Product_Order_Virtual_DedicatedHost{}, err
	}
	if len(packages) != 1 {
		return datatypes.Container_Product_Order_Virtual_DedicatedHost{}, errors.New(T("Ordering package is not found"))
	}
	hourly := billing == "hourly"
	location, err := GetLocation(packages[0], datacenter)
	if err != nil {
		return datatypes.Container_Product_Order_Virtual_DedicatedHost{}, err
	}
	priceId, err := GetDedicatedHostPriceId(packages[0].Items, size, hourly, location)
	if err != nil {
		return datatypes.Container_Product_Order_Virtual_DedicatedHost{}, err
	}
	order := datatypes.Container_Product_Order_Virtual_DedicatedHost{
		Container_Product_Order: datatypes.Container_Product_Order{
			Location: location.Keyname,
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(priceId)},
			},
			PackageId:        packages[0].Id,
			UseHourlyPricing: sl.Bool(hourly),
			Hardware: []datatypes.Hardware{
				datatypes.Hardware{
					Hostname: sl.String(hostname),
					Domain:   sl.String(domain),
					PrimaryBackendNetworkComponent: &datatypes.Network_Component{
						Router: &datatypes.Hardware{
							Id: sl.Int(routerId),
						},
					},
				},
			},
		},
	}
	return order, nil
}

//Verify the dedicated host order.
func (d dedicatedhostManager) VerifyInstanceCreation(orderTemplate datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order, error) {
	return d.OrderService.VerifyOrder(&orderTemplate)
}

//Order a dedicated host.
func (d dedicatedhostManager) OrderInstance(orderTemplate datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order_Receipt, error) {
	return d.OrderService.PlaceOrder(&orderTemplate, sl.Bool(false))
}
