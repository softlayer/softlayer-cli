package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	GetScaleGroup(id int, mask string) (datatypes.Scale_Group, error)
	EditScaleGroup(id int, autoScaleTemplate *datatypes.Scale_Group) (bool, error)
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

//Get details about specific autoscale group
func (as autoScaleManager) GetScaleGroup(id int, mask string) (datatypes.Scale_Group, error) {
	if mask == "" {
		mask = `mask[virtualGuestMembers[id,virtualGuest[id,hostname,domain,provisionDate]],terminationPolicy,
		virtualGuestMemberCount,virtualGuestMemberTemplate[sshKeys],
		policies[id,name,createDate,cooldown,actions,triggers,scaleActions],
		networkVlans[networkVlanId,networkVlan[networkSpace,primaryRouter[hostname]]],
		loadBalancers,regionalGroup[locations]]`
	}
	return as.AutoScaleService.Id(id).Mask(mask).GetObject()
}

//Edit specific autoscale group
//id: Auto Sacale Group Id
//autoScaleTemplate: New Auto Scale Group data
func (as autoScaleManager) EditScaleGroup(id int, autoScaleTemplate *datatypes.Scale_Group) (bool, error) {
	return as.AutoScaleService.Id(id).EditObject(autoScaleTemplate)
}
