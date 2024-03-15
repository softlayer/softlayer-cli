package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type NasNetworkStorageManager interface {
	ListNasNetworkStorages(mask string) ([]datatypes.Network_Storage, error)
	GetNasNetworkStorage(netwokStorageId int, mask string) (datatypes.Network_Storage, error)
}

type nasNetworkStorageManager struct {
	NasNetworkStorageService services.Network_Storage
	AccountService           services.Account
}

func NewNasNetworkStorageManager(session *session.Session) *nasNetworkStorageManager {
	return &nasNetworkStorageManager{
		services.GetNetworkStorageService(session),
		services.GetAccountService(session),
	}
}

// List all NAS Network Storages
// mask: object mask
func (nas nasNetworkStorageManager) ListNasNetworkStorages(mask string) ([]datatypes.Network_Storage, error) {
	if mask == "" {
		mask = "mask[eventCount,serviceResource[datacenter.name]]"
	}
	return nas.AccountService.Mask(mask).GetNasNetworkStorage()
}

// Get NAS Network Storage object
// netwokStorageId: Id of Network Storage object
// mask: object mask
func (nas nasNetworkStorageManager) GetNasNetworkStorage(netwokStorageId int, mask string) (datatypes.Network_Storage, error) {
	if mask == "" {
		mask = "mask[accountId,capacityGb,createDate,id,nasType]"
	}
	return nas.NasNetworkStorageService.Id(netwokStorageId).Mask(mask).GetObject()
}
