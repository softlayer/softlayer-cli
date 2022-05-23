package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	EditScaleGroup(id int, autoScaleTemplate *datatypes.Scale_Group) (bool, error)
	GetVirtualGuestMembers(id int, mask string) ([]datatypes.Scale_Member_Virtual_Guest, error)
	GetLogsScaleGroup(id int, mask string, dateFilter string) ([]datatypes.Scale_Group_Log, error)
	GetScaleGroup(id int, mask string) (datatypes.Scale_Group, error)
	ListScaleGroups(mask string) ([]datatypes.Scale_Group, error)
	Scale(id int, delta int) ([]datatypes.Scale_Member, error)
	ScaleTo(id int, delta int) ([]datatypes.Scale_Member, error)
	Delete(id int) (bool, error)
	CreateScaleGroup(autoScaleTemplate *datatypes.Scale_Group) (datatypes.Scale_Group, error)
	GetDatacenterByName(name string, typeName string) ([]datatypes.Location, error)
}

type autoScaleManager struct {
	AutoScaleService services.Scale_Group
	AccountService   services.Account
	LocationService  services.Location
}

func NewAutoScaleManager(session *session.Session) *autoScaleManager {
	return &autoScaleManager{
		services.GetScaleGroupService(session),
		services.GetAccountService(session),
		services.GetLocationService(session),
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

//Edit specific autoscale group
//id: Auto Sacale Group Id
//autoScaleTemplate: New Auto Scale Group data
func (as autoScaleManager) EditScaleGroup(id int, autoScaleTemplate *datatypes.Scale_Group) (bool, error) {
	return as.AutoScaleService.Id(id).EditObject(autoScaleTemplate)
}

//List all scale groups
//mask: object mask
func (as autoScaleManager) ListScaleGroups(mask string) ([]datatypes.Scale_Group, error) {
	if mask == "" {
		mask = "mask[id,cooldown,createDate,maximumMemberCount,minimumMemberCount,name,virtualGuestMemberTemplate,status,virtualGuestMembers]"
	}
	return as.AccountService.Mask(mask).GetScaleGroups()
}

//Scale this group up or down by the amount given.
//id: Auto Sacale Group Id
//delta: amount given
func (as autoScaleManager) Scale(id int, delta int) ([]datatypes.Scale_Member, error) {
	return as.AutoScaleService.Id(id).Scale(&delta)
}

//Scale this group up or down to the number given.
//id: Auto Sacale Group Id
//delta: amount given
func (as autoScaleManager) ScaleTo(id int, delta int) ([]datatypes.Scale_Member, error) {
	return as.AutoScaleService.Id(id).ScaleTo(&delta)
}

//Delete this group and members.
//id: Auto Sacale Group Id
func (as autoScaleManager) Delete(id int) (bool, error) {
	return as.AutoScaleService.Id(id).ForceDeleteObject()
}

//Create autoscale group
//autoScaleTemplate: New Auto Scale Group data
func (as autoScaleManager) CreateScaleGroup(autoScaleTemplate *datatypes.Scale_Group) (datatypes.Scale_Group, error) {
	return as.AutoScaleService.CreateObject(autoScaleTemplate)
}

//Get location using the name as filter
//name: location name
func (as autoScaleManager) GetDatacenterByName(name string, typeName string) ([]datatypes.Location, error) {
	objectfilter := filter.New()
	if typeName == "shortname" {
		objectfilter = append(objectfilter, filter.Path("name").Eq(name))
	}
	if typeName == "longname" {
		objectfilter = append(objectfilter, filter.Path("longName").Eq(name))
	}
	return as.LocationService.Filter(objectfilter.Build()).GetDatacenters()
}
