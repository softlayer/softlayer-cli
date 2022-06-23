package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type ObjectStorageManager interface {
	GetAccounts(mask string, limit int) ([]datatypes.Network_Storage, error)
}

type objectStorageManager struct {
	ObjectStorageService services.Account
	Session              *session.Session
}

func NewObjectStorageManager(session *session.Session) *objectStorageManager {
	return &objectStorageManager{
		ObjectStorageService: services.GetAccountService(session),
		Session:              session,
	}
}

/*
Gets an account’s associated Virtual Storage volumes.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getHubNetworkStorage/
*/
func (a objectStorageManager) GetAccounts(mask string, limit int) ([]datatypes.Network_Storage, error) {
	return a.ObjectStorageService.Mask(mask).Limit(limit).GetHubNetworkStorage()
}
