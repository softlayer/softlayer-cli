package utils

import "github.com/softlayer/softlayer-go/datatypes"

type VlanById []datatypes.Network_Vlan

func (a VlanById) Len() int {
	return len(a)
}
func (a VlanById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *(a[i].Id) < *(a[j].Id)
	}
	return false
}

type VlanByNumber []datatypes.Network_Vlan

func (a VlanByNumber) Len() int {
	return len(a)
}
func (a VlanByNumber) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByNumber) Less(i, j int) bool {
	if a[i].VlanNumber != nil && a[j].VlanNumber != nil {
		return *(a[i].VlanNumber) < *(a[j].VlanNumber)
	}
	return false
}

type VlanByName []datatypes.Network_Vlan

func (a VlanByName) Len() int {
	return len(a)
}
func (a VlanByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByName) Less(i, j int) bool {
	if a[i].Name != nil && a[j].Name != nil {
		return *(a[i].Name) < *(a[j].Name)
	}
	return false

}

type VlanByFirewall []datatypes.Network_Vlan

func (a VlanByFirewall) Len() int {
	return len(a)
}
func (a VlanByFirewall) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByFirewall) Less(i, j int) bool {
	return len(a[i].FirewallInterfaces) < len(a[j].FirewallInterfaces)
}

type VlanByDatacenter []datatypes.Network_Vlan

func (a VlanByDatacenter) Len() int {
	return len(a)
}
func (a VlanByDatacenter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByDatacenter) Less(i, j int) bool {
	if a[i].PrimaryRouter != nil &&
		a[i].PrimaryRouter.Datacenter != nil &&
		a[i].PrimaryRouter.Datacenter.Name != nil &&
		a[j].PrimaryRouter != nil &&
		a[j].PrimaryRouter.Datacenter != nil &&
		a[j].PrimaryRouter.Datacenter.Name != nil {
		return *(a[i].PrimaryRouter.Datacenter.Name) < *(a[j].PrimaryRouter.Datacenter.Name)
	}
	return false
}

type VlanByHardwareCount []datatypes.Network_Vlan

func (a VlanByHardwareCount) Len() int {
	return len(a)
}
func (a VlanByHardwareCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByHardwareCount) Less(i, j int) bool {
	if a[i].HardwareCount != nil && a[j].HardwareCount != nil {
		return *(a[i].HardwareCount) < *(a[j].HardwareCount)
	}
	return false
}

type VlanByVirtualServerCount []datatypes.Network_Vlan

func (a VlanByVirtualServerCount) Len() int {
	return len(a)
}
func (a VlanByVirtualServerCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByVirtualServerCount) Less(i, j int) bool {
	if a[i].VirtualGuestCount != nil && a[j].VirtualGuestCount != nil {
		return *(a[i].VirtualGuestCount) < *(a[j].VirtualGuestCount)
	}
	return false
}

type VlanByPublicIPCount []datatypes.Network_Vlan

func (a VlanByPublicIPCount) Len() int {
	return len(a)
}
func (a VlanByPublicIPCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VlanByPublicIPCount) Less(i, j int) bool {
	if a[i].TotalPrimaryIpAddressCount != nil && a[j].TotalPrimaryIpAddressCount != nil {
		return *(a[i].TotalPrimaryIpAddressCount) < *(a[j].TotalPrimaryIpAddressCount)
	}
	return false
}
