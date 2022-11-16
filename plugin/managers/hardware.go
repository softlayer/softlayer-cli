package managers

import (
	"errors"
	"time"

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
	DEFAULT_HARDWARE_MASK = "id,hostname,domain,hardwareStatusId,globalIdentifier,fullyQualifiedDomainName,hardwareStatus,processorPhysicalCoreAmount,provisionDate,memoryCapacity,primaryBackendIpAddress,primaryIpAddress,networkManagementIpAddress,datacenter,operatingSystem[softwareLicense[softwareDescription[manufacturer,name,version,referenceCode]]],billingItem[id,nextInvoiceTotalRecurringAmount,children[nextInvoiceTotalRecurringAmount],orderItem.order.userRecord[username]],tagReferences[id,tag[name,id]]"
	DEFAULT_SERVER_MASK   = "mask[" + DEFAULT_HARDWARE_MASK + ",mask(SoftLayer_Hardware_Server)[activeTransaction[id,transactionStatus[name,friendlyName]]]]"
	DETAIL_HARDWARE_MASK  = "id,globalIdentifier,fullyQualifiedDomainName,hostname,domain,provisionDate,hardwareStatus,processorPhysicalCoreAmount," +
		"memoryCapacity,notes,privateNetworkOnlyFlag,primaryBackendIpAddress,primaryIpAddress,networkManagementIpAddress,userData,datacenter," +
		"networkComponents[id,status,speed,maxSpeed,name,ipmiMacAddress,ipmiIpAddress,macAddress,primaryIpAddress,port," +
		"primarySubnet[id,netmask,broadcastAddress,networkIdentifier,gateway]],hardwareChassis[id,name],activeTransaction[id,transactionStatus[friendlyName,name]]," +
		"operatingSystem[softwareLicense[softwareDescription[manufacturer,name,version,referenceCode]],passwords[username,password]]," +
		"billingItem[id,nextInvoiceTotalRecurringAmount,children[nextInvoiceTotalRecurringAmount],nextInvoiceChildren[description,categoryCode,nextInvoiceTotalRecurringAmount],orderItem.order.userRecord[username]]," +
		"hourlyBillingFlag,tagReferences[id,tag[name,id]],networkVlans[id,vlanNumber,networkSpace],remoteManagementAccounts[username,password],lastTransaction[transactionGroup],activeComponents"

	KEY_SIZES      = "sizes"
	KEY_OS         = "operating_systems"
	KEY_PORT_SPEED = "port_speed"
	KEY_LOCATIONS  = "locations"
	KEY_EXTRAS     = "extras"
)

var DEFAULT_CATEGORIES = []string{"pri_ip_addresses", "vpn_management", "remote_management"}
var EXTRA_CATEGORIES = []string{"pri_ipv6_addresses", "static_ipv6_addresses", "sec_ip_addresses"}

type HardwareServerManager interface {
	AuthorizeStorage(id int, storageId string) (bool, error)
	CancelHardware(hardwareId int, reason string, comment string, immediate bool) error
	ListHardware(tags []string, cpus int, memory int, hostname string, domain string, datacenter string, nicSpeed int, publicIP string, privateIP string, owner string, orderId int, mask string) ([]datatypes.Hardware_Server, error)
	GetHardware(hardwareId int, mask string) (datatypes.Hardware_Server, error)
	GetStorageDetails(id int, nasType string) ([]datatypes.Network_Storage, error)
	Reload(hardwareId int, postInstallURL string, sshKeys []int, upgradeBIOS bool, upgradeFirmware bool) error
	Rescure(hardwareId int) error
	PowerCycle(hardwareId int) error
	PowerOff(hardwareId int) error
	PowerOn(hardwareId int) error
	Reboot(hardwareId int, hard bool, soft bool) error
	GetStorageCredentials(id int) (datatypes.Network_Storage_Allowed_Host, error)
	GetHardDrives(id int) ([]datatypes.Hardware_Component, error)
	GetCancellationReasons() map[string]string
	GetCreateOptions(productPackage datatypes.Product_Package) map[string]map[string]string
	GenerateCreateTemplate(productPackage datatypes.Product_Package, params map[string]interface{}) (datatypes.Container_Product_Order, error)
	PlaceOrder(orderTemplate datatypes.Container_Product_Order) (datatypes.Container_Product_Order_Receipt, error)
	VerifyOrder(orderTemplate datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error)
	GetPackage() (datatypes.Product_Package, error)
	Edit(hardwareId int, userdata, hostname, domain, notes string, tags string, publicPortSpeed, privatePortSpeed int) ([]bool, []string)
	UpdateFirmware(hardwareId int, ipmi bool, raidController bool, bios bool, hardDrive bool) error
	GetExtraPriceId(items []datatypes.Product_Item, keyName string, hourly bool, location datatypes.Location_Region) (int, error)
	GetDefaultPriceId(items []datatypes.Product_Item, option string, hourly bool, location datatypes.Location_Region) (int, error)
	GetOSPriceId(items []datatypes.Product_Item, os string, location datatypes.Location_Region) (int, error)
	GetBandwidthPriceId(items []datatypes.Product_Item, hourly bool, noPublic bool, location datatypes.Location_Region) (int, error)
	GetPortSpeedPriceId(items []datatypes.Product_Item, portSpeed int, noPublic bool, location datatypes.Location_Region) (int, error)
	ToggleIPMI(hardwareID int, enabled bool) error
	GetBandwidthData(id int, startDate time.Time, endDate time.Time, period int) ([]datatypes.Metric_Tracking_Object_Data, error)
	GetHardwareComponents(id int) ([]datatypes.Hardware_Component, error)
	GetSensorData(id int, mask string) ([]datatypes.Container_RemoteManagement_SensorReading, error)
	CreateFirmwareReflashTransaction(id int) (bool, error)
	GetUserCustomerNotificationsByHardwareId(id int, mask string) ([]datatypes.User_Customer_Notification_Hardware, error)
	CreateUserCustomerNotification(hardwareId int, userId int) (datatypes.User_Customer_Notification_Hardware, error)
	GetBandwidthAllotmentDetail(hardwareId int, mask string) (datatypes.Network_Bandwidth_Version1_Allotment_Detail, error)
	GetBillingCycleBandwidthUsage(hardwareId int, mask string) ([]datatypes.Network_Bandwidth_Usage, error)
}

type hardwareServerManager struct {
	HardwareService services.Hardware_Server
	AccountService  services.Account
	PackageService  services.Product_Package
	OrderService    services.Product_Order
	LocationService services.Location_Datacenter
	BillingService  services.Billing_Item
	Session         *session.Session
	StorageManager  StorageManager
}

func NewHardwareServerManager(session *session.Session) *hardwareServerManager {
	return &hardwareServerManager{
		services.GetHardwareServerService(session),
		services.GetAccountService(session),
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetLocationDatacenterService(session),
		services.GetBillingItemService(session),
		session,
		NewStorageManager(session),
	}
}

//Authorize File or Block Storage to a Hardware Server.
//int id: Hardware server id.
//string storageUsername: Storage username.
func (hw hardwareServerManager) AuthorizeStorage(id int, storageUsername string) (bool, error) {
	storageResult, err := hw.StorageManager.GetVolumeByUsername(storageUsername)
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
	return hw.HardwareService.Id(id).AllowAccessToNetworkStorageList(networkStorageTemplate)
}

//Returns the hardware server storage credentials.
//int id: Id of the hardware server
func (hw hardwareServerManager) GetStorageCredentials(id int) (datatypes.Network_Storage_Allowed_Host, error) {
	mask := "mask[credential]"
	return hw.HardwareService.Id(id).Mask(mask).GetAllowedHost()
}

//Returns the hardware server hard drives.
//int id: Id of the hardware server
func (hw hardwareServerManager) GetHardDrives(id int) ([]datatypes.Hardware_Component, error) {
	return hw.HardwareService.Id(id).GetHardDrives()
}

//Returns the hardware server attached network storage.
//int id: Id of the hardware server
//nas_type: storage type.
func (hw hardwareServerManager) GetStorageDetails(id int, nasType string) ([]datatypes.Network_Storage, error) {
	mask := "mask[id,username,capacityGb,notes,serviceResourceBackendIpAddress,allowedHardware[id,datacenter]]"
	return hw.HardwareService.Id(id).Mask(mask).GetAttachedNetworkStorages(&nasType)
}

//Cancels the specified dedicated server.
//hardwareId: The ID of the hardware to be cancelled.
//reason: The reason code for the cancellation. This should come from :func:`GetCancellationReasons`.
//comment: An optional comment to include with the cancellation.
func (hw hardwareServerManager) CancelHardware(hardwareId int, reason string, comment string, immediate bool) error {
	reasons := hw.GetCancellationReasons()
	cancelReason := reasons[reason]
	if cancelReason == "" {
		cancelReason = reasons["unneeded"]
	}
	hwBilling, err := hw.GetHardware(hardwareId, "id,billingItem.id")
	if err != nil {
		return errors.New(T("Failed to get hardware {{.ID}}.\n", map[string]interface{}{"ID": hardwareId}) + err.Error())
	}
	if hwBilling.BillingItem == nil || hwBilling.BillingItem.Id == nil {
		return errors.New(T("No billing item found for hardware {{.ID}}.", map[string]interface{}{"ID": hardwareId}))
	}
	billingId := *hwBilling.BillingItem.Id
	_, err = hw.BillingService.Id(billingId).CancelItem(&immediate, sl.Bool(false), &cancelReason, &comment)
	return err
}

//List all hardware (servers and bare metal computing instances).
//tags: filter based on tags
//cpus: filter based on number of CPUs
//memory: filter based on amount of memory in gigabytes
//hostname: filter based on hostname
//domain: filter based on domain
//datacenter: filter based on datacenter
//nicSpeed: filter on network speed(in MBPS)
//publicIP: filter based on public IP address
//privateIP: filter based on private IP adress
//mask: mask to control what properties are returned
func (hw hardwareServerManager) ListHardware(tags []string, cpus int, memory int, hostname string, domain string, datacenter string, nicSpeed int, publicIP string, privateIP string, owner string, orderId int, mask string) ([]datatypes.Hardware_Server, error) {
	if mask == "" {
		mask = DEFAULT_HARDWARE_MASK
	}
	filters := filter.New()
	if len(tags) > 0 {
		tagInterfaces := make([]interface{}, len(tags))
		for i, v := range tags {
			tagInterfaces[i] = v
		}
		filters = append(filters, filter.Path("hardware.tagReferences.tag.name").In(tagInterfaces...))
	}
	if cpus > 0 {
		filters = append(filters, filter.Path("hardware.processorPhysicalCoreAmount").Eq(cpus))
	}
	if memory > 0 {
		filters = append(filters, filter.Path("hardware.memoryCapacity").Eq(memory))
	}
	if hostname != "" {
		filters = append(filters, utils.QueryFilter(hostname, "hardware.hostname"))
	}
	if domain != "" {
		filters = append(filters, utils.QueryFilter(domain, "hardware.domain"))
	}
	if datacenter != "" {
		filters = append(filters, utils.QueryFilter(datacenter, "hardware.datacenter.name"))
	}
	if nicSpeed > 0 {
		filters = append(filters, filter.Path("hardware.networkComponents.maxSpeed").Eq(nicSpeed))
	}
	if publicIP != "" {
		filters = append(filters, utils.QueryFilter(publicIP, "hardware.primaryIpAddress"))
	}
	if privateIP != "" {
		filters = append(filters, utils.QueryFilter(privateIP, "hardware.primaryBackendIpAddress"))
	}
	if owner != "" {
		filters = append(filters, filter.Path("hardware.billingItem.orderItem.order.userRecord.username").Eq(owner))
	}
	if orderId != 0 {
		filters = append(filters, filter.Path("hardware.billingItem.orderItem.order.id").Eq(orderId))
	}

	filters = append(filters, filter.Path("hardware.id").OrderBy("DESC"))

	i := 0
	resourceList := []datatypes.Hardware{}
	for {
		resp, err := hw.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetHardware()
		i++
		if err != nil {
			return []datatypes.Hardware_Server{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}

	servers := []datatypes.Hardware_Server{}
	for _, hw := range resourceList {
		server := datatypes.Hardware_Server{
			Hardware: hw,
		}
		servers = append(servers, server)
	}
	return servers, nil
}

//Get details about a hardware device
//hardwareId: the ID of the hardware
//mask: mask to control what properties are returned
func (hw hardwareServerManager) GetHardware(hardwareId int, mask string) (datatypes.Hardware_Server, error) {
	if mask == "" {
		mask = DETAIL_HARDWARE_MASK
	}
	return hw.HardwareService.Id(hardwareId).Mask(mask).GetObject()
}

//Perform an OS reload of a server with its current configuration.
//hardwareId: the instance ID to reload
//postInstallURL: The URI of the post-install script to run after reload
//sshKeys: The ID of SSH keys to add to the root user
//upgradeBIOS: upgrade BIOS
//upgradeFirmware: upgrade hard drive's firmware
func (hw hardwareServerManager) Reload(hardwareId int, postInstallURL string, sshKeys []int, upgradeBIOS bool, upgradeFirmware bool) error {
	config := datatypes.Container_Hardware_Server_Configuration{
		UpgradeBios:              sl.Int(utils.Bool2Int(upgradeBIOS)),
		UpgradeHardDriveFirmware: sl.Int(utils.Bool2Int(upgradeFirmware)),
	}
	if postInstallURL != "" {
		config.CustomProvisionScriptUri = sl.String(postInstallURL)
	}
	if len(sshKeys) > 0 {
		config.SshKeyIds = sshKeys
	}
	_, err := hw.HardwareService.Id(hardwareId).ReloadOperatingSystem(sl.String("FORCE"), &config)
	return err
}

//Reboot a server into the a recsue kernel.
//hardwareId: the instance ID to rescure
func (hw hardwareServerManager) Rescure(hardwareId int) error {
	_, err := hw.HardwareService.Id(hardwareId).BootToRescueLayer(nil)
	return err
}

func (hw hardwareServerManager) PowerCycle(hardwareId int) error {
	_, err := hw.HardwareService.Id(hardwareId).PowerCycle()
	return err
}
func (hw hardwareServerManager) PowerOff(hardwareId int) error {
	_, err := hw.HardwareService.Id(hardwareId).PowerOff()
	return err
}
func (hw hardwareServerManager) PowerOn(hardwareId int) error {
	_, err := hw.HardwareService.Id(hardwareId).PowerOn()
	return err
}
func (hw hardwareServerManager) Reboot(hardwareId int, hard bool, soft bool) error {
	var err error
	if hard && !soft {
		_, err = hw.HardwareService.Id(hardwareId).RebootHard()
	} else if !hard && soft {
		_, err = hw.HardwareService.Id(hardwareId).RebootSoft()
	} else if !hard && !soft {
		_, err = hw.HardwareService.Id(hardwareId).RebootDefault()
	} else {
		return errors.New(T("Can not specify both --hard and --soft"))
	}
	return err
}

//size: server size name
//hostname: server hostname
//domain: server domain name
//location: datacenter name
//os: operating system name
//portSpeed, port speed in Mbps
//sshKeys: list of IDs of SSH key
//postInstallURL: The URI of the post-install script to run after reload
//hourly:  True if using hourly pricing (default). False for monthly
//noPublic: True if this server should only have private interfaces
//extras: List of extra feature names
func (hw hardwareServerManager) GenerateCreateTemplate(productPackage datatypes.Product_Package, params map[string]interface{}) (datatypes.Container_Product_Order, error) {
	hourly := params["billing"].(string) == "hourly"
	noPublic := false
	if params["noPublic"] != nil {
		noPublic = params["noPublic"].(bool)
	}
	portSpeed := params["portSpeed"].(int)
	datacenter, err := GetLocation(productPackage, params["datacenter"].(string))
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	var prices []datatypes.Product_Item_Price
	for _, category := range DEFAULT_CATEGORIES {
		priceId, err := hw.GetDefaultPriceId(productPackage.Items, category, hourly, datacenter)
		if err != nil {
			return datatypes.Container_Product_Order{}, err
		}
		prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(priceId)})
	}
	osPriceId, err := hw.GetOSPriceId(productPackage.Items, params["osName"].(string), datacenter)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(osPriceId)})
	bandwidthPriceId, err := hw.GetBandwidthPriceId(productPackage.Items, hourly, noPublic, datacenter)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(bandwidthPriceId)})
	portSpeedPriceId, err := hw.GetPortSpeedPriceId(productPackage.Items, portSpeed, noPublic, datacenter)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(portSpeedPriceId)})

	if params["extras"] != nil {
		for _, extra := range params["extras"].([]string) {
			extraPriceId, err := hw.GetExtraPriceId(productPackage.Items, extra, hourly, datacenter)
			if err != nil {
				return datatypes.Container_Product_Order{}, err
			}
			prices = append(prices, datatypes.Product_Item_Price{Id: sl.Int(extraPriceId)})
		}
	}

	order := datatypes.Container_Product_Order{
		Hardware: []datatypes.Hardware{
			datatypes.Hardware{
				Hostname: sl.String(params["hostname"].(string)),
				Domain:   sl.String(params["domain"].(string)),
			},
		},
		Location:         datacenter.Keyname,
		Prices:           prices,
		PackageId:        productPackage.Id,
		UseHourlyPricing: sl.Bool(hourly),
	}
	presetId, err := GetPresetId(productPackage, params["size"].(string))
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	order.PresetId = sl.Int(presetId)
	if params["postInstallURL"] != nil {
		order.ProvisionScripts = []string{params["postInstallURL"].(string)}
	}
	if params["sshKeys"] != nil {
		order.SshKeys = []datatypes.Container_Product_Order_SshKeys{
			datatypes.Container_Product_Order_SshKeys{
				SshKeyIds: params["sshKeys"].([]int),
			},
		}
	}
	return order, nil
}

//Places an order for a hardware
func (hw hardwareServerManager) PlaceOrder(orderTemplate datatypes.Container_Product_Order) (datatypes.Container_Product_Order_Receipt, error) {
	return hw.OrderService.PlaceOrder(&orderTemplate, sl.Bool(false))
}

//Verify an order for a hardware
func (hw hardwareServerManager) VerifyOrder(orderTemplate datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error) {
	return hw.OrderService.VerifyOrder(&orderTemplate)
}

//Returns a dictionary of valid cancellation reasons.
func (hw hardwareServerManager) GetCancellationReasons() map[string]string {
	return map[string]string{
		"unneeded":        "No longer needed",
		"closing":         "Business closing down",
		"cost":            "Server / Upgrade Costs",
		"migrate_larger":  "Migrating to larger server",
		"migrate_smaller": "Migrating to smaller server",
		"datacenter":      "Migrating to a different SoftLayer datacenter",
		"performance":     "Network performance / latency",
		"support":         "Support response / timing",
		"sales":           "Sales process / upgrades",
		"moving":          "Moving to competitor",
	}
}

//Returns valid options for ordering hardware.
func (hw hardwareServerManager) GetCreateOptions(productPackage datatypes.Product_Package) map[string]map[string]string {
	//locations
	locations := make(map[string]string)
	for _, region := range productPackage.Regions {
		if region.Location != nil && region.Location.Location != nil && region.Location.Location.Name != nil && region.Location.Location.LongName != nil {
			locations[*region.Location.Location.Name] = *region.Location.Location.LongName
		}
	}
	//Sizes
	sizes := make(map[string]string)
	for _, preset := range productPackage.ActivePresets {
		if preset.KeyName != nil && preset.Description != nil {
			sizes[*preset.KeyName] = *preset.Description
		}
	}
	//operating system
	operatingSystems := make(map[string]string)
	portSpeeds := make(map[string]string)
	extras := make(map[string]string)
	for _, item := range productPackage.Items {
		if item.ItemCategory != nil && item.ItemCategory.CategoryCode != nil {
			if *item.ItemCategory.CategoryCode == "os" && item.SoftwareDescription != nil && item.SoftwareDescription.ReferenceCode != nil && item.SoftwareDescription.LongDescription != nil {
				operatingSystems[*item.SoftwareDescription.ReferenceCode] = *item.SoftwareDescription.LongDescription
			} else if *item.ItemCategory.CategoryCode == "port_speed" {
				if !IsPrivatePortSpeedItem(item) && IsBonded(item) && item.Description != nil {
					portSpeeds[utils.FormatSLFloatPointerToInt(item.Capacity)] = *item.Description
				}
			} else if utils.StringInSlice(*item.ItemCategory.CategoryCode, EXTRA_CATEGORIES) > -1 && item.KeyName != nil && item.Description != nil {
				extras[*item.KeyName] = *item.Description
			}
		}
	}
	return map[string]map[string]string{
		KEY_LOCATIONS:  locations,
		KEY_SIZES:      sizes,
		KEY_OS:         operatingSystems,
		KEY_PORT_SPEED: portSpeeds,
		KEY_EXTRAS:     extras,
	}
}

//Get the package related to simple hardware ordering
func (hw hardwareServerManager) GetPackage() (datatypes.Product_Package, error) {
	mask := "items[keyName,capacity,description,attributes[id,attributeTypeKeyName],itemCategory[id,categoryCode],softwareDescription[id,referenceCode,longDescription],prices],activePresets,regions[location[location[priceGroups]]]"
	filters := filter.New()
	filters = append(filters, filter.Path("keyName").Eq("BARE_METAL_SERVER"))
	packages, err := hw.PackageService.Mask(mask).Filter(filters.Build()).GetAllObjects()
	if err != nil {
		return datatypes.Product_Package{}, err
	}
	if len(packages) != 1 {
		return datatypes.Product_Package{}, errors.New(T("Ordering package is not found"))
	}
	return packages[0], nil
}

//Edit hostname, domain name, notes, user data of the hardware.
//hardwareId: the instance ID to edit
//userdata: user data on the hardware to edit. If none exist it will be created
//hostname: valid hostname
//domain: valid domain name
//notes: notes about this particular hardware
//tags: tags to set on the hardware as a comma separated list. Use the empty string to remove all tags.
func (hw hardwareServerManager) Edit(hardwareId int, userdata, hostname, domain, notes string, tags string, publicPortSpeed, privatePortSpeed int) ([]bool, []string) {
	var successes []bool
	var messages []string
	var err error
	if userdata != "" {
		_, err := hw.HardwareService.Id(hardwareId).SetUserMetadata([]string{userdata})
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the user data of hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The user data of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
		}
	}
	if tags != "" {
		_, err := hw.HardwareService.Id(hardwareId).SetTags(&tags)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the tags of hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The tags of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
		}
	}
	if hostname != "" || domain != "" || notes != "" {
		var hardware datatypes.Hardware_Server
		if hostname != "" {
			hardware.Hostname = sl.String(hostname)
		}
		if domain != "" {
			hardware.Domain = sl.String(domain)
		}
		if notes != "" {
			hardware.Notes = sl.String(notes)
		}
		_, err := hw.HardwareService.Id(hardwareId).EditObject(&hardware)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the hostname/domain of hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error())
		} else {
			if hostname != "" {
				successes = append(successes, true)
				messages = append(messages, T("The hostname of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
			}
			if domain != "" {
				successes = append(successes, true)
				messages = append(messages, T("The domain of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
			}
			if notes != "" {
				successes = append(successes, true)
				messages = append(messages, T("The note of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
			}
		}
	}

	if publicPortSpeed > 0 {
		_, err = hw.HardwareService.Id(hardwareId).SetPublicNetworkInterfaceSpeed(sl.Int(publicPortSpeed), nil)

		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the public network speed of hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The public network speed of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
		}
	}
	if privatePortSpeed > 0 {
		_, err = hw.HardwareService.Id(hardwareId).SetPrivateNetworkInterfaceSpeed(sl.Int(privatePortSpeed), nil)
		if err != nil {
			successes = append(successes, false)
			messages = append(messages, T("Failed to update the private network speed of hardware server: {{.ID}}.\n", map[string]interface{}{"ID": hardwareId})+err.Error())
		} else {
			successes = append(successes, true)
			messages = append(messages, T("The private network speed of hardware server: {{.ID}} was updated.", map[string]interface{}{"ID": hardwareId}))
		}
	}
	return successes, messages
}

//Update hardware firmware
//hardwareId: The ID of the hardware to have its firmware updatd
//ipmi: update the ipmi firmware
//raidController: update the raid controller firmware
//bios: update the bios firmware
//hardDrive: update the hard drive firmware
func (hw hardwareServerManager) UpdateFirmware(hardwareId int, ipmi bool, raidController bool, bios bool, hardDrive bool) error {
	_, err := hw.HardwareService.Id(hardwareId).CreateFirmwareUpdateTransaction(
		sl.Int(utils.Bool2Int(ipmi)),
		sl.Int(utils.Bool2Int(raidController)),
		sl.Int(utils.Bool2Int(bios)),
		sl.Int(utils.Bool2Int(hardDrive)))
	return err
}

//Return a price ID attached to item with the given keyName
func (hw hardwareServerManager) GetExtraPriceId(items []datatypes.Product_Item, keyName string, hourly bool, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		if item.KeyName != nil && *item.KeyName != keyName {
			continue
		}
		for _, price := range item.Prices {
			if !matchesBilling(price, hourly) {
				continue
			}
			if !matchesLocation(price, location) {
				continue
			}
			return *price.Id, nil
		}
	}
	return 0, errors.New(T("Could not find valid price for extra option {{.KeyName}}", map[string]interface{}{"KeyName": keyName}))
}

//Returns a 'free' price id given an option
func (hw hardwareServerManager) GetDefaultPriceId(items []datatypes.Product_Item, option string, hourly bool, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != option {
			continue
		}
		for _, price := range item.Prices {
			if price.RecurringFee != nil && *price.RecurringFee == 0 &&
				price.HourlyRecurringFee != nil && *price.HourlyRecurringFee == 0 &&
				matchesBilling(price, hourly) &&
				matchesLocation(price, location) {
				return *price.Id, nil
			}
		}
	}
	return 0, errors.New(T("Could not find valid price for {{.KeyName}} option", map[string]interface{}{"KeyName": option}))
}

//Choose a valid price id for bandwidth.
func (hw hardwareServerManager) GetBandwidthPriceId(items []datatypes.Product_Item, hourly bool, noPublic bool, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		var capacity float64
		if item.Capacity != nil {
			capacity = float64(*item.Capacity)
		}
		if item.ItemCategory == nil ||
			item.ItemCategory.CategoryCode == nil ||
			*item.ItemCategory.CategoryCode != "bandwidth" ||
			(hourly || noPublic) && capacity != 0.0 ||
			!(hourly || noPublic) && capacity == 0.0 {
			continue
		}
		for _, price := range item.Prices {
			if !matchesBilling(price, hourly) {
				continue
			}
			if !matchesLocation(price, location) {
				continue
			}
			return *price.Id, nil
		}
	}
	return 0, errors.New(T("Could not find valid price for bandwidth option"))
}

//Returns the price id matching
func (hw hardwareServerManager) GetOSPriceId(items []datatypes.Product_Item, os string, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != "os" ||
			item.SoftwareDescription == nil || item.SoftwareDescription.ReferenceCode == nil || *item.SoftwareDescription.ReferenceCode != os {
			continue
		}
		for _, price := range item.Prices {
			if !matchesLocation(price, location) {
				continue
			}
			return *price.Id, nil
		}
	}
	return 0, errors.New(T("Could not find valid price for os {{.OS}}", map[string]interface{}{"OS": os}))
}

//Choose a valid price id for port speed
func (hw hardwareServerManager) GetPortSpeedPriceId(items []datatypes.Product_Item, portSpeed int, noPublic bool, location datatypes.Location_Region) (int, error) {
	for _, item := range items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != "port_speed" {
			continue
		}
		if item.Capacity == nil || float64(*item.Capacity) != float64(portSpeed) ||
			IsPrivatePortSpeedItem(item) != noPublic ||
			!IsBonded(item) {
			continue
		}
		for _, price := range item.Prices {
			if !matchesLocation(price, location) {
				continue
			}
			return *price.Id, nil
		}
	}
	return 0, errors.New(T("Could not find valid price for port speed:"))
}

// ToggleIPMI will create a transaction to enable/disable IPMI interface of the given hardware
func (hw hardwareServerManager) ToggleIPMI(hardwareID int, enabled bool) error {
	_, err := hw.HardwareService.Id(hardwareID).ToggleManagementInterface(&enabled)
	return err
}

// Finds the MetricTrackingObjectId for a hardware server then calls
// SoftLayer_Metric_Tracking_Object::getBandwidthData()
func (hw hardwareServerManager) GetBandwidthData(id int, startDate time.Time, endDate time.Time, period int) ([]datatypes.Metric_Tracking_Object_Data, error) {
	trackingId, err := hw.HardwareService.Id(id).GetMetricTrackingObjectId()
	if err != nil {
		return nil, err
	}

	trackingService := services.GetMetricTrackingObjectService(hw.Session)
	startTime := datatypes.Time{Time: startDate}
	endTime := datatypes.Time{Time: endDate}
	bandwidthData, err := trackingService.Id(trackingId).GetBandwidthData(&startTime, &endTime, nil, &period)
	return bandwidthData, err
}

//Return True if the price object is hourly and/or monthly
func matchesBilling(price datatypes.Product_Item_Price, hourly bool) bool {
	if hourly && price.HourlyRecurringFee != nil {
		return true
	}
	if !hourly && price.RecurringFee != nil {
		return true
	}
	return false
}

//Return True if the price object matches the location.
func matchesLocation(price datatypes.Product_Item_Price, location datatypes.Location_Region) bool {
	if price.LocationGroupId == nil {
		return true
	}
	for _, group := range location.Location.Location.PriceGroups {
		if group.Id != nil && price.LocationGroupId != nil && *group.Id == *price.LocationGroupId {
			return true
		}
	}
	return false
}

//Determine if the port speed item is private network only
func IsPrivatePortSpeedItem(item datatypes.Product_Item) bool {
	for _, attribute := range item.Attributes {
		if attribute.AttributeTypeKeyName != nil && *attribute.AttributeTypeKeyName == "IS_PRIVATE_NETWORK_ONLY" {
			return true
		}
	}
	return false
}

//Determine if the item refers to a bonded port
func IsBonded(item datatypes.Product_Item) bool {
	for _, attribute := range item.Attributes {
		if attribute.AttributeTypeKeyName != nil && *attribute.AttributeTypeKeyName == "NON_LACP" {
			return false
		}
	}
	return true
}

// Get the longer key with a short location name
func GetLocation(productPackage datatypes.Product_Package, location string) (datatypes.Location_Region, error) {
	for _, region := range productPackage.Regions {
		if region.Location != nil && region.Location.Location != nil && *region.Location.Location.Name == location {
			return region, nil
		}
	}
	return datatypes.Location_Region{}, errors.New(T("Invalid datacenter name specified."))
}

//Get the preset id given the keyName of the preset
func GetPresetId(productPackage datatypes.Product_Package, size string) (int, error) {
	for _, preset := range productPackage.ActivePresets {
		if preset.KeyName != nil && *preset.KeyName == size {
			return *preset.Id, nil
		}
	}
	return 0, errors.New(T("Could not find valid size for: {{.Size}}", map[string]interface{}{"Size": size}))
}

//Returns hardware server components.
//int id: The hardware server identifier.
func (hw hardwareServerManager) GetHardwareComponents(id int) ([]datatypes.Hardware_Component, error) {
	objectMask := "mask[id,hardwareComponentModel[longDescription,hardwareGenericComponentModel[description,hardwareComponentType[keyName]],firmwares[createDate,version]]]"
	objectFilter := filter.New()
	objectFilter = append(objectFilter, filter.Path("components.hardwareComponentModel.firmwares.createDate").OrderBy("DESC"))
	return hw.HardwareService.Id(id).Mask(objectMask).Filter(objectFilter.Build()).GetComponents()
}

//Returns hardware server sensor data.
//int id: The hardware server identifier.
//string mask: Object mask.
func (hw hardwareServerManager) GetSensorData(id int, mask string) ([]datatypes.Container_RemoteManagement_SensorReading, error) {
	if mask == "" {
		mask = "mask[sensorId,status,sensorReading,lowerCritical,lowerNonCritical,upperNonCritical,upperCritical]"
	}
	return hw.HardwareService.Id(id).Mask(mask).GetSensorData()
}

//Create a transaction to reflash firmware.
//int id: The hardware server identifier.
func (hw hardwareServerManager) CreateFirmwareReflashTransaction(id int) (bool, error) {
	// int ipmi: Reflash the ipmi firmware.
	// int raid_controller: Reflash the raid controller firmware.
	// int bios: Reflash the bios firmware.
	// the values were set as 1 to represent true
	ipmi := 1
	raidController := 1
	bios := 1
	return hw.HardwareService.Id(id).CreateFirmwareReflashTransaction(&ipmi, &raidController, &bios)
}

//Return all hardware notifications associated with the passed hardware ID
//int id: The hardware server identifier.
//string mask: Object mask.
func (hw hardwareServerManager) GetUserCustomerNotificationsByHardwareId(id int, mask string) ([]datatypes.User_Customer_Notification_Hardware, error) {
	UserCustomerNotificationHardwareService := services.GetUserCustomerNotificationHardwareService(hw.Session)
	if mask == "" {
		mask = "mask[hardwareId,user[firstName,lastName,email,username]]"
	}
	return UserCustomerNotificationHardwareService.Mask(mask).FindByHardwareId(&id)
}

//Create a user hardware notification entry
//int hardwareId: The hardware server identifier.
//int userId: The user identifier.
func (hw hardwareServerManager) CreateUserCustomerNotification(hardwareId int, userId int) (datatypes.User_Customer_Notification_Hardware, error) {
	userCustomerNotificationTemplate := datatypes.User_Customer_Notification_Hardware{
		HardwareId: sl.Int(hardwareId),
		UserId:     sl.Int(userId),
	}
	userCustomerNotificationHardwareService := services.GetUserCustomerNotificationHardwareService(hw.Session)
	return userCustomerNotificationHardwareService.CreateObject(&userCustomerNotificationTemplate)
}

// Return hardware’s allotted detail record.
// int hardwareId: The hardware server identifier.
// string mask: The Object mask.
func (hw hardwareServerManager) GetBandwidthAllotmentDetail(hardwareId int, mask string) (datatypes.Network_Bandwidth_Version1_Allotment_Detail, error) {
	if mask == "" {
		mask = "mask[allocation[amount]]"
	}
	return hw.HardwareService.Id(hardwareId).Mask(mask).GetBandwidthAllotmentDetail()
}

// Retrieve The raw bandwidth usage data for the current billing cycle.
// int hardwareId: The hardware server identifier.
// string mask: The Object mask.
func (hw hardwareServerManager) GetBillingCycleBandwidthUsage(hardwareId int, mask string) ([]datatypes.Network_Bandwidth_Usage, error) {
	if mask == "" {
		mask = "mask[amountIn,amountOut,type]"
	}
	return hw.HardwareService.Id(hardwareId).Mask(mask).GetBillingCycleBandwidthUsage()
}
