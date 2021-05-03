package utils

import "github.com/softlayer/softlayer-go/datatypes"

type GroupById []datatypes.Network_SecurityGroup

func (a GroupById) Len() int {
	return len(a)
}
func (a GroupById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a GroupById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type GroupByName []datatypes.Network_SecurityGroup

func (a GroupByName) Len() int {
	return len(a)
}
func (a GroupByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a GroupByName) Less(i, j int) bool {
	if a[i].Name != nil && a[j].Name != nil {
		return *a[i].Name < *a[j].Name
	}
	return *a[i].Id < *a[j].Id
}

type GroupByDescription []datatypes.Network_SecurityGroup

func (a GroupByDescription) Len() int {
	return len(a)
}
func (a GroupByDescription) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a GroupByDescription) Less(i, j int) bool {
	if a[i].Description != nil && a[j].Description != nil {
		return *a[i].Description < *a[j].Description
	}
	return false
}

type GroupByCreated []datatypes.Network_SecurityGroup

func (a GroupByCreated) Len() int {
	return len(a)
}
func (a GroupByCreated) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a GroupByCreated) Less(i, j int) bool {
	if a[i].CreateDate != nil && a[j].CreateDate != nil {
		return (*a[i].CreateDate).Time.Before((*a[j].CreateDate).Time)
	}
	return false
}

type RuleById []datatypes.Network_SecurityGroup_Rule

func (a RuleById) Len() int {
	return len(a)
}
func (a RuleById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type RuleByRemoteIp []datatypes.Network_SecurityGroup_Rule

func (a RuleByRemoteIp) Len() int {
	return len(a)
}
func (a RuleByRemoteIp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByRemoteIp) Less(i, j int) bool {
	if a[i].RemoteIp != nil && a[j].RemoteIp != nil {
		return *a[i].RemoteIp < *a[j].RemoteIp
	}
	return false
}

type RuleByRemoteGroupId []datatypes.Network_SecurityGroup_Rule

func (a RuleByRemoteGroupId) Len() int {
	return len(a)
}
func (a RuleByRemoteGroupId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByRemoteGroupId) Less(i, j int) bool {
	if a[i].RemoteGroupId != nil && a[j].RemoteGroupId != nil {
		return *a[i].RemoteGroupId < *a[j].RemoteGroupId
	}
	return false
}

type RuleByDirection []datatypes.Network_SecurityGroup_Rule

func (a RuleByDirection) Len() int {
	return len(a)
}
func (a RuleByDirection) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByDirection) Less(i, j int) bool {
	if a[i].Direction != nil && a[j].Direction != nil {
		return *a[i].Direction < *a[j].Direction
	}
	return false
}

type RuleByEtherType []datatypes.Network_SecurityGroup_Rule

func (a RuleByEtherType) Len() int {
	return len(a)
}
func (a RuleByEtherType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByEtherType) Less(i, j int) bool {
	if a[i].Ethertype != nil && a[j].Ethertype != nil {
		return *a[i].Ethertype < *a[j].Ethertype
	}
	return false
}

type RuleByMinPort []datatypes.Network_SecurityGroup_Rule

func (a RuleByMinPort) Len() int {
	return len(a)
}
func (a RuleByMinPort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByMinPort) Less(i, j int) bool {
	if a[i].PortRangeMin != nil && a[j].PortRangeMin != nil {
		return *a[i].PortRangeMin < *a[j].PortRangeMin
	}
	return false
}

type RuleByMaxPort []datatypes.Network_SecurityGroup_Rule

func (a RuleByMaxPort) Len() int {
	return len(a)
}
func (a RuleByMaxPort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByMaxPort) Less(i, j int) bool {
	if a[i].PortRangeMax != nil && a[j].PortRangeMax != nil {
		return *a[i].PortRangeMax < *a[j].PortRangeMax
	}
	return false
}

type RuleByProtocol []datatypes.Network_SecurityGroup_Rule

func (a RuleByProtocol) Len() int {
	return len(a)
}
func (a RuleByProtocol) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RuleByProtocol) Less(i, j int) bool {
	if a[i].Protocol != nil && a[j].Protocol != nil {
		return *a[i].Protocol < *a[j].Protocol
	}
	return false
}

type InterfaceByInterfaceId []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding

func (a InterfaceByInterfaceId) Len() int {
	return len(a)
}
func (a InterfaceByInterfaceId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a InterfaceByInterfaceId) Less(i, j int) bool {
	if a[i].NetworkComponent != nil && a[j].NetworkComponent != nil {
		return *a[i].NetworkComponent.Id < *a[j].NetworkComponent.Id
	}
	return false
}

type InterfaceByVSId []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding

func (a InterfaceByVSId) Len() int {
	return len(a)
}
func (a InterfaceByVSId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a InterfaceByVSId) Less(i, j int) bool {
	if a[i].NetworkComponent != nil && a[j].NetworkComponent != nil && a[i].NetworkComponent.Guest != nil && a[j].NetworkComponent.Guest != nil {
		return *a[i].NetworkComponent.Guest.Id < *a[j].NetworkComponent.Guest.Id
	}
	return false
}

type InterfaceByVSHost []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding

func (a InterfaceByVSHost) Len() int {
	return len(a)
}
func (a InterfaceByVSHost) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a InterfaceByVSHost) Less(i, j int) bool {
	if a[i].NetworkComponent != nil &&
		a[j].NetworkComponent != nil &&
		a[i].NetworkComponent.Guest != nil &&
		a[j].NetworkComponent.Guest != nil &&
		a[i].NetworkComponent.Guest.Hostname != nil &&
		a[j].NetworkComponent.Guest.Hostname != nil {
		return *a[i].NetworkComponent.Guest.Hostname < *a[j].NetworkComponent.Guest.Hostname
	}
	return false
}
