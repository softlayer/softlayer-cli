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
	GetEvents(typeEvent string, mask string, dateFilter string) ([]datatypes.Notification_Occurrence_Event, error)
	GetEventDetail(identifier int, mask string) (datatypes.Notification_Occurrence_Event, error)
	GetInvoices(limit int, closed bool, getAll bool) ([]datatypes.Billing_Invoice, error)
}

type accountManager struct {
	AccountService services.Account
	Session        *session.Session
}

func NewAccountManager(session *session.Session) *accountManager {
	return &accountManager{
		AccountService: services.GetAccountService(session),
		Session:        session,
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
Gets all events with the potential to cause a service interruption with a specific keyName.
https://sldn.softlayer.com/reference/services/SoftLayer_Notification_Occurrence_Event/getAllObjects/
*/
func (a accountManager) GetEvents(typeEvent string, mask string, dateFilter string) ([]datatypes.Notification_Occurrence_Event, error) {
	NotificationOccurrenceEventService := services.GetNotificationOccurrenceEventService(a.Session)
	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("ASC"))
	filters = append(filters, filter.Path("notificationOccurrenceEventType.keyName").Eq(typeEvent))
	if dateFilter != "" {
		if typeEvent == "PLANNED"{
			filters = append(filters, filter.Path("endDate").DateAfter(dateFilter))
		}
		if typeEvent == "UNPLANNED_INCIDENT"{
			filters = append(filters, filter.Path("modifyDate").DateAfter(dateFilter))
		}
	}
	if typeEvent == "ANNOUNCEMENT"{
		filters = append(filters, filter.Path("statusCode.keyName").Eq("PUBLISHED"))
	}

	resourceList, err := NotificationOccurrenceEventService.Mask(mask).Filter(filters.Build()).GetAllObjects()
	if err != nil {
		return []datatypes.Notification_Occurrence_Event{}, err
	}
	return resourceList, err
}

/*
Gets a event with the potential to cause a service interruption.
https://sldn.softlayer.com/reference/services/SoftLayer_Notification_Occurrence_Event/getObject/
*/
func (a accountManager) GetEventDetail(identifier int, mask string) (datatypes.Notification_Occurrence_Event, error) {
	NotificationOccurrenceEventService := services.GetNotificationOccurrenceEventService(a.Session)
	
	resourceList, err := NotificationOccurrenceEventService.Mask(mask).Id(identifier).GetObject()
	if err != nil {
		return datatypes.Notification_Occurrence_Event{}, err
	}
	return resourceList, err
}

/*
Gets all invoices from the account
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getInvoices/
*/
func (a accountManager) GetInvoices(limit int, closed bool, getAll bool) ([]datatypes.Billing_Invoice, error) {
	mask := "mask[invoiceTotalAmount, itemCount]"
	filters := filter.New()
	filters = append(filters, filter.Path("invoices.id").OrderBy("DESC"))
	if !closed {
		filters = append(filters, filter.Path("invoices.statusCode").Eq("OPEN"))
	}
	resourceList := []datatypes.Billing_Invoice{}
	if getAll {
		i := 0
		for {
			resp, err := a.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetInvoices()
			i++
			if err != nil {
				return []datatypes.Billing_Invoice{}, err
			}
			resourceList = append(resourceList, resp...)
			if len(resp) < metadata.LIMIT {
				break
			}
		}
	} else {
		resp, err := a.AccountService.Mask(mask).Filter(filters.Build()).Limit(limit).GetInvoices()
		if err != nil {
			return []datatypes.Billing_Invoice{}, err
		}
		resourceList = append(resourceList, resp...)
	}

	return resourceList, nil
}
