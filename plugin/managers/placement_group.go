package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type PlaceGroupManager interface {
	List(mask string) ([]datatypes.Virtual_PlacementGroup, error)
	Create(templateObject *datatypes.Virtual_PlacementGroup) (datatypes.Virtual_PlacementGroup, error)
	GetRouter(datacenterId *int, mask string) ([]datatypes.Hardware, error)
	GetObject(placementgroupID int, mask string) (datatypes.Virtual_PlacementGroup, error)
	Delete(placementgroupID int) (bool, error)
	GetRules() ([]datatypes.Virtual_PlacementGroup_Rule, error)
	GetBackendRouterFromHostName() ([]datatypes.Network_Pod, error)
}

type placeGroupManager struct {
	Account        services.Account
	PlaceGroup     services.Virtual_PlacementGroup
	PlaceGroupRule services.Virtual_PlacementGroup_Rule
	NetworkPod     services.Network_Pod
}

func NewPlaceGroupManager(session *session.Session) *placeGroupManager {
	return &placeGroupManager{
		services.GetAccountService(session),
		services.GetVirtualPlacementGroupService(session),
		services.GetVirtualPlacementGroupRuleService(session),
		services.GetNetworkPodService(session),
	}
}

func (p placeGroupManager) List(mask string) ([]datatypes.Virtual_PlacementGroup, error) {
	if mask == "" {
		mask = "mask[id, name, createDate, rule, guestCount, backendRouter[id, hostname]]"
	}

	i := 0
	filters := filter.New()
	filters = append(filters, filter.Path("placementGroups.id").OrderBy("DESC"))
	resourceList := []datatypes.Virtual_PlacementGroup{}
	for {
		resp, err := p.Account.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetPlacementGroups()
		i++
		if err != nil {
			return []datatypes.Virtual_PlacementGroup{}, err
		}
		resourceList = append(resourceList, resp...)

		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil

	//return p.Account.Mask(mask).GetPlacementGroups()
}

func (i placeGroupManager) Create(templateObject *datatypes.Virtual_PlacementGroup) (datatypes.Virtual_PlacementGroup, error) {
	return i.PlaceGroup.CreateObject(templateObject)
}

func (p placeGroupManager) GetRouter(datacenterId *int, mask string) ([]datatypes.Hardware, error) {
	if mask == "" {
		mask = "mask[id, topLevelLocation[longName], hostname]"
	}

	return p.PlaceGroup.Mask(mask).GetAvailableRouters(datacenterId)

}

func (i placeGroupManager) GetObject(placementgroupID int, mask string) (datatypes.Virtual_PlacementGroup, error) {
	if mask == "" {
		mask = "mask[id, name, createDate, rule, backendRouter[id, hostname], guests[activeTransaction[id,transactionStatus[name,friendlyName]]]]"
	}
	return i.PlaceGroup.Mask(mask).Id(placementgroupID).GetObject()
}

func (i placeGroupManager) Delete(placementgroupID int) (bool, error) {
	return i.PlaceGroup.Id(placementgroupID).DeleteObject()
}

func (i placeGroupManager) GetRules() ([]datatypes.Virtual_PlacementGroup_Rule, error) {
	return i.PlaceGroupRule.GetAllObjects()
}

func (i placeGroupManager) GetBackendRouterFromHostName() ([]datatypes.Network_Pod, error) {
	return i.NetworkPod.GetAllObjects()
}
