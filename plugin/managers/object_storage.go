package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ObjectStorageManager interface {
	GetAccounts(mask string) ([]datatypes.Network_Storage, error)
	GetEndpoints(HubNetworkStorageId int) ([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint, error)
	ListCredential(StorageId int, mask string) ([]datatypes.Network_Storage_Credential, error)
	CreateCredential(StorageId int, mask string) ([]datatypes.Network_Storage_Credential, error)
}

type objectStorageManager struct {
	ObjectStorageService services.Network_Storage_Hub_Cleversafe_Account
	Session              *session.Session
}

func NewObjectStorageManager(session *session.Session) *objectStorageManager {
	return &objectStorageManager{
		ObjectStorageService: services.GetNetworkStorageHubCleversafeAccountService(session),
		Session:              session,
	}
}

/*
Gets an account’s associated Virtual Storage volumes.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getHubNetworkStorage/
*/
func (a objectStorageManager) GetAccounts(mask string) ([]datatypes.Network_Storage, error) {
	AccountService := services.GetAccountService(a.Session)

	if mask == "" {
		mask = "mask[id,username,notes,vendorName,serviceResource]"
	}

	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("ASC"))

	i := 0
	resourceList := []datatypes.Network_Storage{}
	for {
		resp, err := AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetHubNetworkStorage()
		i++
		if err != nil {
			return []datatypes.Network_Storage{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Returns a collection of endpoint URLs available to this IBM Cloud Object Storage account.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Storage_Hub_Cleversafe_Account/getEndpoints/
*/
func (a objectStorageManager) GetEndpoints(HubNetworkStorageId int) ([]datatypes.Container_Network_Storage_Hub_ObjectStorage_Endpoint, error) {
	return a.ObjectStorageService.Id(HubNetworkStorageId).GetEndpoints(nil)
}

/*
Gets credentials used for generating an AWS signature. Max of 2.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Storage_Hub_Cleversafe_Account/getCredentials/
*/
func (a objectStorageManager) ListCredential(StorageId int, mask string) ([]datatypes.Network_Storage_Credential, error) {
	return a.ObjectStorageService.Mask(mask).Id(StorageId).GetCredentials()
}

/*
Create credentials for an IBM Cloud Object Storage Account.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Storage_Hub_Cleversafe_Account/credentialCreate/
*/
func (a objectStorageManager) CreateCredential(StorageId int, mask string) ([]datatypes.Network_Storage_Credential, error) {
	return a.ObjectStorageService.Mask(mask).Id(StorageId).CredentialCreate()
}
