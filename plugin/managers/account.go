package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type AccountManager interface {
	SummaryByDatacenter() (map[string]map[string]int, error)
	GetBandwidthPools() ([]datatypes.Network_Bandwidth_Version1_Allotment, error)
	GetBandwidthPoolServers(identifier int) (int, error)
	GetInvoiceDetail(identifier int) ([]datatypes.Billing_Invoice_Item, error)
}

type accountManager struct {
	AccountService        services.Account
	Session               *session.Session
}

func NewAccountManager(session *session.Session) *accountManager {
	return &accountManager{
		AccountService:        services.GetAccountService(session),
		Session:               session,
	}
}

//Summary of the networks on the account, grouped by data center.
//returns a map, the key of the map is datacenter name, the value of the map is another map
//the keys of the inner map are: vlan_count, public_ip_count, subnet_count, hardware_count, virtual_guest_count
//the value of the innter map are the count of those resources
func (a accountManager) SummaryByDatacenter() (map[string]map[string]int, error) {
	datacenters := make(map[string](map[string]int))
	vlans, err := a.AccountService.Mask(DEFAULT_VLAN_MASK).GetNetworkVlans()
	if err != nil {
		return datacenters, err
	}
	for _, vlan := range vlans {
		if vlan.PrimaryRouter != nil && vlan.PrimaryRouter.Datacenter != nil && vlan.PrimaryRouter.Datacenter.Name != nil {
			name := *vlan.PrimaryRouter.Datacenter.Name
			if datacenters[name] == nil {
				datacenters[name] = make(map[string]int)
			}
			datacenters[name]["vlan_count"]++
			if vlan.TotalPrimaryIpAddressCount != nil {
				datacenters[name]["public_ip_count"] += int(*vlan.TotalPrimaryIpAddressCount)
			}
			if vlan.SubnetCount != nil {
				datacenters[name]["subnet_count"] += int(*vlan.SubnetCount)
			}
			if vlan.HardwareCount != nil {
				datacenters[name]["hardware_count"] += int(*vlan.HardwareCount)
			}
			if vlan.VirtualGuestCount != nil {
				datacenters[name]["virtual_guest_count"] += int(*vlan.VirtualGuestCount)
			}
		}
	}
	return datacenters, nil
}

// https://sldn.softlayer.com/reference/services/SoftLayer_Account/getBandwidthAllotments/
func (a accountManager) GetBandwidthPools() ([]datatypes.Network_Bandwidth_Version1_Allotment, error) {
	mask := "mask[totalBandwidthAllocated,locationGroup, id, name, projectedPublicBandwidthUsage, " +
		"billingCyclePublicBandwidthUsage[amountOut,amountIn]]"
	pools, err := a.AccountService.Mask(mask).GetBandwidthAllotments()
	return pools, err
}

/*
Gets a count of all servers in a bandwidth pool
Getting the server counts individually is significantly faster than pulling them in
with the GetBandwidthPools api call.
*/
func (a accountManager) GetBandwidthPoolServers(identifier int) (int, error) {
	mask := "mask[id, bareMetalInstanceCount, hardwareCount, virtualGuestCount]"
	allotmentService := services.GetNetworkBandwidthVersion1AllotmentService(a.Session)
	counts, err := allotmentService.Mask(mask).Id(identifier).GetObject()
	total := 0
	if counts.BareMetalInstanceCount != nil {
		total += int(*counts.BareMetalInstanceCount)
	}
	if counts.HardwareCount != nil {
		total += int(*counts.HardwareCount)
	}
	if counts.VirtualGuestCount != nil {
		total += int(*counts.VirtualGuestCount)
	}
	return total, err
}

/*
Gets a list of top-level invoice items that are on the currently pending invoice.
https://sldn.softlayer.com/reference/services/SoftLayer_Billing_Invoice/getInvoiceTopLevelItems/
*/
func (a accountManager) GetInvoiceDetail(identifier int) ([]datatypes.Billing_Invoice_Item, error) {
	BillingInoviceService := services.GetBillingInvoiceService(a.Session)
	
	mask := "mask[id, description, hostName, domainName, oneTimeAfterTaxAmount, recurringAfterTaxAmount,createDate,categoryCode,category[name],location[name],children[id, category[name], description, oneTimeAfterTaxAmount, recurringAfterTaxAmount]]"

	filters := filter.New()
	filters = append(filters, filter.Path("hardware.id").OrderBy("DESC"))

	i := 0
	resourceList := []datatypes.Billing_Invoice_Item{}
	for {
		resp, err := BillingInoviceService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).Id(identifier).GetInvoiceTopLevelItems()
		i++
		if err != nil {
			return []datatypes.Billing_Invoice_Item{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}
