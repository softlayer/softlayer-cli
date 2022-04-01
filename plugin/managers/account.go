package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AccountManager interface {
	SummaryByDatacenter() (map[string]map[string]int, error)
	GetBandwidthPools() ([]datatypes.Network_Bandwidth_Version1_Allotment, error)
	GetBandwidthPoolServers(identifier int) (int, error)
	GetEventDetail(identifier int) (datatypes.Notification_Occurrence_Event, error)
}

type accountManager struct {
	AccountService                     services.Account
	NotificationOccurrenceEventService services.Notification_Occurrence_Event
	Session                            *session.Session
}

func NewAccountManager(session *session.Session) *accountManager {
	return &accountManager{
		AccountService:                     services.GetAccountService(session),
		NotificationOccurrenceEventService: services.GetNotificationOccurrenceEventService(session),
		Session:                            session,
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

func (a accountManager) GetEventDetail(identifier int) (datatypes.Notification_Occurrence_Event, error) {
	mask := "mask[acknowledgedFlag,attachments,impactedResources,statusCode,updates,notificationOccurrenceEventType]"
	resourceList, err := a.NotificationOccurrenceEventService.Mask(mask).Id(identifier).GetObject()
	if err != nil {
		return datatypes.Notification_Occurrence_Event{}, err
	}
	return resourceList, err
}
