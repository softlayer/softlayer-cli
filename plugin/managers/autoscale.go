package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	GetVirtualGuestMembers(id int, mask string) ([]datatypes.Scale_Member_Virtual_Guest, error)
	GetLogsScaleGroup(id int, mask string, dateFilter string) ([]datatypes.Scale_Group_Log, error)
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

//Get virtual guest members about specific autoscale group
//id: Auto Sacale Group Id
//mask: object mask
func (as autoScaleManager) GetVirtualGuestMembers(id int, mask string) ([]datatypes.Scale_Member_Virtual_Guest, error) {
	if mask == "" {
		mask = "mask[id, createDate, scaleGroup]"
	}
	return as.AutoScaleService.Id(id).Mask(mask).GetVirtualGuestMembers()
}

//Get logs about specific autoscale group
//id: Auto Sacale Group Id
//mask: object mask
//dateFilter: Earliest date to retrieve logs for [YYYY-MM-DD]
func (as autoScaleManager) GetLogsScaleGroup(id int, mask string, dateFilter string) ([]datatypes.Scale_Group_Log, error) {
	if mask == "" {
		mask = "mask[id,createDate,description,scaleGroup]"
	}
	if dateFilter != "" {
		filter := filter.New(filter.Path("logs.createDate").DateAfter(dateFilter))
		return as.AutoScaleService.Filter(filter.Build()).Id(id).Mask(mask).GetLogs()
	}
	return as.AutoScaleService.Id(id).Mask(mask).GetLogs()
}

//Get details about specific autoscale group
//id: Auto Sacale Group Id
//mask: object mask
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
//mask: object mask
func (as autoScaleManager) ListScaleGroups(mask string) ([]datatypes.Scale_Group, error) {
	if mask == "" {
		mask = "mask[id,cooldown,createDate,maximumMemberCount,minimumMemberCount,name,virtualGuestMemberTemplate,status,virtualGuestMembers]"
	}
	return as.AccountService.Mask(mask).GetScaleGroups()
}