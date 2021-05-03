package utils

import "github.com/softlayer/softlayer-go/datatypes"

type HardwareById []datatypes.Hardware_Server

func (a HardwareById) Len() int {
	return len(a)
}
func (a HardwareById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *(a[i].Id) < *(a[j].Id)
	}
	return false
}

type HardwareByGuid []datatypes.Hardware_Server

func (a HardwareByGuid) Len() int {
	return len(a)
}
func (a HardwareByGuid) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByGuid) Less(i, j int) bool {
	if a[i].GlobalIdentifier != nil && a[j].GlobalIdentifier != nil {
		return *(a[i].GlobalIdentifier) < *(a[j].GlobalIdentifier)
	}
	return false
}

type HardwareByHostname []datatypes.Hardware_Server

func (a HardwareByHostname) Len() int {
	return len(a)
}
func (a HardwareByHostname) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByHostname) Less(i, j int) bool {
	if a[i].Hostname != nil && a[j].Hostname != nil {
		return *(a[i].Hostname) < *(a[j].Hostname)
	}
	return false
}

type HardwareByDomain []datatypes.Hardware_Server

func (a HardwareByDomain) Len() int {
	return len(a)
}
func (a HardwareByDomain) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByDomain) Less(i, j int) bool {
	if a[i].Domain != nil && a[j].Domain != nil {
		return *(a[i].Domain) < *(a[j].Domain)
	}
	return false
}

type HardwareByCPU []datatypes.Hardware_Server

func (a HardwareByCPU) Len() int {
	return len(a)
}
func (a HardwareByCPU) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByCPU) Less(i, j int) bool {
	if a[i].ProcessorPhysicalCoreAmount != nil && a[j].ProcessorPhysicalCoreAmount != nil {
		return *(a[i].ProcessorPhysicalCoreAmount) < *(a[j].ProcessorPhysicalCoreAmount)
	}
	return false
}

type HardwareByMemory []datatypes.Hardware_Server

func (a HardwareByMemory) Len() int {
	return len(a)
}
func (a HardwareByMemory) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByMemory) Less(i, j int) bool {
	if a[i].MemoryCapacity != nil && a[j].MemoryCapacity != nil {
		return *(a[i].MemoryCapacity) < *(a[j].MemoryCapacity)
	}
	return false
}

type HardwareByPublicIP []datatypes.Hardware_Server

func (a HardwareByPublicIP) Len() int {
	return len(a)
}
func (a HardwareByPublicIP) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByPublicIP) Less(i, j int) bool {
	if a[i].PrimaryIpAddress != nil && a[j].PrimaryIpAddress != nil {
		return *(a[i].PrimaryIpAddress) < *(a[j].PrimaryIpAddress)
	}
	return false
}

type HardwareByPrivateIP []datatypes.Hardware_Server

func (a HardwareByPrivateIP) Len() int {
	return len(a)
}
func (a HardwareByPrivateIP) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByPrivateIP) Less(i, j int) bool {
	if a[i].PrimaryBackendIpAddress != nil && a[j].PrimaryBackendIpAddress != nil {
		return *(a[i].PrimaryBackendIpAddress) < *(a[j].PrimaryBackendIpAddress)
	}
	return false
}

type HardwareByRemoteIP []datatypes.Hardware_Server

func (a HardwareByRemoteIP) Len() int {
	return len(a)
}
func (a HardwareByRemoteIP) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByRemoteIP) Less(i, j int) bool {
	if a[i].NetworkManagementIpAddress != nil && a[j].NetworkManagementIpAddress != nil {
		return *(a[i].NetworkManagementIpAddress) < *(a[j].NetworkManagementIpAddress)
	}
	return false
}

type HardwareByStatus []datatypes.Hardware_Server

func (a HardwareByStatus) Len() int {
	return len(a)
}
func (a HardwareByStatus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByStatus) Less(i, j int) bool {
	if a[i].HardwareStatus != nil && a[i].HardwareStatus.Status != nil && a[j].HardwareStatus != nil && a[j].HardwareStatus.Status != nil {
		return *(a[i].HardwareStatus.Status) < *(a[j].HardwareStatus.Status)
	}
	return false
}

type HardwareByLocation []datatypes.Hardware_Server

func (a HardwareByLocation) Len() int {
	return len(a)
}
func (a HardwareByLocation) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByLocation) Less(i, j int) bool {
	if a[i].Datacenter != nil && a[i].Datacenter.Name != nil && a[j].Datacenter != nil && a[j].Datacenter.Name != nil {
		return *(a[i].Datacenter.Name) < *(a[j].Datacenter.Name)
	}
	return false
}

type HardwareByCreated []datatypes.Hardware_Server

func (a HardwareByCreated) Len() int {
	return len(a)
}
func (a HardwareByCreated) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByCreated) Less(i, j int) bool {
	if a[i].ProvisionDate != nil && a[j].ProvisionDate != nil {
		iTime := a[i].ProvisionDate.Time
		jTime := a[j].ProvisionDate.Time
		return iTime.Before(jTime)
	}
	return false
}

type HardwareByCreatedBy []datatypes.Hardware_Server

func (a HardwareByCreatedBy) Len() int {
	return len(a)
}
func (a HardwareByCreatedBy) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByCreatedBy) Less(i, j int) bool {
	if a[i].BillingItem != nil &&
		a[i].BillingItem.OrderItem != nil &&
		a[i].BillingItem.OrderItem.Order.UserRecord != nil &&
		a[i].BillingItem.OrderItem.Order.UserRecord.Username != nil &&
		a[j].BillingItem != nil &&
		a[j].BillingItem.OrderItem != nil &&
		a[j].BillingItem.OrderItem.Order.UserRecord != nil &&
		a[j].BillingItem.OrderItem.Order.UserRecord.Username != nil {
		return *a[i].BillingItem.OrderItem.Order.UserRecord.Username < *a[j].BillingItem.OrderItem.Order.UserRecord.Username
	}
	return false
}

type HardwareByOS []datatypes.Hardware_Server

func (a HardwareByOS) Len() int {
	return len(a)
}
func (a HardwareByOS) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HardwareByOS) Less(i, j int) bool {
	if a[i].OperatingSystem != nil &&
		a[i].OperatingSystem.SoftwareLicense != nil &&
		a[i].OperatingSystem.SoftwareLicense.SoftwareDescription != nil &&
		a[i].OperatingSystem.SoftwareLicense.SoftwareDescription.Name != nil &&
		a[j].OperatingSystem != nil &&
		a[j].OperatingSystem.SoftwareLicense != nil &&
		a[j].OperatingSystem.SoftwareLicense.SoftwareDescription != nil &&
		a[j].OperatingSystem.SoftwareLicense.SoftwareDescription.Name != nil {
		return *a[i].OperatingSystem.SoftwareLicense.SoftwareDescription.Name < *a[j].OperatingSystem.SoftwareLicense.SoftwareDescription.Name
	}
	return false
}

type RouterByHostname []datatypes.Hardware

func (a RouterByHostname) Len() int {
	return len(a)
}
func (a RouterByHostname) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RouterByHostname) Less(i, j int) bool {
	if a[i].TopLevelLocation != nil && a[j].TopLevelLocation != nil && a[i].TopLevelLocation.LongName != nil && a[j].TopLevelLocation.LongName != nil {
		return *(a[i].TopLevelLocation.LongName) < *(a[j].TopLevelLocation.LongName)
	}
	return false
}
