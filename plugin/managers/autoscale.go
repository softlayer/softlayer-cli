package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

//Manages SoftLayer server images.
//See product information here: https://knowledgelayer.softlayer.com/topic/image-templates
type AutoScaleManager interface {
	GetScaleGroup(id int, mask string) (datatypes.Scale_Group, error)
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

//List all scale groups
func (as autoScaleManager) ListScaleGroups(mask string) ([]datatypes.Scale_Group, error) {
	if mask == "" {
		mask = "mask[id,cooldown,createDate,maximumMemberCount,minimumMemberCount,name,virtualGuestMemberTemplate,status,virtualGuestMembers]"
	}
	return as.AccountService.Mask(mask).GetScaleGroups()
}
