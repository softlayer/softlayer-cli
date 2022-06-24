package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
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
	if mask == "" {
		mask = "mask[id,username,notes,vendorName,serviceResource]"
	}

	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("ASC"))

	return a.ObjectStorageService.Mask(mask).Filter(filters.Build()).Limit(limit).GetHubNetworkStorage()
}
