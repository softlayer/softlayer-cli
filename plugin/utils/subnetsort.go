package utils

import "github.com/softlayer/softlayer-go/datatypes"

type SubnetById []datatypes.Network_Subnet

func (a SubnetById) Len() int {
	return len(a)
}
func (a SubnetById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *(a[i].Id) < *(a[j].Id)
	}
	return false
}

type SubnetByIdentifier []datatypes.Network_Subnet

func (a SubnetByIdentifier) Len() int {
	return len(a)
}
func (a SubnetByIdentifier) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByIdentifier) Less(i, j int) bool {
	if a[i].NetworkIdentifier != nil && a[j].NetworkIdentifier != nil {
		return *(a[i].NetworkIdentifier) < *(a[j].NetworkIdentifier)
	}
	return false
}

type SubnetByType []datatypes.Network_Subnet

func (a SubnetByType) Len() int {
	return len(a)
}
func (a SubnetByType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByType) Less(i, j int) bool {
	if a[i].SubnetType != nil && a[j].SubnetType != nil {
		return *(a[i].SubnetType) < *(a[j].SubnetType)
	}
	return false
}

type SubnetByNetworkSpace []datatypes.Network_Subnet

func (a SubnetByNetworkSpace) Len() int {
	return len(a)
}
func (a SubnetByNetworkSpace) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByNetworkSpace) Less(i, j int) bool {
	if a[i].NetworkVlan != nil && a[i].NetworkVlan.NetworkSpace != nil &&
		a[j].NetworkVlan != nil && a[j].NetworkVlan.NetworkSpace != nil {
		return *(a[i].NetworkVlan.NetworkSpace) < *(a[j].NetworkVlan.NetworkSpace)
	}
	return false
}

type SubnetByDatacenter []datatypes.Network_Subnet

func (a SubnetByDatacenter) Len() int {
	return len(a)
}
func (a SubnetByDatacenter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByDatacenter) Less(i, j int) bool {
	if a[i].Datacenter != nil && a[i].Datacenter.Name != nil &&
		a[j].Datacenter != nil && a[j].Datacenter.Name != nil {
		return *(a[i].Datacenter.Name) < *(a[j].Datacenter.Name)
	}
	return false
}

type SubnetByVlanId []datatypes.Network_Subnet

func (a SubnetByVlanId) Len() int {
	return len(a)
}
func (a SubnetByVlanId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByVlanId) Less(i, j int) bool {
	if a[i].NetworkVlan != nil && a[i].NetworkVlan.Id != nil &&
		a[j].NetworkVlan != nil && a[j].NetworkVlan.Id != nil {
		return *(a[i].NetworkVlan.Id) < *(a[j].NetworkVlan.Id)
	}
	return false
}

type SubnetByIpCount []datatypes.Network_Subnet

func (a SubnetByIpCount) Len() int {
	return len(a)
}
func (a SubnetByIpCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByIpCount) Less(i, j int) bool {
	if a[i].IpAddresses != nil && a[j].IpAddresses != nil {
		return len(a[i].IpAddresses) < len(a[j].IpAddresses)
	}
	return false
}

type SubnetByHardwareCount []datatypes.Network_Subnet

func (a SubnetByHardwareCount) Len() int {
	return len(a)
}
func (a SubnetByHardwareCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByHardwareCount) Less(i, j int) bool {
	if a[i].Hardware != nil && a[j].Hardware != nil {
		return len(a[i].Hardware) < len(a[j].Hardware)
	}
	return false
}

type SubnetByVSCount []datatypes.Network_Subnet

func (a SubnetByVSCount) Len() int {
	return len(a)
}
func (a SubnetByVSCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SubnetByVSCount) Less(i, j int) bool {
	if a[i].VirtualGuests != nil && a[j].VirtualGuests != nil {
		return len(a[i].VirtualGuests) < len(a[j].VirtualGuests)
	}
	return false
}
