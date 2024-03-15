package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

//counterfeiter:generate -o ../testhelpers/ . BandwidthManager
type BandwidthManager interface {
	GetLocationGroup() ([]datatypes.Location_Group, error)
	CreatePool(name string, regionId int) (datatypes.Network_Bandwidth_Version1_Allotment, error)
	DeletePool(bandwidthPoolID int) error
	EditPool(bandwidthPoolID int, regionId int, newPoolName string) (bool, error)
}

type bandwidthManager struct {
	BandwidthService     services.Network_Bandwidth_Version1_Allotment
	LocationGroupService services.Location_Group
	AccountService       services.Account
	Session              *session.Session
}

func NewBandwidthManager(session *session.Session) *bandwidthManager {
	return &bandwidthManager{
		BandwidthService:     services.GetNetworkBandwidthVersion1AllotmentService(session),
		LocationGroupService: services.GetLocationGroupService(session),
		AccountService:       services.GetAccountService(session),
		Session:              session,
	}
}

/*
https://sldn.softlayer.com/reference/services/SoftLayer_Location_Group/getAllObjects/
*/
func (a bandwidthManager) GetLocationGroup() ([]datatypes.Location_Group, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("locationGroupTypeId").Eq(1))
	return a.LocationGroupService.Filter(filters.Build()).GetAllObjects()
}

/*
Creates a Bandwidth Pool.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Bandwidth_Version1_Allotment/createObject/
*/
func (a bandwidthManager) CreatePool(name string, regionId int) (datatypes.Network_Bandwidth_Version1_Allotment, error) {
	currentAccount, err := a.AccountService.GetCurrentUser()
	if err != nil {
		return datatypes.Network_Bandwidth_Version1_Allotment{}, errors.NewAPIError(T("Failed to get currect user."), err.Error(), 2)
	}
	var template = datatypes.Network_Bandwidth_Version1_Allotment{
		AccountId:                sl.Int(*currentAccount.AccountId),
		BandwidthAllotmentTypeId: sl.Int(2),
		LocationGroupId:          sl.Int(regionId),
		Name:                     sl.String(name),
	}
	return a.BandwidthService.CreateObject(&template)
}

/*
Deletes a Bandwidth Pool.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Bandwidth_Version1_Allotment/requestVdrCancellation/
*/
func (a bandwidthManager) DeletePool(bandwidthPoolId int) error {

	_, err := a.BandwidthService.Id(bandwidthPoolId).RequestVdrCancellation()
	return err
}

/*
Edit a Bandwidth Pool.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Bandwidth_Version1_Allotment/requestVdrCancellation/
*/

func (a bandwidthManager) EditPool(bandwidthPoolId int, regionId int, newPoolName string) (bool, error) {
	mask := ""

	actualBandwidth, err := a.BandwidthService.Id(bandwidthPoolId).Mask(mask).GetObject()

	if err != nil {
		return false, errors.NewAPIError(T("Failed to get currect user."), err.Error(), 2)
	}

	var template = datatypes.Network_Bandwidth_Version1_Allotment{
		AccountId:                sl.Int(*actualBandwidth.AccountId),
		BandwidthAllotmentTypeId: sl.Int(2),
		LocationGroupId:          sl.Int(regionId),
		Name:                     sl.String(newPoolName),
	}
	return a.BandwidthService.Id(bandwidthPoolId).EditObject(&template)
}
