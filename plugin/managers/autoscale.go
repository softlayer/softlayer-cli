package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	GetVirtualGuestMembers(id int, mask string) ([]datatypes.Scale_Member_Virtual_Guest, error)
}

type autoScaleManager struct {
	AutoScaleService services.Scale_Group
	AccountService   services.Account
}

func NewAutoScaleManager(session *session.Session) *autoScaleManager {
	return &autoScaleManager{
		services.GetScaleGroupService(session),
		services.GetAccountService(session),
	}
}

//Get virtual guest members about specific autoscale group
//id: Auto Sacale Group Id
//mask: object mask
func (as autoScaleManager) GetVirtualGuestMembers(id int, mask string) ([]datatypes.Scale_Member_Virtual_Guest, error) {
	return as.AutoScaleService.Id(id).Mask(mask).GetVirtualGuestMembers()
}
