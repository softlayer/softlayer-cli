package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

//Manages SoftLayer server images.
//See product information here: https://knowledgelayer.softlayer.com/topic/image-templates
type AutoScaleManager interface {
	ListScaleGroups(mask string) ([]datatypes.Scale_Group, error)
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

//List all scale groups
func (as autoScaleManager) ListScaleGroups(mask string) ([]datatypes.Scale_Group, error) {
	return as.AccountService.Mask(mask).GetScaleGroups()
}
