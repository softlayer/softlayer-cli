package utils

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
)

type VolumeById []datatypes.Network_Storage

func (a VolumeById) Len() int {
	return len(a)
}
func (a VolumeById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type VolumeByUsername []datatypes.Network_Storage

func (a VolumeByUsername) Len() int {
	return len(a)
}
func (a VolumeByUsername) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByUsername) Less(i, j int) bool {
	if a[i].Username != nil && a[j].Username != nil {
		return *a[i].Username < *a[j].Username
	}
	return false
}

type VolumeByDatacenter []datatypes.Network_Storage

func (a VolumeByDatacenter) Len() int {
	return len(a)
}
func (a VolumeByDatacenter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByDatacenter) Less(i, j int) bool {
	if a[i].ServiceResource != nil &&
		a[i].ServiceResource.Datacenter != nil &&
		a[i].ServiceResource.Datacenter.Name != nil &&
		a[j].ServiceResource != nil &&
		a[j].ServiceResource.Datacenter != nil &&
		a[j].ServiceResource.Datacenter.Name != nil {
		return *a[i].ServiceResource.Datacenter.Name < *a[j].ServiceResource.Datacenter.Name
	}
	return true
}

type VolumeByStorageType []datatypes.Network_Storage

func (a VolumeByStorageType) Len() int {
	return len(a)
}
func (a VolumeByStorageType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByStorageType) Less(i, j int) bool {
	if a[i].StorageType != nil &&
		a[i].StorageType.KeyName != nil &&
		a[j].StorageType != nil &&
		a[j].StorageType.KeyName != nil {
		return *a[i].StorageType.KeyName < *a[j].StorageType.KeyName
	}
	return true
}

type VolumeByCapacity []datatypes.Network_Storage

func (a VolumeByCapacity) Len() int {
	return len(a)
}
func (a VolumeByCapacity) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByCapacity) Less(i, j int) bool {
	if a[i].CapacityGb != nil && a[j].CapacityGb != nil {
		return *a[i].CapacityGb < *a[j].CapacityGb
	}
	return false
}

type VolumeByBytesUsed []datatypes.Network_Storage

func (a VolumeByBytesUsed) Len() int {
	return len(a)
}
func (a VolumeByBytesUsed) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByBytesUsed) Less(i, j int) bool {
	if a[i].BytesUsed != nil && a[j].BytesUsed != nil {
		return *a[i].BytesUsed > *a[j].BytesUsed
	}
	return false
}

type VolumeByIPAddress []datatypes.Network_Storage

func (a VolumeByIPAddress) Len() int {
	return len(a)
}
func (a VolumeByIPAddress) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByIPAddress) Less(i, j int) bool {
	if a[i].ServiceResourceBackendIpAddress != nil && a[j].ServiceResourceBackendIpAddress != nil {
		return *a[i].ServiceResourceBackendIpAddress < *a[j].ServiceResourceBackendIpAddress
	}
	return false
}

type VolumeByLunId []datatypes.Network_Storage

func (a VolumeByLunId) Len() int {
	return len(a)
}
func (a VolumeByLunId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByLunId) Less(i, j int) bool {
	if a[i].LunId != nil && a[j].LunId != nil {
		aLunId, err1 := strconv.Atoi(*a[i].LunId)
		bLunId, err2 := strconv.Atoi(*a[j].LunId)
		if err1 == nil && err2 == nil {
			return aLunId < bLunId
		}
		return *a[i].LunId < *a[j].LunId
	}
	return false
}

type VolumeByTxnCount []datatypes.Network_Storage

func (a VolumeByTxnCount) Len() int {
	return len(a)
}
func (a VolumeByTxnCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByTxnCount) Less(i, j int) bool {
	if a[i].ActiveTransactionCount != nil && a[j].ActiveTransactionCount != nil {
		return int(*a[i].ActiveTransactionCount) < int(*a[j].ActiveTransactionCount)
	}
	return false
}

type VolumeByCreatedBy []datatypes.Network_Storage

func (a VolumeByCreatedBy) Len() int {
	return len(a)
}
func (a VolumeByCreatedBy) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByCreatedBy) Less(i, j int) bool {
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
	return true
}

type VolumeByMountAddr []datatypes.Network_Storage

func (a VolumeByMountAddr) Len() int {
	return len(a)
}
func (a VolumeByMountAddr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a VolumeByMountAddr) Less(i, j int) bool {
	if a[i].FileNetworkMountAddress != nil && a[j].FileNetworkMountAddress != nil {
		return *a[i].FileNetworkMountAddress < *a[j].FileNetworkMountAddress
	}
	return true
}

type SnapshotsById []datatypes.Network_Storage

func (a SnapshotsById) Len() int {
	return len(a)
}
func (a SnapshotsById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SnapshotsById) Less(i, j int) bool {
	return *a[i].Id < *a[j].Id
}

type SnapshotsByName []datatypes.Network_Storage

func (a SnapshotsByName) Len() int {
	return len(a)
}
func (a SnapshotsByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SnapshotsByName) Less(i, j int) bool {
	if a[i].Username != nil && a[j].Username != nil {
		return *a[i].Username < *a[j].Username
	}
	return true
}

type SnapshotsByCreated []datatypes.Network_Storage

func (a SnapshotsByCreated) Len() int {
	return len(a)
}
func (a SnapshotsByCreated) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SnapshotsByCreated) Less(i, j int) bool {
	if a[i].SnapshotCreationTimestamp != nil && a[j].SnapshotCreationTimestamp != nil {
		return *a[i].SnapshotCreationTimestamp < *a[j].SnapshotCreationTimestamp
	}
	return false
}

type SnapshotsBySize []datatypes.Network_Storage

func (a SnapshotsBySize) Len() int {
	return len(a)
}
func (a SnapshotsBySize) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SnapshotsBySize) Less(i, j int) bool {
	if a[i].SnapshotSizeBytes != nil && a[j].SnapshotSizeBytes != nil {
		ai, err1 := strconv.Atoi(*a[i].SnapshotSizeBytes)
		aj, err2 := strconv.Atoi(*a[j].SnapshotSizeBytes)
		if err1 == nil && err2 == nil {
			return ai < aj
		}
		return *a[i].SnapshotSizeBytes < *a[j].SnapshotSizeBytes
	}
	return false
}
