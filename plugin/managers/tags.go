package managers

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

//counterfeiter:generate -o ../testhelpers/ . TagsManager
type TagsManager interface {
	GetTagByTagName(tagName string) ([]datatypes.Tag, error)
	ListTags() ([]datatypes.Tag, error)
	ListEmptyTags() ([]datatypes.Tag, error)
	GetTagReferences(tagId int) ([]datatypes.Tag_Reference, error)
	ReferenceLookup(resourceType string, resourceId int) string
	DeleteTag(tagName string) (bool, error)
	SetTags(tags string, keyName string, resourceId int) (bool, error)
	GetUnattachedTags(mask string) ([]datatypes.Tag, error)
}

type tagsManager struct {
	TagService services.Tag
	Session    *session.Session
}

func NewTagsManager(session *session.Session) *tagsManager {
	return &tagsManager{
		services.GetTagService(session),
		session,
	}
}

// Calls SoftLayer_Tag::getTagByTagName(tagName)
func (tag tagsManager) GetTagByTagName(tagName string) ([]datatypes.Tag, error) {

	tagDetails, err := tag.TagService.GetTagByTagName(&tagName)
	return tagDetails, err
}

func (tag tagsManager) ListTags() ([]datatypes.Tag, error) {

	tags := []datatypes.Tag{}
	objectMask := "mask[id,name,referenceCount]"
	i := 0
	for {
		resp, err := tag.TagService.Mask(objectMask).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).
			GetAttachedTagsForCurrentUser()
		i++
		if err != nil {
			return tags, err
		}
		tags = append(tags, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}

	return tags, nil
}

func (tag tagsManager) ListEmptyTags() ([]datatypes.Tag, error) {

	tags := []datatypes.Tag{}
	objectMask := "mask[id,name,referenceCount]"
	i := 0
	for {
		resp, err := tag.TagService.Mask(objectMask).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).
			GetUnattachedTagsForCurrentUser()
		i++
		if err != nil {
			return tags, err
		}
		tags = append(tags, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return tags, nil
}

func (tag tagsManager) GetTagReferences(tagId int) ([]datatypes.Tag_Reference, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("references.id").OrderBy("ASC"))

	objectMask := "mask[tagType]"
	tagReferences := []datatypes.Tag_Reference{}
	i := 0
	for {
		references, err := tag.TagService.Mask(objectMask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).
			Id(tagId).GetReferences()
		i++
		if err != nil {
			return references, err
		}
		tagReferences = append(tagReferences, references...)
		if len(references) < metadata.LIMIT {
			break
		}
	}

	return tagReferences, nil
}

func (tag tagsManager) DeleteTag(tagName string) (bool, error) {
	return tag.TagService.DeleteTag(&tagName)
}

// Sadly there isn't a great way to go from a Tag to an actual resource.
// So we need to use the tagType->keyName to figure out which service we need to use.
// From  SoftLayer_Tag::getAllTagTypes()
// |Type                             |Service |
// | -----------------------------   | ------ |
// |Hardware                         |HARDWARE|
// |CCI                              |GUEST|
// |Account Document                 |ACCOUNT_DOCUMENT|
// |Ticket                           |TICKET|
// |Vlan Firewall                    |NETWORK_VLAN_FIREWALL|
// |Contract                         |CONTRACT|
// |Image Template                   |IMAGE_TEMPLATE|
// |Application Delivery Controller  |APPLICATION_DELIVERY_CONTROLLER|
// |Vlan                             |NETWORK_VLAN|
// |Dedicated Host                   |DEDICATED_HOST|
func (tag tagsManager) ReferenceLookup(resourceType string, resourceId int) string {

	name := fmt.Sprintf("Unable to lookup %v", resourceType)
	switch resourceType {
	case "HARDWARE":
		service := services.GetHardwareService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.FullyQualifiedDomainName, err)
	case "GUEST":
		service := services.GetVirtualGuestService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.FullyQualifiedDomainName, err)
	case "TICKET":
		service := services.GetTicketService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.Title, err)
	case "NETWORK_VLAN_FIREWALL":
		service := services.GetNetworkVlanFirewallService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.PrimaryIpAddress, err)
	case "IMAGE_TEMPLATE":
		//SoftLayer_Virtual_Guest_Block_Device_Template_Group
		service := services.GetVirtualGuestBlockDeviceTemplateGroupService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.Name, err)
	case "APPLICATION_DELIVERY_CONTROLLER":
		//SoftLayer_Network_Application_Delivery_Controller
		service := services.GetNetworkApplicationDeliveryControllerService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.Name, err)
	case "NETWORK_VLAN":
		service := services.GetNetworkVlanService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.Name, err)
	case "NETWORK_SUBNET":
		service := services.GetNetworkSubnetService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.NetworkIdentifier, err)
	case "DEDICATED_HOST":
		service := services.GetVirtualDedicatedHostService(tag.Session)
		resource, err := service.Id(resourceId).GetObject()
		name = NameCheck(resource.Name, err)
	}

	return name

}

func NameCheck(name *string, err error) string {

	checked_name := ""
	if err != nil {
		apiError := err.(sl.Error)
		if apiError.StatusCode == 404 {
			checked_name = "Not Found"
		} else {
			checked_name = fmt.Sprintf("%v", err)
		}
	} else {
		checked_name = utils.FormatStringPointer(name)
	}
	return checked_name
}

func (tag tagsManager) SetTags(tags string, keyName string, resourceId int) (bool, error) {
	return tag.TagService.SetTags(&tags, &keyName, &resourceId)
}

// Returns unattached tags for current user
// mask: Object mask
func (tag tagsManager) GetUnattachedTags(mask string) ([]datatypes.Tag, error) {
	return tag.TagService.Mask(mask).GetUnattachedTagsForCurrentUser()
}
