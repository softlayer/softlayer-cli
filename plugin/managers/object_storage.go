package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type ObjectStorageManager interface {
	GetAccounts(mask string, limit int) ([]datatypes.Network_Storage, error)
	GetEndpoints(HubNetworkStorageId int) ([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint, error)
	ListCredential(StorageId int, mask string)([]datatypes.Network_Storage_Credential, error)
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

/*
Returns a collection of endpoint URLs available to this IBM Cloud Object Storage account.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Storage_Hub_Cleversafe_Account/getEndpoints/
*/
func (a objectStorageManager) GetEndpoints(HubNetworkStorageId int) ([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint, error) {
	NetworkStorageHubCleversafeAccountService := services.GetNetworkStorageHubCleversafeAccountService(a.Session)

	return NetworkStorageHubCleversafeAccountService.Id(HubNetworkStorageId).GetEndpoints(nil)
}

/*
Gets credentials used for generating an AWS signature. Max of 2.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Storage_Hub_Cleversafe_Account/getCredentials/
*/
func (a objectStorageManager) ListCredential(StorageId int, mask string) ([]datatypes.Network_Storage_Credential, error) {
	NetworkStorageHubCleversafeAccountService := services.GetNetworkStorageHubCleversafeAccountService(a.Session)

	return NetworkStorageHubCleversafeAccountService.Mask(mask).Id(StorageId).GetCredentials()
}
