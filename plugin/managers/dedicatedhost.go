package managers

import (
	"errors"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	LOCATIONS      = "locations"
	DEDICATED_HOST = "dedicatedhost"
)

var existDatacenter = false

//Manages SoftLayer Dedicated host.
type DedicatedHostManager interface {
	ListGuests(identifier int, cpu int, domain string, hostname string, memory int, tags []string, mask string) ([]datatypes.Virtual_Guest, error)
	GetPackage() (datatypes.Product_Package, error)
	GetCreateOptions(productPackage datatypes.Product_Package) map[string]map[string]string
	GetVlansOptions(datacenter string, flavor string, productPackage datatypes.Product_Package) ([]datatypes.Network_Vlan, error)
}

type dedicatedhostManager struct {
	AccountService       services.Account
	VirtualDedicatedHost services.Virtual_DedicatedHost
	PackageService       services.Product_Package
}

func NewDedicatedhostManager(session *session.Session) *dedicatedhostManager {
	return &dedicatedhostManager{
		services.GetAccountService(session),
		services.GetVirtualDedicatedHostService(session),
		services.GetProductPackageService(session),
	}
}

//Get the package related to simple dedicatedhost ordering
func (d dedicatedhostManager) GetPackage() (datatypes.Product_Package, error) {
	mask := "items[id,description,prices,capacity,keyName,itemCategory[categoryCode],bundleItems[capacity,keyName,categories[categoryCode],hardwareGenericComponentModel[id,hardwareComponentType[keyName]]]],regions[location[location[priceGroups]]]"
	filters := filter.New()
	filters = append(filters, filter.Path("keyName").Eq("DEDICATED_HOST"))
	packages, err := d.PackageService.Mask(mask).Filter(filters.Build()).GetAllObjects()
	if err != nil {
		return datatypes.Product_Package{}, err
	}
	if len(packages) != 1 {
		return datatypes.Product_Package{}, errors.New(T("Ordering package is not found"))
	}
	return packages[0], nil
}

//Returns valid options for ordering hardware.
func (d dedicatedhostManager) GetCreateOptions(productPackage datatypes.Product_Package) map[string]map[string]string {
	//locations
	locations := make(map[string]string)
	for _, region := range productPackage.Regions {
		if region.Location != nil && region.Location.Location != nil && region.Location.Location.Name != nil && region.Location.Location.LongName != nil {
			locations[*region.Location.Location.Name] = *region.Location.Location.LongName
		}
	}
	//dedicatedhost
	dedicatedhost := make(map[string]string)
	for _, item := range productPackage.Items {
		if item.ItemCategory != nil && item.ItemCategory.CategoryCode != nil {
			if *item.ItemCategory.CategoryCode == "dedicated_virtual_hosts" {
				dedicatedhost[*item.KeyName] = *item.Description
			}
		}
	}

	return map[string]map[string]string{
		LOCATIONS:      locations,
		DEDICATED_HOST: dedicatedhost,
	}
}

//Get the private vlans in the account.
func (d dedicatedhostManager) GetVlansOptions(datacenter string, flavor string, productPackage datatypes.Product_Package) ([]datatypes.Network_Vlan, error) {
	maskVlans := "primaryRouter[datacenter]"
	maskItemPrices := "pricingLocationGroup[locations]"
	filters := filter.New()
	filters = append(filters, filter.Path("privateNetworkVlans.primaryRouter.datacenter.name").Eq(datacenter))
	dedicatedhostItems, err := d.PackageService.Id(*productPackage.Id).Mask(maskItemPrices).GetItemPrices()
	if err != nil {
		return []datatypes.Network_Vlan{}, err
	}

	for _, itemDedicatedHost := range dedicatedhostItems {
		if *itemDedicatedHost.Item.KeyName == flavor {
			if itemDedicatedHost.PricingLocationGroup != nil {
				for _, location := range itemDedicatedHost.PricingLocationGroup.Locations {
					if *location.Name == datacenter {
						existDatacenter = true
						break
					}
				}
			}
		}
	}

	if existDatacenter {
		return d.AccountService.Mask(maskVlans).Filter(filters.Build()).GetPrivateNetworkVlans()
	} else {
		return []datatypes.Network_Vlan{}, errors.New(T("There are not private vlans available for this datacenter."))
	}
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
