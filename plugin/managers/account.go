package managers

import (
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AccountManager interface {
	SummaryByDatacenter() (map[string]map[string]int, error)
}

type accountManager struct {
	AccountService services.Account
}

func NewAccountManager(session *session.Session) *accountManager {
	return &accountManager{
		services.GetAccountService(session),
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
