package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

const (
	IMAGE_DEFAULT_MASK = "id, name, imageType, accountId"
	//IMAGE_DEFAULT_MASK = "id,accountId,name,globalIdentifier,blockDevices,parentId,createDate,transaction"
	IMAGE_DETAIL_MASK = `id, globalIdentifier, name, datacenter.name, status.name,
accountId, publicFlag, imageType, flexImageFlag, note, createDate, blockDevicesDiskSpaceTotal,
children[transaction, blockDevicesDiskSpaceTotal, datacenter.name]`
)

//Manages SoftLayer server images.
//See product information here: https://knowledgelayer.softlayer.com/topic/image-templates
type ImageManager interface {
	GetImage(imageId int) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	AddLocation(imageId int, locations []datatypes.Location) (bool, error)
	DeleteLocation(imageId int, locations []datatypes.Location) (bool, error)
	DeleteImage(imageId int) error
	ListPrivateImages(name string, mask string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	ListPublicImages(name string, mask string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	EditImage(imageId int, name string, note string, tag string) ([]bool, []string)
	ExportImage(imageId int, config datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (bool, error)
	ImportImage(config datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	GetDatacenters(imageId int) ([]datatypes.Location, error)
}

type imageManager struct {
	ImageService   services.Virtual_Guest_Block_Device_Template_Group
	AccountService services.Account
}

func NewImageManager(session *session.Session) *imageManager {
	return &imageManager{
		services.GetVirtualGuestBlockDeviceTemplateGroupService(session),
		services.GetAccountService(session),
	}
}

func (i imageManager) ExportImage(imageId int, config datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (bool, error) {
	return i.ImageService.Id(imageId).CopyToIcos(&config)
}
func (i imageManager) ImportImage(config datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	return i.ImageService.CreateFromIcos(&config)
}

//Get details about an image
//image: The ID of the image.
func (i imageManager) GetImage(imageId int) (datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	return i.ImageService.Id(imageId).Mask(IMAGE_DETAIL_MASK).GetObject()
}

//Delete the image by its ID
//imageId: The ID of the image.
func (i imageManager) DeleteImage(imageId int) error {
	_, err := i.ImageService.Id(imageId).DeleteObject()
	return err
}

//List all public images, fitler by its name
//name: filter based on name
func (i imageManager) ListPrivateImages(name string, mask string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	filters := filter.New()
	if name != "" {
		filters = append(filters, filter.Path("privateBlockDeviceTemplateGroups.name").Eq(name))
	}

	if mask == "" {
		mask = IMAGE_DEFAULT_MASK
	}
	//n := 0
	//var resourceList []datatypes.Virtual_Guest_Block_Device_Template_Group
	//for {
	//	resp, err := i.AccountService.Mask(IMAGE_DEFAULT_MASK).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(n * metadata.LIMIT).GetPrivateBlockDeviceTemplateGroups()
	//	n++
	//	if err != nil {
	//		return []datatypes.Virtual_Guest_Block_Device_Template_Group{}, err
	//	}
	//	resourceList = append(resourceList, resp...)
	//	if len(resp) < metadata.LIMIT {
	//		break
	//	}
	//}

	resourceList, err := i.AccountService.Mask(IMAGE_DEFAULT_MASK).Filter(filters.Build()).GetPrivateBlockDeviceTemplateGroups()
	if err != nil {
		return []datatypes.Virtual_Guest_Block_Device_Template_Group{}, err
	}
	return resourceList, nil
}

//List all public images,fitler by its name
//name: filter based on name
func (i imageManager) ListPublicImages(name string, mask string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	filters := filter.New()
	if name != "" {
		filters = append(filters, filter.Path("name").Eq(name))
	}
	if mask == "" {
		mask = IMAGE_DEFAULT_MASK
	}

	//n := 0
	//var resourceList []datatypes.Virtual_Guest_Block_Device_Template_Group
	//for {
	//	resp, err := i.ImageService.Mask(IMAGE_DEFAULT_MASK).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(n * metadata.LIMIT).GetPublicImages()
	//	n++
	//	if err != nil {
	//		return []datatypes.Virtual_Guest_Block_Device_Template_Group{}, err
	//	}
	//	resourceList = append(resourceList, resp...)
	//	if len(resp) < metadata.LIMIT {
	//		break
	//	}
	//}

	resourceList, err := i.ImageService.Mask(IMAGE_DEFAULT_MASK).Filter(filters.Build()).GetPublicImages()
	if err != nil {
		return []datatypes.Virtual_Guest_Block_Device_Template_Group{}, err
	}
	return resourceList, nil

}

//Edit image related details
//imageId: The ID of the image
//name: Name of the Image.
//note: Note of the image.
//tag: Tags of the image to be updated to.
func (i imageManager) EditImage(imageId int, name string, note string, tag string) ([]bool, []string) {
	var succeed []bool
	var messages []string
	image := datatypes.Virtual_Guest_Block_Device_Template_Group{}
	if name != "" {
		image.Name = sl.String(name)
	}
	if note != "" {
		image.Note = sl.String(note)
	}
	if name != "" || note != "" {
		_, err := i.ImageService.Id(imageId).EditObject(&image)
		if err != nil {
			succeed = append(succeed, false)
			messages = append(messages, err.Error()+"\n"+T("Failed to update the image {{.ID}}.", map[string]interface{}{"ID": imageId}))
		} else {
			if name != "" {
				succeed = append(succeed, true)
				messages = append(messages, T("The name of the image {{.ID}} is updated.", map[string]interface{}{"ID": imageId}))
			}
			if note != "" {
				succeed = append(succeed, true)
				messages = append(messages, T("The note of the image {{.ID}} is updated.", map[string]interface{}{"ID": imageId}))
			}
		}
	}
	if tag != "" {
		_, err := i.ImageService.Id(imageId).SetTags(sl.String(tag))
		if err != nil {
			succeed = append(succeed, false)
			messages = append(messages, err.Error()+"\n"+T("Failed to update the tag of the image {{.ID}}.", map[string]interface{}{"ID": imageId}))
		} else {
			succeed = append(succeed, true)
			messages = append(messages, T("The tag of the image {{.ID}} is updated.", map[string]interface{}{"ID": imageId}))
		}
	}

	return succeed, messages
}

//Add the location of the image
//imageId: The ID of the image
//location: location to remove of the image.
func (i imageManager) AddLocation(imageId int, locations []datatypes.Location) (bool, error) {
	return i.ImageService.Id(imageId).AddLocations(locations)
}

//Remove the location of the image
//imageId: The ID of the image
//location: location to remove of the image.
func (i imageManager) DeleteLocation(imageId int, locations []datatypes.Location) (bool, error) {
	return i.ImageService.Id(imageId).RemoveLocations(locations)
}

//Remove the location of the image
//imageId: The ID of the image
//location: location to remove of the image.
func (i imageManager) GetDatacenters(imageId int) ([]datatypes.Location, error) {
	return i.ImageService.Id(imageId).GetStorageLocations()
}
