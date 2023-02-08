package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	INSTANCE_DEFAULT_MASK = "id, globalIdentifier, hostname, hourlyBillingFlag, domain, fullyQualifiedDomainName, status.name, " +
		"powerState.name, activeTransaction, datacenter.name, account.id, " +
		"maxCpu, maxMemory, primaryIpAddress, primaryBackendIpAddress, " +
		"privateNetworkOnlyFlag, dedicatedAccountHostOnlyFlag, createDate, modifyDate, " +
		"billingItem[orderItem[id,order.userRecord[username]], recurringFee,nextInvoiceChildren], notes, tagReferences.tag.name"
	INSTANCE_DETAIL_MASK = "id,globalIdentifier,fullyQualifiedDomainName,hostname,domain,createDate,modifyDate,provisionDate,notes," +
		"dedicatedAccountHostOnlyFlag,privateNetworkOnlyFlag,primaryBackendIpAddress,primaryIpAddress," +
		"networkComponents[id,status,speed,maxSpeed,name,macAddress,primaryIpAddress,port,primarySubnet,securityGroupBindings[securityGroup[id, name]]]," +
		"lastKnownPowerState.name,powerState,status,maxCpu,maxMemory,datacenter,activeTransaction[id, transactionStatus[friendlyName,name]]," +
		"lastOperatingSystemReload.id,blockDevices,blockDeviceTemplateGroup[id, name, globalIdentifier],postInstallScriptUri," +
		"operatingSystem[passwords[username,password],softwareLicense.softwareDescription[manufacturer,name,version,referenceCode]]," +
		"softwareComponents[passwords[username,password,notes],softwareLicense[softwareDescription[manufacturer,name,version,referenceCode]]]," +
		"hourlyBillingFlag,userData," +
		"billingItem[id,package,nextInvoiceTotalRecurringAmount,nextInvoiceChildren[description,categoryCode,recurringFee,nextInvoiceTotalRecurringAmount],children[description,categoryCode,nextInvoiceTotalRecurringAmount],orderItem[id,order.userRecord[username],preset.keyName]]," +
		"tagReferences[id,tag[name,id]],networkVlans[id,vlanNumber,networkSpace],dedicatedHost.id,transientGuestFlag,lastTransaction[transactionGroup]"
	HOST_DEFAULT_MASK = "id,name,createDate,cpuCount,diskCapacity,memoryCapacity,guestCount,datacenter,backendRouter,allocationStatus"

	KEY_DATABASE      = "databases"
	KEY_GUEST         = "guests"
	HOST_DEFAULT_SIZE = "56_CORES_X_242_RAM_X_1_4_TB"
)

var (
	FlavorKeys   = []string{"B1", "BL1", "BL2", "C1", "M1"}
	FlavorLabels = map[string]string{
		"B1":  T("balanced"),
		"BL1": T("balanced local - hdd"),
		"BL2": T("balanced local - ssd"),
		"C1":  T("compute"),
		"M1":  T("memory"),
	}
)

//Manages SoftLayer Virtual Servers.
//See product information here: http://www.softlayer.com/virtual-servers
type VirtualServerManager interface {
	AttachPortableStorage(id int, portableStorageId int) (datatypes.Provisioning_Version1_Transaction, error)
	AuthorizeStorage(id int, storageId string) (bool, error)
	CancelInstance(id int) error
	MigrateInstance(id int) (datatypes.Provisioning_Version1_Transaction, error)
	MigrateDedicatedHost(id int, hostId int) error
	CreateDedicatedHost(size, hostname, domain, datacenter string, billing string, routerId int) (datatypes.Container_Product_Order_Receipt, error)
	CreateInstance(template *datatypes.Virtual_Guest) (datatypes.Virtual_Guest, error)
	CreateInstances(template []datatypes.Virtual_Guest) ([]datatypes.Virtual_Guest, error)
	GenerateInstanceCreationTemplate(virtualGuest *datatypes.Virtual_Guest, params map[string]interface{}) (*datatypes.Virtual_Guest, error)
	VerifyInstanceCreation(template datatypes.Virtual_Guest) (datatypes.Container_Product_Order, error)
	GetCreateOptions(vsiType string, datacenter string) (map[string]map[string]string, error)
	GetInstance(id int, mask string) (datatypes.Virtual_Guest, error)
	GetDedicatedHost(hostId int) (datatypes.Virtual_DedicatedHost, error)
	GetLikedInstance(virtualGuest *datatypes.Virtual_Guest, id int) (*datatypes.Virtual_Guest, error)
	CaptureImage(vsId int, imageName string, imageNote string, allDisk bool) (datatypes.Provisioning_Version1_Transaction, error)
	ListInstances(hourly bool, monthly bool, domain string, hostname string, datacenter string, publicIP string, privateIP string, owner string, cpu int, memory int, network int, orderId int, tags []string, mask string) ([]datatypes.Virtual_Guest, error)
	GetInstances(mask string, objFilter filter.Filters) ([]datatypes.Virtual_Guest, error)
	PauseInstance(id int) error
	PowerOnInstance(id int) error
	PowerOffInstance(id int, soft bool, hard bool) error
	RebootInstance(id int, soft bool, hard bool) error
	ReloadInstance(id int, postURI string, sshKeys []int, imageID int) error
	ResumeInstance(id int) error
	RescueInstance(id int) error
	UpgradeInstance(id int, cpu int, memory int, network int, addDisk int, resizeDisk []int, privateCPU bool, flavor string) (datatypes.Container_Product_Order_Receipt, error)
	InstanceIsReady(id int, until time.Time) (bool, string, error)
	SetUserMetadata(id int, userdata []string) error
	SetTags(id int, tags string) error
	SetNetworkPortSpeed(id int, public bool, portSpeed int) error
	EditInstance(id int, hostname string, domain string, userdata string, tags string, publicSpeed *int, privateSpeed *int) ([]bool, []string)
	GetBandwidthData(id int, startDate time.Time, endDate time.Time, period int) ([]datatypes.Metric_Tracking_Object_Data, error)
	GetStorageDetails(id int, nasType string) ([]datatypes.Network_Storage, error)
	GetStorageCredentials(id int) (datatypes.Network_Storage_Allowed_Host, error)
	GetPortableStorage(id int) ([]datatypes.Virtual_Disk_Image, error)
	GetLocalDisks(id int) ([]datatypes.Virtual_Guest_Block_Device, error)
	GetCapacityDetail(id int) (datatypes.Virtual_ReservedCapacityGroup, error)
	CapacityList(mask string) ([]datatypes.Virtual_ReservedCapacityGroup, error)
	GetRouters(packageName string) ([]datatypes.Location_Region, error)
	GetCapacityCreateOptions(packageName string) ([]datatypes.Product_Item, error)
	GetPods() ([]datatypes.Network_Pod, error)
	GenerateInstanceCapacityCreationTemplate(reservedCapacity *datatypes.Container_Product_Order_Virtual_ReservedCapacity, params map[string]interface{}) (interface{}, error)
	GetSummaryUsage(id int, startDate time.Time, endDate time.Time, validType string, periodic int) (resp []datatypes.Metric_Tracking_Object_Data, err error)
	PlacementsGroupList(mask string) ([]datatypes.Virtual_PlacementGroup, error)
	GetPlacementGroupDetail(id int) (datatypes.Virtual_PlacementGroup, error)
	PlacementCreate(templateObject *datatypes.Virtual_PlacementGroup) (datatypes.Virtual_PlacementGroup, error)
	GetDatacenters() ([]datatypes.Location, error)
	GetAvailablePlacementRouters(id int) ([]datatypes.Hardware, error)
	GetRules() ([]datatypes.Virtual_PlacementGroup_Rule, error)
	GetUserCustomerNotificationsByVirtualGuestId(id int, mask string) ([]datatypes.User_Customer_Notification_Virtual_Guest, error)
	CreateUserCustomerNotification(virtualServerId int, userId int) (datatypes.User_Customer_Notification_Virtual_Guest, error)
	DeleteUserCustomerNotification(userCustomerNotificationId int) (resp bool, err error)
}

type virtualServerManager struct {
	VirtualGuestService  services.Virtual_Guest
	AccountService       services.Account
	PackageService       services.Product_Package
	OrderService         services.Product_Order
	DedicatedHostService services.Virtual_DedicatedHost
	OrderManager         OrderManager
	Session              *session.Session
	StorageManager       StorageManager
}

type virtualCreateOptions struct {
	Locations        []datatypes.Location
	Sizes            []datatypes.Product_Item
	Ram              []datatypes.Product_Item
	Database         []datatypes.Product_Item
	OperatingSystems []datatypes.Product_Item
	GuestCore        []datatypes.Product_Item
	PortSpeed        []datatypes.Product_Item
	GuestDisk        []datatypes.Product_Item
	Extras           []datatypes.Product_Item
}

type extras struct {
	Name   string
	Key    string
	Prices []datatypes.Product_Item_Price
}

func NewVirtualServerManager(session *session.Session) *virtualServerManager {
	return &virtualServerManager{
		services.GetVirtualGuestService(session),
		services.GetAccountService(session),
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetVirtualDedicatedHostService(session),
		NewOrderManager(session),
		session,
		NewStorageManager(session),
	}
}

//Attach portable storage to a Virtual Server.
//int id: Virtual server id.
//int portableStorageId: Portable storage id.
func (vs virtualServerManager) AttachPortableStorage(id int, portableStorageId int) (datatypes.Provisioning_Version1_Transaction, error) {
	return vs.VirtualGuestService.Id(id).AttachDiskImage(&portableStorageId)
}

//Authorize File or Block Storage to a Virtual Server.
//int id: Virtual server id.
//string storageUsername: Storage username.
func (vs virtualServerManager) AuthorizeStorage(id int, storageUsername string) (bool, error) {
	storageResult, err := vs.StorageManager.GetVolumeByUsername(storageUsername)
	if err != nil {
		return false, err
	}
	if len(storageResult) == 0 {
		return false, errors.New(T("The Storage {{.Storage}} was not found.", map[string]interface{}{"Storage": storageUsername}))
	}
	networkStorageTemplate := []datatypes.Network_Storage{
		{
			Id: storageResult[0].Id,
		},
	}
	return vs.VirtualGuestService.Id(id).AllowAccessToNetworkStorageList(networkStorageTemplate)
}

//Cancel an instance immediately, deleting all its data.
//id: the instance ID to cancel
func (vs virtualServerManager) CancelInstance(id int) error {
	_, err := vs.VirtualGuestService.Id(id).DeleteObject()
	return err
}

//Migrate an instance.
//id: the instance ID to migrate.
func (vs virtualServerManager) MigrateInstance(id int) (datatypes.Provisioning_Version1_Transaction, error) {
	resourceList, err := vs.VirtualGuestService.Id(id).Migrate()
	return resourceList, err
}

//Migrate a dedicated Host instance.
func (vs virtualServerManager) MigrateDedicatedHost(id int, hostId int) (err error) {
	return vs.VirtualGuestService.Id(id).MigrateDedicatedHost(&hostId)
}

func GetDedicatedHostPriceId(items []datatypes.Product_Item, size string, hourly bool, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		if item.KeyName == nil || *item.KeyName != size {
			continue
		}
		for _, price := range item.Prices {
			if !matchesBilling(price, hourly) {
				continue
			}
			if !matchesLocation(price, location) {
				continue
			}
			if price.Id != nil {
				return *price.Id, nil
			} else {
				return 0, errors.New(T("Price ID not found"))
			}
		}
	}
	return 0, errors.New(T("Could not find valid price for dedicated host with size= {{.KeyName}}", map[string]interface{}{"KeyName": size}))
}

//Create a dedicated host for dedicated virtual server
func (vs virtualServerManager) CreateDedicatedHost(size, hostname, domain, datacenter string, billing string, routerId int) (datatypes.Container_Product_Order_Receipt, error) {
	mask := "items[keyName,capacity,description,attributes[id,attributeTypeKeyName],itemCategory[id,categoryCode],softwareDescription[id,referenceCode,longDescription],prices],activePresets,regions[location[location[priceGroups]]]"
	packages, err := vs.PackageService.Mask(mask).Filter(filter.Path("keyName").Eq("DEDICATED_HOST").Build()).GetAllObjects()
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if len(packages) != 1 {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Ordering package is not found"))
	}
	hourly := billing == "hourly"
	location, err := GetLocation(packages[0], datacenter)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	priceId, err := GetDedicatedHostPriceId(packages[0].Items, size, hourly, location)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	order := datatypes.Container_Product_Order_Virtual_DedicatedHost{
		Container_Product_Order: datatypes.Container_Product_Order{
			Location: location.Keyname,
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{Id: sl.Int(priceId)},
			},
			PackageId:        packages[0].Id,
			UseHourlyPricing: sl.Bool(hourly),
			Hardware: []datatypes.Hardware{
				datatypes.Hardware{
					Hostname: sl.String(hostname),
					Domain:   sl.String(domain),
					PrimaryBackendNetworkComponent: &datatypes.Network_Component{
						Router: &datatypes.Hardware{
							Id: sl.Int(routerId),
						},
					},
				},
			},
		},
	}
	return vs.OrderService.PlaceOrder(&order, sl.Bool(false))
}

//Creates a new virtual server instance.
//template: the template virtual service instance to be created
func (vs virtualServerManager) CreateInstance(template *datatypes.Virtual_Guest) (datatypes.Virtual_Guest, error) {
	return vs.VirtualGuestService.CreateObject(template)
}

func (vs virtualServerManager) CreateInstances(template []datatypes.Virtual_Guest) ([]datatypes.Virtual_Guest, error) {
	return vs.VirtualGuestService.CreateObjects(template)
}

//Generate a new virtual server instance template from parameters for creation
func (vs virtualServerManager) GenerateInstanceCreationTemplate(virtualGuest *datatypes.Virtual_Guest, params map[string]interface{}) (*datatypes.Virtual_Guest, error) {
	var err error
	if params["template"] != nil {
		template := params["template"].(string)
		if _, err = os.Stat(template); os.IsNotExist(err) {
			return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("Template file: {{.Location}} does not exist.", map[string]interface{}{"Location": template}))
		}
		virtualGuest, err = getParamsFromTemplate(virtualGuest, template)
		if err != nil {
			return &datatypes.Virtual_Guest{}, err
		}
	}

	if params["like"] != nil {
		like := params["like"].(int)
		virtualGuest, err = vs.GetLikedInstance(virtualGuest, like)
		if err != nil {
			return &datatypes.Virtual_Guest{}, err
		}
	}

	if params["hostname"] != nil {
		virtualGuest.Hostname = sl.String(params["hostname"].(string))
	}
	if params["domain"] != nil {
		virtualGuest.Domain = sl.String(params["domain"].(string))
	}
	if params["flavor"] != nil {
		virtualGuest.SupplementalCreateObjectOptions =
			&datatypes.Virtual_Guest_SupplementalCreateObjectOptions{
				FlavorKeyName: sl.String(params["flavor"].(string)),
			}
	} else {
		if params["cpu"] != nil {
			virtualGuest.StartCpus = sl.Int(params["cpu"].(int))
		}
		if params["memory"] != nil {
			virtualGuest.MaxMemory = sl.Int(params["memory"].(int))
		}
		if params["dedicated"] != nil {
			virtualGuest.DedicatedAccountHostOnlyFlag = sl.Bool(params["dedicated"].(bool))
		}
	}

	if params["datacenter"] != nil {
		virtualGuest.Datacenter = &datatypes.Location{
			Name: sl.String(params["datacenter"].(string)),
		}
	}
	if params["os"] != nil {
		virtualGuest.OperatingSystemReferenceCode = sl.String(params["os"].(string))
	}

	if params["image"] != nil {
		virtualGuest.BlockDeviceTemplateGroup =
			&datatypes.Virtual_Guest_Block_Device_Template_Group{
				GlobalIdentifier: sl.String(params["image"].(string)),
			}
	}

	if params["billing"] != nil {
		virtualGuest.HourlyBillingFlag = sl.Bool(params["billing"].(bool))
	}

	if params["private"] != nil {
		virtualGuest.PrivateNetworkOnlyFlag = sl.Bool(params["private"].(bool))
	}

	if params["san"] == nil && params["flavor"] != nil {
		virtualGuest.LocalDiskFlag = sl.Bool(false)
	} else if params["san"] != nil {
		virtualGuest.LocalDiskFlag = sl.Bool(!params["san"].(bool))
	} else {
		virtualGuest.LocalDiskFlag = sl.Bool(false)
	}

	if params["i"] != nil {
		virtualGuest.PostInstallScriptUri = sl.String(params["i"].(string))
	}

	if params["sshkeys"] != nil {
		sshkeys := params["sshkeys"].([]int)
		var securityKeys []datatypes.Security_Ssh_Key
		for _, sshkey := range sshkeys {
			key := datatypes.Security_Ssh_Key{
				Id: sl.Int(sshkey),
			}
			securityKeys = append(securityKeys, key)
		}
		virtualGuest.SshKeys = securityKeys
	}

	if params["disks"] != nil {
		disks := params["disks"].([]int)
		if virtualGuest != nil && virtualGuest.LocalDiskFlag != nil && *virtualGuest.LocalDiskFlag == true && len(disks) > 2 {
			return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("Local disk number cannot excceed two."))
		}
		if virtualGuest != nil && virtualGuest.LocalDiskFlag != nil && *virtualGuest.LocalDiskFlag == false && len(disks) > 5 {
			return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("San disk number cannot excceed five."))
		}
		var blockDevices []datatypes.Virtual_Guest_Block_Device
		for index, disk := range disks {
			deviceNum := index
			capacity := disk
			if index > 0 {
				deviceNum = index + 1
			}
			deviceNumStr := strconv.Itoa(deviceNum)
			blockDevices = append(blockDevices, datatypes.Virtual_Guest_Block_Device{
				Device: &deviceNumStr,
				DiskImage: &datatypes.Virtual_Disk_Image{
					Capacity: &capacity,
				},
			})
		}
		virtualGuest.BlockDevices = blockDevices
	}

	if params["network"] != nil {
		virtualGuest.NetworkComponents =
			[]datatypes.Virtual_Guest_Network_Component{
				datatypes.Virtual_Guest_Network_Component{
					MaxSpeed: sl.Int(params["network"].(int)),
				},
			}
	}

	if params["vlan-public"] != nil {
		virtualGuest.PrimaryNetworkComponent =
			&datatypes.Virtual_Guest_Network_Component{
				NetworkVlan: &datatypes.Network_Vlan{
					Id: sl.Int(params["vlan-public"].(int)),
				},
			}

		if params["subnet-public"] != nil {
			virtualGuest.PrimaryNetworkComponent.NetworkVlan.PrimarySubnet =
				&datatypes.Network_Subnet{
					Id: sl.Int(params["subnet-public"].(int)),
				}
		}
	}

	if params["vlan-private"] != nil {
		virtualGuest.PrimaryBackendNetworkComponent =
			&datatypes.Virtual_Guest_Network_Component{
				NetworkVlan: &datatypes.Network_Vlan{
					Id: sl.Int(params["vlan-private"].(int)),
				},
			}

		if params["subnet-private"] != nil {
			virtualGuest.PrimaryBackendNetworkComponent.NetworkVlan.PrimarySubnet =
				&datatypes.Network_Subnet{
					Id: sl.Int(params["subnet-private"].(int)),
				}
		}
	}

	if params["boot-mode"] != nil {
		if virtualGuest.SupplementalCreateObjectOptions == nil {
			virtualGuest.SupplementalCreateObjectOptions =
				&datatypes.Virtual_Guest_SupplementalCreateObjectOptions{
					BootMode: sl.String(params["boot-mode"].(string)),
				}
		} else {
			virtualGuest.SupplementalCreateObjectOptions.BootMode = sl.String(params["boot-mode"].(string))
		}
	}

	if params["placement-group-id"] != nil {
		virtualGuest.PlacementGroupId = sl.Int(params["placement-group-id"].(int))
	}

	if params["transient"] != nil {
		virtualGuest.TransientGuestFlag = sl.Bool(params["transient"].(bool))
	}

	if params["host-id"] != nil {
		virtualGuest.DedicatedHost = &datatypes.Virtual_DedicatedHost{Id: sl.Int(params["host-id"].(int))}
	}

	if params["public-security-group"] != nil {
		groupIds := params["public-security-group"].([]int)
		var groups []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding
		for _, id := range groupIds {
			groups = append(groups, datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
				SecurityGroup: &datatypes.Network_SecurityGroup{
					Id: sl.Int(id),
				},
			})
		}
		if virtualGuest.PrimaryNetworkComponent == nil {
			virtualGuest.PrimaryNetworkComponent = &datatypes.Virtual_Guest_Network_Component{
				SecurityGroupBindings: groups,
			}
		} else {
			virtualGuest.PrimaryNetworkComponent.SecurityGroupBindings = groups
		}

	}
	if params["private-security-group"] != nil {
		groupIds := params["private-security-group"].([]int)
		var groups []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding
		for _, id := range groupIds {
			groups = append(groups, datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
				SecurityGroup: &datatypes.Network_SecurityGroup{
					Id: sl.Int(id),
				},
			})
		}
		if virtualGuest.PrimaryBackendNetworkComponent == nil {
			virtualGuest.PrimaryBackendNetworkComponent = &datatypes.Virtual_Guest_Network_Component{
				SecurityGroupBindings: groups,
			}
		} else {
			virtualGuest.PrimaryBackendNetworkComponent.SecurityGroupBindings = groups
		}
	}

	if params["userdata"] != nil {
		virtualGuest.UserData = []datatypes.Virtual_Guest_Attribute{
			datatypes.Virtual_Guest_Attribute{
				Value: sl.String(params["userdata"].(string)),
			},
		}
	}

	if params["postURL"] != nil {
		virtualGuest.PostInstallScriptUri = sl.String(params["postURL"].(string))
	}

	if virtualGuest.Hostname == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewMissingInputError("-H|--hostname")
	}
	if virtualGuest.Domain == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewMissingInputError("-D|--domain")
	}
	if virtualGuest.StartCpus == nil && virtualGuest.SupplementalCreateObjectOptions.FlavorKeyName == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("either [-c|--cpu] or [--flavor] is required."))
	}
	if virtualGuest.MaxMemory == nil && virtualGuest.SupplementalCreateObjectOptions.FlavorKeyName == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("either [-m|--memory] or [--flavor] is required."))
	}
	if virtualGuest.Datacenter == nil || virtualGuest.Datacenter.Name == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewMissingInputError("--datacenter")
	}
	if virtualGuest.OperatingSystemReferenceCode == nil && virtualGuest.BlockDeviceTemplateGroup == nil {
		return &datatypes.Virtual_Guest{}, bmxErr.NewInvalidUsageError(T("either [-o|--os] or [--image] is required."))
	}
	return virtualGuest, nil
}

func getParamsFromTemplate(virtualGuest *datatypes.Virtual_Guest, templateFile string) (*datatypes.Virtual_Guest, error) {
	data, err := ioutil.ReadFile(templateFile) // #nosec
	if err != nil {
		return &datatypes.Virtual_Guest{}, err
	}
	err = json.Unmarshal(data, virtualGuest)
	return virtualGuest, nil
}

//Verifies an instance creation command, without actually placing an order.
//template: the template virtual service instance to be verified
func (vs virtualServerManager) VerifyInstanceCreation(template datatypes.Virtual_Guest) (datatypes.Container_Product_Order, error) {
	return vs.VirtualGuestService.GenerateOrderTemplate(&template)
}

//Retrieves the available options for creating a virtual server instance
func (vs virtualServerManager) GetCreateOptions(vsiType string, datacenter string) (map[string]map[string]string, error) {

	virtualCreateOptionsResult := virtualCreateOptions{}
	var locationGroupId int
	packageData, err := vs.GetPackage(vsiType)
	if err != nil {
		return nil, err
	}

	var dc_detail datatypes.Location
	locations := make(map[string]string)
	if datacenter != "" {
		//filter :=`{"name": {"operation": datacenter}}`
		filter := filter.Build(filter.Path("name").Eq(datacenter))
		mask := "mask[priceGroups]"
		locationService := services.GetLocationService(vs.Session)
		dc_details, err := locationService.Mask(mask).Filter(filter).Limit(1).GetDatacenters()
		if err != nil {
			return nil, err
		}
		dc_detail = dc_details[0]
		// A DC will have several price groups, no good way to deal with this other than checking each.
		// An item should only belong to one type of price group.

		for _, group := range dc_detail.PriceGroups {
			if strings.HasPrefix(*group.Description, "Location") {
				locations[*group.Name] = *group.Description
			}
		}
	}

	for _, region := range packageData.Regions {
		if region.Location != nil && region.Location.Location != nil && region.Location.Location.Name != nil && region.Location.Location.LongName != nil {
			locations[*region.Location.Location.Name] = *region.Location.Location.LongName
		}
	}
	virtualCreateOptionsResult.Sizes = getOptionSizes(packageData, locationGroupId)

	operatingSystems := make(map[string]string)
	database := make(map[string]string)
	portSpeeds := make(map[string]string)
	guests := make(map[string]string)
	extras := make(map[string]string)

	sizes := make(map[string]string)
	for _, preset := range packageData.ActivePresets {
		if preset.KeyName != nil && preset.Description != nil {
			sizes[*preset.KeyName] = *preset.Description
		}
	}

	items := packageData.Items
	for _, item := range items {
		category := utils.FormatStringPointerName(item.ItemCategory.CategoryCode)
		item.Prices = getItemPriceByLocationGroup(item.Prices, locationGroupId)

		extraList := []string{"pri_ipv6_addresses", "static_ipv6_addresses",
			"sec_ip_addresses", "trusted_platform_module", "software_guard_extensions"}
		switch {
		case category == "os":
			operatingSystems[*item.SoftwareDescription.ReferenceCode] = *item.SoftwareDescription.LongDescription

		case category == "database":
			database[*item.KeyName] = *item.Description
		case category == "port_speed":
			portSpeeds[utils.FormatStringPointer(item.KeyName)] = *item.Description

		case strings.Contains(category, "guest_disk"):
			guests[string(*item.KeyName)] = *item.Description
		case utils.WordInList(extraList, category):
			extras[*item.KeyName] = *item.Description
		}
	}

	return map[string]map[string]string{
		KEY_LOCATIONS:  locations,
		KEY_SIZES:      sizes,
		KEY_DATABASE:   database,
		KEY_OS:         operatingSystems,
		KEY_PORT_SPEED: portSpeeds,
		KEY_GUEST:      guests,
		KEY_EXTRAS:     extras,
	}, err

}

func getItemPriceByLocationGroup(prices []datatypes.Product_Item_Price, locationGroupId int) []datatypes.Product_Item_Price {
	return prices
}

func getOptionSizes(packageData datatypes.Product_Package, locationGroupId int) []datatypes.Product_Item {
	var sizes []datatypes.Product_Item
	presets := packageData.ActivePresets
	presets = append(presets, packageData.AccountRestrictedActivePresets...)
	for _, preset := range presets {
		hourlyPrice := datatypes.Product_Item_Price{
			HourlyRecurringFee: sl.Float(getPresetCost(preset, packageData.Items, "hourly", locationGroupId)),
		}
		monthlyPrice := datatypes.Product_Item_Price{
			RecurringFee: sl.Float(getPresetCost(preset, packageData.Items, "monthly", locationGroupId)),
		}
		presetPrices := []datatypes.Product_Item_Price{hourlyPrice, monthlyPrice}
		presetItem := datatypes.Product_Item{
			Description: preset.Description,
			KeyName:     preset.KeyName,
			Prices:      presetPrices,
		}
		sizes = append(sizes, presetItem)
	}
	return sizes
}
func getPresetCost(preset datatypes.Product_Package_Preset, items []datatypes.Product_Item, typeFee string, locationId int) float64 {
	result := 1.0
	return result
}

func printAsJsonFormat(data interface{}) {
	jsonData, jsonErr := json.MarshalIndent(data, "", "    ")
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return
	}
	println(string(jsonData))
}

func (vs virtualServerManager) GetPackage(packageName string) (datatypes.Product_Package, error) {
	itemsMask := "items[id,keyName,capacity,description,attributes[id,attributeTypeKeyName]," +
		"itemCategory[id,categoryCode],softwareDescription[id,referenceCode,longDescription]," +
		"prices[categories]],"
	// The preset prices list will only have default prices. The prices->item->prices will have location specific
	presetsMask := "activePresets[prices],"
	regionMask := "regions[location[location[priceGroups]]]"

	packageMask := "mask[id," +
		itemsMask +
		presetsMask +
		regionMask +
		"]"

	packageData, err := vs.OrderManager.GetPackageByKey(packageName, packageMask)
	if err != nil {
		return datatypes.Product_Package{}, err
	}

	return packageData, nil
}

//Get details about a virtual server instance.
//id: the instance ID
//mask: mask of properties
func (vs virtualServerManager) GetInstance(id int, mask string) (datatypes.Virtual_Guest, error) {
	if mask == "" {
		mask = INSTANCE_DEFAULT_MASK
	}
	return vs.VirtualGuestService.Id(id).Mask(mask).GetObject()
}

func (vs virtualServerManager) GetDedicatedHost(hostId int) (datatypes.Virtual_DedicatedHost, error) {
	return vs.DedicatedHostService.Id(hostId).GetObject()
}

//Assgin template properties from liked instance
//virtualGuest: template instance to be assigned
//id: the ID of liked instance
func (vs virtualServerManager) GetLikedInstance(virtualGuest *datatypes.Virtual_Guest, id int) (*datatypes.Virtual_Guest, error) {
	mask := "id, hostname, domain, datacenter.name, maxCpu, maxMemory, hourlyBillingFlag, localDiskFlag, " +
		"dedicatedAccountHostOnlyFlag, privateNetworkOnlyFlag, postInstallScriptUri, userData, networkComponents[maxSpeed], operatingSystemReferenceCode"
	obj, err := vs.GetInstance(id, mask)
	if err != nil {
		return &datatypes.Virtual_Guest{}, err
	}
	virtualGuest.Id = nil
	virtualGuest.Hostname = obj.Hostname
	virtualGuest.Domain = obj.Domain
	virtualGuest.Datacenter = obj.Datacenter
	virtualGuest.StartCpus = obj.MaxCpu
	virtualGuest.MaxMemory = obj.MaxMemory
	virtualGuest.HourlyBillingFlag = obj.HourlyBillingFlag
	virtualGuest.DedicatedAccountHostOnlyFlag = obj.DedicatedAccountHostOnlyFlag
	virtualGuest.PrivateNetworkOnlyFlag = obj.PrivateNetworkOnlyFlag
	virtualGuest.PostInstallScriptUri = obj.PostInstallScriptUri
	virtualGuest.UserData = obj.UserData
	virtualGuest.NetworkComponents = []datatypes.Virtual_Guest_Network_Component{obj.NetworkComponents[0]}
	virtualGuest.OperatingSystemReferenceCode = obj.OperatingSystemReferenceCode
	virtualGuest.LocalDiskFlag = obj.LocalDiskFlag
	return virtualGuest, nil
}

//Capture one or all disks from a VS to a SoftLayer image.
//vsId: ID of instance
//imageName: name of the image to be created
//imageNote: note of the image to be created
//allDisk: set to true to include all additional attached storage devices
func (vs virtualServerManager) CaptureImage(vsId int, imageName string, imageNote string, allDisk bool) (datatypes.Provisioning_Version1_Transaction, error) {
	vsi, err := vs.GetInstance(vsId, "id,blockDevices[id,device,mountType,diskImage[id,metadataFlag,type[keyName]]]")
	if err != nil {
		return datatypes.Provisioning_Version1_Transaction{}, err
	}
	return vs.VirtualGuestService.Id(vsId).CreateArchiveTransaction(sl.String(imageName), getDisks(vsi, allDisk), sl.String(imageNote))
}

func getDisks(vs datatypes.Virtual_Guest, all bool) []datatypes.Virtual_Guest_Block_Device {
	disks := []datatypes.Virtual_Guest_Block_Device{}
	for _, disk := range vs.BlockDevices {
		//We never want metadata disks
		if disk.DiskImage != nil && disk.DiskImage.MetadataFlag != nil && *disk.DiskImage.MetadataFlag == true {
			continue
		}
		//We never want swap devices
		if disk.DiskImage != nil && disk.DiskImage.Type != nil && disk.DiskImage.Type.KeyName != nil && *disk.DiskImage.Type.KeyName == "SWAP" {
			continue
		}
		//We never want CD images
		if disk.MountType != nil && *disk.MountType == "CD" {
			continue
		}
		//Only use the first block device if we don't want additional disks
		if !all && disk.Device != nil && *disk.Device != "0" {
			continue
		}
		disks = append(disks, disk)
	}
	return disks
}

//Retrieve a list of all virtual servers on the account.
//hourly: include hourly instances
//monthly: include monthly instances
//domain: filter based on domain
//hostname: filter based on hostname
//datacenter: filter based on datacenter
//publicIP: filter based on public IP address
//privateIP: filter based on private IP address
//createdby: filter based on ID of creator
//cpu: filter based on number of CPUS
//memory: filter based on amount of memory
//network: filter based on network speed (in MBPS)
//orderId: filter based on the ID of the order which purchased this instance
//tags: filter based on list of tags
func (vs virtualServerManager) ListInstances(hourly bool, monthly bool, domain string, hostname string, datacenter string, publicIP string, privateIP string, owner string, cpu int, memory int, network int, orderID int, tags []string, mask string) ([]datatypes.Virtual_Guest, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("virtualGuests.id").OrderBy("DESC"))
	if domain != "" {
		filters = append(filters, utils.QueryFilter(domain, "virtualGuests.domain"))
	}
	if hostname != "" {
		filters = append(filters, utils.QueryFilter(hostname, "virtualGuests.hostname"))
	}
	if datacenter != "" {
		filters = append(filters, utils.QueryFilter(datacenter, "virtualGuests.datacenter.name"))
	}
	if publicIP != "" {
		filters = append(filters, utils.QueryFilter(publicIP, "virtualGuests.primaryIpAddress"))
	}
	if privateIP != "" {
		filters = append(filters, utils.QueryFilter(privateIP, "virtualGuests.primaryBackendIpAddress"))
	}
	if owner != "" {
		filters = append(filters, filter.Path("virtualGuests.billingItem.orderItem.order.userRecord.username").Eq(owner))
	}
	if cpu != 0 {
		filters = append(filters, filter.Path("virtualGuests.maxCpu").Eq(cpu))
	}
	if memory != 0 {
		filters = append(filters, filter.Path("virtualGuests.maxMemory").Eq(memory))
	}
	if network != 0 {
		filters = append(filters, filter.Path("virtualGuests.networkComponents.maxSpeed").Eq(network))
	}
	if orderID != 0 {
		filters = append(filters, filter.Path("virtualGuests.billingItem.orderItem.order.id").Eq(orderID))
	}
	if len(tags) > 0 {
		tagInterfaces := make([]interface{}, len(tags))
		for i, v := range tags {
			tagInterfaces[i] = v
		}
		filters = append(filters, filter.Path("virtualGuests.tagReferences.tag.name").In(tagInterfaces...))
	}

	if mask == "" {
		mask = INSTANCE_DEFAULT_MASK
	}

	if hourly == false && monthly == true {
		i := 0
		resourceList := []datatypes.Virtual_Guest{}
		for {
			resp, err := vs.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetMonthlyVirtualGuests()
			i++
			if err != nil {
				return []datatypes.Virtual_Guest{}, err
			}
			resourceList = append(resourceList, resp...)
			if len(resp) < metadata.LIMIT {
				break
			}
		}
		return resourceList, nil

	} else if hourly == true && monthly == false {
		i := 0
		resourceList := []datatypes.Virtual_Guest{}
		for {
			resp, err := vs.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetHourlyVirtualGuests()
			i++
			if err != nil {
				return []datatypes.Virtual_Guest{}, err
			}
			resourceList = append(resourceList, resp...)
			if len(resp) < metadata.LIMIT {
				break
			}
		}
		return resourceList, nil
	}

	i := 0
	resourceList := []datatypes.Virtual_Guest{}
	for {
		resp, err := vs.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetVirtualGuests()
		i++
		if err != nil {
			return []datatypes.Virtual_Guest{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil

}

//This method support a mask and a filter as parameters to retrieve a list of all virtual servers on the account.
func (vs virtualServerManager) GetInstances(mask string, objFilter filter.Filters) ([]datatypes.Virtual_Guest, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("virtualGuests.id").OrderBy("ASC"))
	if mask == "" {
		mask = INSTANCE_DEFAULT_MASK
	}
	if len(objFilter) > 0 {
		filters = objFilter
	}

	i := 0
	resourceList := []datatypes.Virtual_Guest{}
	for {
		resp, err := vs.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetVirtualGuests()
		i++
		if err != nil {
			return []datatypes.Virtual_Guest{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

//Pause an active virtual server.
//id: ID of virtual server instance
func (vs virtualServerManager) PauseInstance(id int) error {
	_, err := vs.VirtualGuestService.Id(id).Pause()
	return err
}

//Power on a virtual server.
//id: ID of virtual server instance
func (vs virtualServerManager) PowerOnInstance(id int) error {
	_, err := vs.VirtualGuestService.Id(id).PowerOn()
	return err
}

//Power off an active virtual server.
//id: ID of virtual server instance
//sort: perform a soft poweroff
//hard: perform a hard poweroff
func (vs virtualServerManager) PowerOffInstance(id int, soft bool, hard bool) error {
	var err error
	if soft == true && hard == false {
		_, err = vs.VirtualGuestService.Id(id).PowerOffSoft()
	} else {
		_, err = vs.VirtualGuestService.Id(id).PowerOff()
	}
	return err
}

//Reboot an active virtual server.
//id: ID of virtual server instance
//sort: perform a soft reboot
//hard: perform a hard reboot
func (vs virtualServerManager) RebootInstance(id int, soft bool, hard bool) error {
	var err error
	if soft == false && hard == false {
		_, err = vs.VirtualGuestService.Id(id).RebootDefault()
	} else if soft == true && hard == false {
		_, err = vs.VirtualGuestService.Id(id).RebootSoft()
	} else if soft == false && hard == true {
		_, err = vs.VirtualGuestService.Id(id).RebootHard()
	}
	return err
}

//Reload operating system on a virtual server
//id: ID of virtual server instance
//postURI:The URI of the post-install script to run after reload
//sshKeys: The SSH key IDs to add to the root user
//imageID: The ID of the image to load onto the server
func (vs virtualServerManager) ReloadInstance(id int, postURI string, sshKeys []int, imageID int) error {
	config := datatypes.Container_Hardware_Server_Configuration{}
	if postURI != "" {
		config.CustomProvisionScriptUri = sl.String(postURI)
	}
	if len(sshKeys) > 0 {
		config.SshKeyIds = sshKeys
	}
	if imageID != 0 {
		config.ImageTemplateId = sl.Int(imageID)
	}
	_, err := vs.VirtualGuestService.Id(id).ReloadOperatingSystem(sl.String("FORCE"), &config)
	return err
}

//Resumes a paused virtual server.
//id: ID of virtual server instance
func (vs virtualServerManager) ResumeInstance(id int) error {
	_, err := vs.VirtualGuestService.Id(id).Resume()
	return err
}

//Reboot a virtual server into a rescue image.
//id: ID of virtual server instance
func (vs virtualServerManager) RescueInstance(id int) error {
	_, err := vs.VirtualGuestService.Id(id).ExecuteRescueLayer()
	return err
}

//Upgrades a virtual server instance
//id: ID of virtual server instance
//cpu: The number of virtual CPUs to upgrade to
//memory: RAM of the virtual server to be upgraded to
//network: The port speed to set
//privateCPU: CPU will be in Private Node.
func (vs virtualServerManager) UpgradeInstance(id int, cpu int, memory int, network int, addDisk int, resizeDisk []int, privateCPU bool, flavor string) (datatypes.Container_Product_Order_Receipt, error) {
	upgradeOptions := make(map[string]int)
	public := true
	if cpu != 0 {
		upgradeOptions["guest_core"] = cpu
	}
	if memory != 0 {
		upgradeOptions["ram"] = memory / 1024
	}
	if network != -1 {
		upgradeOptions["port_speed"] = network
	}
	if privateCPU == true {
		public = false
	}

	packageItems, err := vs.VirtualGuestService.Id(id).Mask("mask[id,locationGroupId,categories[name,id,categoryCode],item[description,capacity,units]]").GetUpgradeItemPrices(sl.Bool(true))
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	prices := []datatypes.Product_Item_Price{}
	for option, value := range upgradeOptions {
		priceID := getPriceIdForUpgrade(packageItems, option, value, public)
		if priceID == -1 {
			return datatypes.Container_Product_Order_Receipt{},
				errors.New(T("Unable to find {{.Option}} option with value {{.Value}}.", map[string]interface{}{"Option": option, "Value": value}))
		}
		prices = append(prices, datatypes.Product_Item_Price{Id: &priceID})
	}

	order := datatypes.Container_Product_Order{
		ComplexType: sl.String("SoftLayer_Container_Product_Order_Virtual_Guest_Upgrade"),
		Prices:      prices,
		Properties: []datatypes.Container_Product_Order_Property{
			datatypes.Container_Product_Order_Property{
				Name:  sl.String("MAINTENANCE_WINDOW"),
				Value: sl.String(time.Now().UTC().Format(time.RFC3339)),
			},
		},
		VirtualGuests: []datatypes.Virtual_Guest{
			datatypes.Virtual_Guest{
				Id: &id,
			},
		},
	}

	if flavor != "" {
		vsObject, err := vs.GetInstance(id, "billingItem.package")
		if err != nil {
			return datatypes.Container_Product_Order_Receipt{}, err
		}
		if vsObject.BillingItem != nil && vsObject.BillingItem.Package != nil && vsObject.BillingItem.Package.KeyName != nil {
			preset, err := vs.OrderManager.GetPresetbyKey(*vsObject.BillingItem.Package.KeyName, flavor)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
			order.PresetId = preset.Id
		}
	}

	upgradeOrder := datatypes.Container_Product_Order_Virtual_Guest_Upgrade{
		Container_Product_Order_Virtual_Guest: datatypes.Container_Product_Order_Virtual_Guest{
			Container_Product_Order_Hardware_Server: datatypes.Container_Product_Order_Hardware_Server{
				Container_Product_Order: order,
			},
		},
	}

	if addDisk != -1 || len(resizeDisk) != 0 {
		itemPrices := []datatypes.Product_Item_Price{}
		if addDisk != -1 {
			//dictionary to match names used by SoftLayer_Virtual_Guest::getUpgradeItemPrices and SoftLayer_Virtual_Guest::getBlockDevices
			diskKeyNames := map[string]string{
				"First Disk":  "Disk 1",
				"Second Disk": "Disk 2",
				"Third Disk":  "Disk 3",
				"Fourth Disk": "Disk 4",
				"Fifth Disk":  "Disk 5",
			}
			diskItemId := 0
			allAvailableDiskCategoriesToNewDisk := []string{}
			description := strconv.Itoa(addDisk) + " GB (SAN)"
			for _, item := range packageItems {
				if *item.Item.Description == description {
					diskItemId = *item.Id
					allAvailableDiskCategoriesToNewDisk = append(allAvailableDiskCategoriesToNewDisk, *item.Categories[0].Name)
				}
			}
			virtualGuestDisks, err := vs.GetLocalDisks(id)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
			diskCategoriesUsedByVS := []string{}
			for _, disk := range virtualGuestDisks {
				diskCategoriesUsedByVS = append(diskCategoriesUsedByVS, *disk.DiskImage.Name)
			}
			availableDiskCategoriesToNewDiskinVS := []string{}
			for _, diskCategoy := range allAvailableDiskCategoriesToNewDisk {
				diskCategoyCounter := utils.StringInSlice(diskKeyNames[diskCategoy], diskCategoriesUsedByVS)
				if diskCategoyCounter == -1 {
					availableDiskCategoriesToNewDiskinVS = append(availableDiskCategoriesToNewDiskinVS, diskCategoy)
				}
			}
			if len(availableDiskCategoriesToNewDiskinVS) == 0 {
				return datatypes.Container_Product_Order_Receipt{},
					errors.New(T("There is not available category to this disk size"))
			}
			categoryId := 0
			for _, item := range packageItems {
				if *item.Item.Description == description && *item.Categories[0].Name == availableDiskCategoriesToNewDiskinVS[0] {
					categoryId = *item.Categories[0].Id
					break
				}
			}

			diskItem := datatypes.Product_Item_Price{
				Categories: []datatypes.Product_Item_Category{
					datatypes.Product_Item_Category{
						Id: sl.Int(categoryId),
					},
				},
				Id: sl.Int(diskItemId),
			}
			itemPrices = append(itemPrices, diskItem)
		}

		if len(resizeDisk) != 0 {
			//disk_key_names dictionary to convert disk numbers (int) to ordinal numbers (string)
			diskKeyNames := map[int]string{
				1: "First Disk",
				2: "Second Disk",
				3: "Third Disk",
				4: "Fourth Disk",
				5: "Fifth Disk",
			}
			capacity := resizeDisk[0]
			diskNumber := resizeDisk[1]
			categoryToRequest := diskKeyNames[diskNumber]
			if categoryToRequest == "" {
				return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid disk number to this disk capacity"))
			}
			description := strconv.Itoa(capacity) + " GB (SAN)"
			diskItemId := 0
			categoryId := 0
			for _, item := range packageItems {
				if *item.Item.Description == description && *item.Categories[0].Name == categoryToRequest {
					diskItemId = *item.Id
					categoryId = *item.Categories[0].Id
					break
				}
			}
			if diskItemId == 0 && categoryId == 0 {
				return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid disk number to this disk capacity"))
			}
			diskItem := datatypes.Product_Item_Price{
				Categories: []datatypes.Product_Item_Category{
					datatypes.Product_Item_Category{
						Id: sl.Int(categoryId),
					},
				},
				Id: sl.Int(diskItemId),
			}
			itemPrices = append(itemPrices, diskItem)
		}
		upgradeOrder.Container_Product_Order_Virtual_Guest.Container_Product_Order_Hardware_Server.Container_Product_Order.Prices = itemPrices
	}

	return vs.OrderService.PlaceOrder(&upgradeOrder, sl.Bool(false))
}

func getPriceIdForUpgrade(packageItems []datatypes.Product_Item_Price, option string, value int, public bool) int {

	for _, price := range packageItems {
		if len(price.Categories) == 0 || price.Item == nil {
			continue
		}
		product := price.Item

		var isPrivate bool
		if product.Units != nil && (*product.Units == "PRIVATE_CORE" || *product.Units == "DEDICATED_CORE") {
			isPrivate = true
		}

		for _, category := range price.Categories {
			if category.CategoryCode != nil && price.Id != nil {
				if !(*category.CategoryCode == option && strconv.FormatFloat(float64(*product.Capacity), 'f', 0, 64) == strconv.Itoa(value)) {
					continue
				}
				if option == "guest_core" {
					if public && !isPrivate {
						return *price.Id
					} else if !public && isPrivate {
						return *price.Id
					}
				} else if option == "port_speed" && product.Description != nil {
					if strings.Contains(*product.Description, "Public") {
						return *price.Id
					}
				} else {
					return *price.Id
				}
			}
		}
	}

	return -1
}

//Check the virtual server instance is ready for use
//param1: bool, indicate whether the instance is ready
//param2: string, indicate a possible reason if the instance is not ready
//param3: error, any error may happen when getting the status of the instance
func (vs virtualServerManager) InstanceIsReady(id int, until time.Time) (bool, string, error) {
	for {
		virtualGuest, err := vs.GetInstance(id, "id, lastOperatingSystemReload[id,modifyDate], activeTransaction[id,transactionStatus.name], provisionDate, powerState.keyName")
		if err != nil {
			return false, T("Failed to get this virtual guest instance."), err
		}

		lastReload := virtualGuest.LastOperatingSystemReload
		activeTxn := virtualGuest.ActiveTransaction
		provisionDate := virtualGuest.ProvisionDate

		// if lastReload != nil && lastReload.ModifyDate != nil {
		// 	fmt.Println("lastReload: ", (*lastReload.ModifyDate).Format(time.RFC3339))
		// }
		// if activeTxn != nil && activeTxn.TransactionStatus != nil && activeTxn.TransactionStatus.Name != nil {
		// 	fmt.Println("activeTxn: ", *activeTxn.TransactionStatus.Name)
		// }
		// if provisionDate != nil {
		// 	fmt.Println("provisionDate: ", (*provisionDate).Format(time.RFC3339))
		// }
		var reloading bool
		if activeTxn != nil && activeTxn.Id != nil && lastReload != nil && lastReload.Id != nil {
			reloading = activeTxn != nil && lastReload != nil && *activeTxn.Id == *lastReload.Id
		}
		if provisionDate != nil && !reloading {
			//fmt.Println("power state:", *virtualGuest.PowerState.KeyName)
			if virtualGuest.PowerState != nil && virtualGuest.PowerState.KeyName != nil && *virtualGuest.PowerState.KeyName == "HALTED" {
				return false, T("Virtual guest instance {{.Id}} is power off.", map[string]interface{}{"Id": id}), nil
			}
			if virtualGuest.PowerState != nil && virtualGuest.PowerState.KeyName != nil && *virtualGuest.PowerState.KeyName == "PAUSED" {
				return false, T("Virtual guest instance {{.Id}} is paused.", map[string]interface{}{"Id": id}), nil
			}

			pingable, err := vs.VirtualGuestService.Id(id).IsPingable()
			if err != nil {
				return false, T("Failed to reach virtual guest instance {{.Id}}.", map[string]interface{}{"Id": id}), err
			}
			//fmt.Println("pingable:", pingable)
			if pingable == false {
				return false, T("Virtual guest instance {{.Id}} is not reachable.", map[string]interface{}{"Id": id}), nil
			}
			return true, "", nil
		}

		now := time.Now()
		if now.After(until) {
			return false, T("Virtual guest instance {{.Id}} is loading operating system.", map[string]interface{}{"Id": id}), nil
		}

		min := math.Min(float64(1.0), float64(until.Sub(now)))
		time.Sleep(time.Duration(min) * time.Second)
	}
}

//Set user metadata for a virtual server
//id: ID of virtual server instance
//userdata: array of user data
func (vs virtualServerManager) SetUserMetadata(id int, userdata []string) error {
	_, err := vs.VirtualGuestService.Id(id).SetUserMetadata(userdata)
	return err
}

//Set tags for a virtual server
//id: ID of virtual server instance
//tags: tags to set on the VS as a comma separated list. Use the empty string to remove all tags.
func (vs virtualServerManager) SetTags(id int, tags string) error {
	_, err := vs.VirtualGuestService.Id(id).SetTags(&tags)
	return err
}

//Set network port speed for a virtual server
//id: ID of virtual server instance
//public: public network port
//portSpeed: the network port speed to be set
func (vs virtualServerManager) SetNetworkPortSpeed(id int, public bool, portSpeed int) error {
	var err error
	if public {
		_, err = vs.VirtualGuestService.Id(id).SetPublicNetworkInterfaceSpeed(&portSpeed)
	} else {
		_, err = vs.VirtualGuestService.Id(id).SetPrivateNetworkInterfaceSpeed(&portSpeed)
	}
	return err
}

//Edit hostname, domain name, notes, and/or the user data of a virtual server
//id: ID of virtual server instance
//hostname: hostname of virtual server to be updated
//domain: domain of virtual server to be updated
//userdata: userdata of virtual server to be updated
//tags: tags of virtual server to be updated
//publicSpeed: public network port spped to be updated
//privateSpeed: private network port spped to be updated
func (vs virtualServerManager) EditInstance(id int, hostname string, domain string, userdata string, tags string, publicSpeed *int, privateSpeed *int) ([]bool, []string) {
	var successes []bool
	var messages []string
	if hostname != "" || domain != "" {
		instance, err := vs.GetInstance(id, "")
		if err != nil {
			return []bool{false}, []string{err.Error()}
		}
		if hostname != "" {
			instance.Hostname = sl.String(hostname)
		}
		if domain != "" {
			instance.Domain = sl.String(domain)
		}
		// if notes != "" {
		// 	instance.Notes = sl.String(notes)
		// }
		_, err = vs.VirtualGuestService.Id(id).EditObject(&instance)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the hostname/domain of virtual server instance: {{.VsId}}.\n", map[string]interface{}{"VsId": id})+err.Error())
		} else {
			if hostname != "" {
				successes = append(successes, true)
				messages = append(messages, T("The hostname of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
			}
			if domain != "" {
				successes = append(successes, true)
				messages = append(messages, T("The domain of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
			}
			// if notes != "" {
			// 	successes = append(successes, true)
			// 	messages = append(messages, T("The note of virtual server instance: {{.VsId}} is updated.", map[string]interface{}{"VsId": id}))
			// }
		}
	}
	if userdata != "" {
		err := vs.SetUserMetadata(id, []string{userdata})
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the user data of virtual server instance: {{.VsId}}.\n", map[string]interface{}{"VsId": id})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The user data of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
		}
	}

	if tags != "" {
		err := vs.SetTags(id, tags)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the tags of virtual server instance: {{.VsId}}.\n", map[string]interface{}{"VsId": id})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The tags of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
		}
	}

	if publicSpeed != nil {
		err := vs.SetNetworkPortSpeed(id, true, *publicSpeed)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the public network speed of virtual server instance: {{.VsId}}.\n", map[string]interface{}{"VsId": id})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The public network speed of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
		}
	}

	if privateSpeed != nil {
		err := vs.SetNetworkPortSpeed(id, false, *privateSpeed)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the private network speed of virtual server instance: {{.VsId}}.\n", map[string]interface{}{"VsId": id})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The private network speed of virtual server instance: {{.VsId}} was updated.", map[string]interface{}{"VsId": id}))
		}
	}

	return successes, messages
}

// Finds the MetricTrackingObjectId for a virtual server then calls
// SoftLayer_Metric_Tracking_Object::getBandwidthData()
func (vs virtualServerManager) GetBandwidthData(id int, startDate time.Time, endDate time.Time, period int) ([]datatypes.Metric_Tracking_Object_Data, error) {
	trackingId, err := vs.VirtualGuestService.Id(id).GetMetricTrackingObjectId()
	if err != nil {
		return nil, err
	}

	trackingService := services.GetMetricTrackingObjectService(vs.Session)
	startTime := datatypes.Time{Time: startDate}
	endTime := datatypes.Time{Time: endDate}
	bandwidthData, err := trackingService.Id(trackingId).GetBandwidthData(&startTime, &endTime, nil, &period)
	return bandwidthData, err
}

//Returns the virtual server storage credentials.
//int id: Id of the virtual server
func (vs virtualServerManager) GetStorageCredentials(id int) (datatypes.Network_Storage_Allowed_Host, error) {
	mask := "mask[credential]"
	return vs.VirtualGuestService.Id(id).Mask(mask).GetAllowedHost()
}

//Returns the virtual server portable storage.
//int id: Id of the virtual server
func (vs virtualServerManager) GetPortableStorage(id int) ([]datatypes.Virtual_Disk_Image, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("portableStorageVolumes.blockDevices.guest.id").Eq(id))
	mask := "mask[billingItem[location]]"
	return vs.AccountService.Mask(mask).Filter(filters.Build()).GetPortableStorageVolumes()
}

//Returns the virtual server local disks.
//int id: Id of the virtual server
func (vs virtualServerManager) GetLocalDisks(id int) ([]datatypes.Virtual_Guest_Block_Device, error) {
	mask := "mask[diskImage]"
	return vs.VirtualGuestService.Id(id).Mask(mask).GetBlockDevices()
}

//Returns the virtual server attached network storage.
//int id: Id of the virtual server
//nas_type: storage type.
func (vs virtualServerManager) GetStorageDetails(id int, nasType string) ([]datatypes.Network_Storage, error) {
	mask := "mask[id,username,capacityGb,notes,serviceResourceBackendIpAddress,allowedVirtualGuests[id,datacenter]]"
	return vs.VirtualGuestService.Id(id).Mask(mask).GetAttachedNetworkStorages(&nasType)
}

func (vs virtualServerManager) GetCapacityDetail(id int) (datatypes.Virtual_ReservedCapacityGroup, error) {
	mask := "mask[instances[billingItem[item[keyName],category], guest], backendRouter[datacenter]]"
	reservedService := services.GetVirtualReservedCapacityGroupService(vs.Session)
	return reservedService.Mask(mask).Id(id).GetObject()
}

// Finds the Reserved Capacity groups of Account
// SoftLayer_Reserved_Capacity_Groups
func (vs virtualServerManager) CapacityList(mask string) ([]datatypes.Virtual_ReservedCapacityGroup, error) {
	if mask == "" {
		mask = "mask[availableInstanceCount, occupiedInstanceCount," +
			"instances[id, billingItem[description, hourlyRecurringFee]]," +
			" instanceCount, backendRouter[datacenter]]"
	}
	return vs.AccountService.Mask(mask).GetReservedCapacityGroups()
}

//Pulls down all backendRouterIds that are available
//A list of locations where product_package
func (vs virtualServerManager) GetRouters(packageName string) ([]datatypes.Location_Region, error) {
	productPackage, err := vs.OrderManager.GetPackageByKey(packageName, "mask[id,locations]")
	if err != nil {
		return nil, err
	}
	regions, err := vs.PackageService.Id(*productPackage.Id).GetRegions()
	return regions, err
}

//List available reserved capacity plans
func (vs virtualServerManager) GetCapacityCreateOptions(packageName string) ([]datatypes.Product_Item, error) {
	productPackage, err := vs.OrderManager.GetPackageByKey(packageName, "mask[id,locations]")
	if err != nil {
		return nil, err
	}
	items, err := vs.PackageService.Id(*productPackage.Id).GetItems()
	return items, err
}

//Get the pod details, which contains the router id
func (vs virtualServerManager) GetPods() ([]datatypes.Network_Pod, error) {
	podService := services.GetNetworkPodService(vs.Session)
	return podService.GetAllObjects()
}

func (vs virtualServerManager) GenerateInstanceCapacityCreationTemplate(reservedCapacity *datatypes.Container_Product_Order_Virtual_ReservedCapacity, params map[string]interface{}) (interface{}, error) {
	flavorId, _ := vs.OrderManager.GetPackageByKey("RESERVED_CAPACITY", "id")
	reservedCapacity.PackageId = flavorId.Id
	reservedCapacity.ComplexType = sl.String("SoftLayer_Container_Product_Order_Virtual_ReservedCapacity")
	if params["name"] != nil {
		reservedCapacity.Name = sl.String(params["name"].(string))
	}
	if params["backendRouterId"] != nil {
		reservedCapacity.BackendRouterId = sl.Int(params["backendRouterId"].(int))
	}
	if params["flavor"] != nil {
		items, _ := vs.PackageService.Id(*flavorId.Id).GetItemPrices()
		for _, item := range items {
			if *item.Item.KeyName == params["flavor"] {
				if item.LocationGroupId == nil {
					priceFlavor := datatypes.Product_Item_Price{
						Id: item.Id,
					}
					reserverItem := datatypes.Product_Item_Price{
						Id: sl.Int(217601),
					}
					reservedCapacity.Prices = append(reservedCapacity.Prices, priceFlavor)
					reservedCapacity.Prices = append(reservedCapacity.Prices, reserverItem)

				}
			}

		}

	}
	if params["test"] == true {
		return vs.OrderService.VerifyOrder(reservedCapacity)
	} else {
		return vs.OrderService.PlaceOrder(reservedCapacity, sl.Bool(false))
	}
}

func (vs virtualServerManager) GetSummaryUsage(id int, startDate time.Time, endDate time.Time, validType string, periodic int) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	trackingInstance, err := vs.VirtualGuestService.Id(id).GetMetricTrackingObject()
	trackingService := services.GetMetricTrackingObjectService(vs.Session)
	if err != nil {
		return nil, err
	}
	startTime := datatypes.Time{Time: startDate}
	endTime := datatypes.Time{Time: endDate}
	data_types := []datatypes.Container_Metric_Data_Type{
		{KeyName: sl.String(validType), SummaryType: sl.String("max")},
	}
	return trackingService.Id(*(trackingInstance.Id)).GetSummaryData(&startTime, &endTime, data_types, &periodic)

}

// Finds the placement groups of Account
// SoftLayer_Virtual_PlacementGroup
func (vs virtualServerManager) PlacementsGroupList(mask string) ([]datatypes.Virtual_PlacementGroup, error) {
	if mask == "" {
		mask = "mask[id, name, createDate, rule, guestCount, backendRouter[id, hostname]]"
	}
	return vs.AccountService.Mask(mask).GetPlacementGroups()
}

func (vs virtualServerManager) GetPlacementGroupDetail(id int) (datatypes.Virtual_PlacementGroup, error) {
	mask := "mask[id, name, createDate, rule, backendRouter[id, hostname],guests[activeTransaction[id,transactionStatus[name,friendlyName]]]]"
	reservedService := services.GetVirtualPlacementGroupService(vs.Session)
	return reservedService.Mask(mask).Id(id).GetObject()
}

func (vs virtualServerManager) GetAvailablePlacementRouters(id int) ([]datatypes.Hardware, error) {
	placementServices := services.GetVirtualPlacementGroupService(vs.Session)
	routers, err := placementServices.GetAvailableRouters(sl.Int(id))
	if err != nil {
		return nil, err
	}
	return routers, err
}

func (vs virtualServerManager) GetDatacenters() ([]datatypes.Location, error) {
	locations := services.GetLocationDatacenterService(vs.Session)
	datacenters, err := locations.GetDatacenters()
	if err != nil {
		return nil, err
	}
	return datacenters, err
}
func (vs virtualServerManager) GetRules() ([]datatypes.Virtual_PlacementGroup_Rule, error) {
	placementServices := services.GetVirtualPlacementGroupRuleService(vs.Session)
	rules, err := placementServices.GetAllObjects()
	if err != nil {
		return nil, err
	}
	return rules, err
}

func (vs virtualServerManager) PlacementCreate(templateObject *datatypes.Virtual_PlacementGroup) (datatypes.Virtual_PlacementGroup, error) {
	placementService := services.GetVirtualPlacementGroupService(vs.Session)
	return placementService.CreateObject(templateObject)
}

//Return all virtual guest notifications associated with the passed hardware ID
//int id: The virtual guest identifier.
//string mask: Object mask.
func (vs virtualServerManager) GetUserCustomerNotificationsByVirtualGuestId(id int, mask string) ([]datatypes.User_Customer_Notification_Virtual_Guest, error) {
	userCustomerNotificationVirtualGuestService := services.GetUserCustomerNotificationVirtualGuestService(vs.Session)
	if mask == "" {
		mask = "mask[id,guestId,userId,user[firstName,lastName,email,username]]"
	}
	return userCustomerNotificationVirtualGuestService.Mask(mask).FindByGuestId(&id)
}

//Create a user virtual server notification entry
//int virtualServerId: The vietual server identifier.
//int userId: The user identifier.
func (vs virtualServerManager) CreateUserCustomerNotification(virtualServerId int, userId int) (datatypes.User_Customer_Notification_Virtual_Guest, error) {
	userCustomerNotificationTemplate := datatypes.User_Customer_Notification_Virtual_Guest{
		GuestId: sl.Int(virtualServerId),
		UserId:  sl.Int(userId),
	}
	mask := "mask[user]"
	userCustomerNotificationVirtualGuestService := services.GetUserCustomerNotificationVirtualGuestService(vs.Session)
	return userCustomerNotificationVirtualGuestService.Mask(mask).CreateObject(&userCustomerNotificationTemplate)
}

//Delete a user virtual server notification entry
//int userCustomerNotificationId: The user customer notification identifier.
func (vs virtualServerManager) DeleteUserCustomerNotification(userCustomerNotificationId int) (resp bool, err error) {
	userCustomerNotificationTemplates := []datatypes.User_Customer_Notification_Virtual_Guest{
		datatypes.User_Customer_Notification_Virtual_Guest{
			Id: sl.Int(userCustomerNotificationId),
		},
	}
	userCustomerNotificationVirtualGuestService := services.GetUserCustomerNotificationVirtualGuestService(vs.Session)
	return userCustomerNotificationVirtualGuestService.Mask(mask).DeleteObjects(userCustomerNotificationTemplates)
}
