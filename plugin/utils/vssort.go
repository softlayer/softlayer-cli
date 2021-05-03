package utils

import "github.com/softlayer/softlayer-go/datatypes"

type VirtualGuestById []datatypes.Virtual_Guest

func (a VirtualGuestById) Len() int {
	return len(a)
}
func (a VirtualGuestById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *(a[i].Id) < *(a[j].Id)
	}
	return false
}

type VirtualGuestByHostname []datatypes.Virtual_Guest

func (a VirtualGuestByHostname) Len() int {
	return len(a)
}
func (a VirtualGuestByHostname) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByHostname) Less(i, j int) bool {
	if a[i].Hostname != nil && a[j].Hostname != nil {
		return *(a[i].Hostname) < *(a[j].Hostname)
	}
	return false
}

type VirtualGuestByDomain []datatypes.Virtual_Guest

func (a VirtualGuestByDomain) Len() int {
	return len(a)
}
func (a VirtualGuestByDomain) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByDomain) Less(i, j int) bool {
	if a[i].Domain != nil && a[j].Domain != nil {
		return *(a[i].Domain) < *(a[j].Domain)
	}
	return false
}

type VirtualGuestByDatacenter []datatypes.Virtual_Guest

func (a VirtualGuestByDatacenter) Len() int {
	return len(a)
}
func (a VirtualGuestByDatacenter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByDatacenter) Less(i, j int) bool {
	if a[i].Datacenter.Name != nil && a[j].Datacenter.Name != nil {
		return *(a[i].Datacenter.Name) < *(a[j].Datacenter.Name)
	}
	return false
}

type VirtualGuestByCPU []datatypes.Virtual_Guest

func (a VirtualGuestByCPU) Len() int {
	return len(a)
}
func (a VirtualGuestByCPU) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByCPU) Less(i, j int) bool {
	if a[i].MaxCpu != nil && a[j].MaxCpu != nil {
		return *(a[i].MaxCpu) < *(a[j].MaxCpu)
	}
	return false

}

type VirtualGuestByMemory []datatypes.Virtual_Guest

func (a VirtualGuestByMemory) Len() int {
	return len(a)
}
func (a VirtualGuestByMemory) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByMemory) Less(i, j int) bool {
	if a[i].MaxMemory != nil && a[j].MaxMemory != nil {
		return *(a[i].MaxMemory) < *(a[j].MaxMemory)
	}
	return false
}

type VirtualGuestByPrimaryIp []datatypes.Virtual_Guest

func (a VirtualGuestByPrimaryIp) Len() int {
	return len(a)
}
func (a VirtualGuestByPrimaryIp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByPrimaryIp) Less(i, j int) bool {
	if a[i].PrimaryIpAddress != nil && a[j].PrimaryIpAddress != nil {
		return *(a[i].PrimaryIpAddress) < *(a[j].PrimaryIpAddress)
	}
	return false
}

type VirtualGuestByBackendIp []datatypes.Virtual_Guest

func (a VirtualGuestByBackendIp) Len() int {
	return len(a)
}
func (a VirtualGuestByBackendIp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VirtualGuestByBackendIp) Less(i, j int) bool {
	if a[i].PrimaryBackendIpAddress != nil && a[j].PrimaryBackendIpAddress != nil {
		return *(a[i].PrimaryBackendIpAddress) < *(a[j].PrimaryBackendIpAddress)
	}
	return false

}
