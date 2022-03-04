package managers

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	VOLUME_TYPE_BLOCK        = "block"
	VOLUME_TYPE_FILE         = "file"
	STORAGE_TYPE_PERFORMANCE = "performance"
	STORAGE_TYPE_ENDURANCE   = "endurance"
	SaaS_Category            = "storage_as_a_service"
	Space_Category           = "performance_storage_space"
	Iops_Category            = "performance_storage_iops"
	Tier_Category            = "storage_tier_level"
	Snapshot_Category        = "storage_snapshot_space"
	Replication_Category     = "performance_storage_replication"

	Saas_Order         = "Container_Product_Order_Network_Storage_AsAService"
	Saas_Order_Upgrade = "Container_Product_Order_Network_Storage_AsAService_Upgrade"

	BLOCK_VOLUME_DEFAULT_MASK = "id,username,lunId,capacityGb,bytesUsed,serviceResource.datacenter.name,serviceResourceBackendIpAddress,storageType.keyName,activeTransactionCount,billingItem.orderItem.order[id,userRecord.username],notes"
	BLOCK_VOLUME_DETAIL_MASK  = "id,username,password,capacityGb,snapshotCapacityGb,parentVolume.snapshotSizeBytes,storageType.keyName," +
		"serviceResource.datacenter.name,serviceResourceBackendIpAddress,storageTierLevel,iops,lunId," +
		"originalVolumeName,originalSnapshotName,originalVolumeSize," +
		"activeTransactionCount,activeTransactions.transactionStatus.friendlyName," +
		"replicationPartnerCount,replicationStatus," +
		"replicationPartners[id,username,serviceResourceBackendIpAddress,serviceResource.datacenter.name,replicationSchedule.type.keyname]"

	FILE_VOLUME_DEFAULT_MASK = "id,username,capacityGb,bytesUsed,serviceResource.datacenter.name,serviceResourceBackendIpAddress,activeTransactionCount,fileNetworkMountAddress,storageType.keyName,notes"
	FILE_VOLUME_DETAIL_MASK  = "id,username,password,capacityGb,bytesUsed,snapshotCapacityGb,parentVolume.snapshotSizeBytes,storageType.keyName,serviceResource.datacenter.name,serviceResourceBackendIpAddress,fileNetworkMountAddress,storageTierLevel,iops,lunId,originalVolumeName,originalSnapshotName,originalVolumeSize,activeTransactionCount,activeTransactions.transactionStatus.friendlyName,replicationPartnerCount,replicationStatus,replicationPartners[id,username,serviceResourceBackendIpAddress,serviceResource.datacenter.name,replicationSchedule.type.keyname]"
)

var (
	ENDURANCE_TIERS = map[float64]int{
		0.25: 100,
		2:    200,
		4:    300,
		10:   1000,
	}

	TIER_PER_IOPS = map[string]float64{
		"LOW_INTENSITY_TIER": 0.25,
		"READHEAVY_TIER":     2,
		"WRITEHEAVY_TIER":    4,
		"10_IOPS_PER_GB":     10,
	}
)

//Manages SoftLayer Block and File Storage volumes.
//See product information here: https://www.ibm.com/cloud-computing/bluemix/block-storage, https://www.ibm.com/cloud-computing/bluemix/file-storage
type StorageManager interface {
	SetSnapshotNotification(volumeID int, enabled bool) error
	GetSnapshotNotificationStatus(volumeId int) (int, error)
	GetVolumeAccessList(volumeId int) (datatypes.Network_Storage, error)
	AuthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error)
	DeauthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error)
	SetCredentialPassword(hostId int, password string) error
	SetLunId(volumeId int, lunId int) (datatypes.Network_Storage_Property, error)

	OrderReplicantVolume(volumeType string, volumeId int, snapshotSchedule string, location string, tier float64, iops int, osType string) (datatypes.Container_Product_Order_Receipt, error)
	FailOverToReplicant(volumeId int, replicantId int) error
	FailBackFromReplicant(volumeId int) error
	DisasterRecoveryFailover(volumeId int, replicantId int) error
	GetReplicationPartners(volumeId int) ([]datatypes.Network_Storage, error)
	GetReplicationLocations(volumeId int) ([]datatypes.Location, error)

	ListVolumes(volumeType string, datacenter string, username string, storageType string, orderId int, mask string) ([]datatypes.Network_Storage, error)
	GetVolumeDetails(volumeType string, volumeId int, mask string) (datatypes.Network_Storage, error)
	GetVolumeByUsername(username string) ([]datatypes.Network_Storage, error)
	OrderVolume(volumeType string, location string, storageType string, osType string, size int, tier float64, iops int, snapshotSize int, billing bool) (datatypes.Container_Product_Order_Receipt, error)
	CancelVolume(volumeType string, volumeId int, reason string, immediate bool) error
	OrderDuplicateVolume(config DuplicateOrderConfig) (datatypes.Container_Product_Order_Receipt, error)
	OrderModifiedVolume(volumeType string, volumeID int, newTier float64, size int, iops int) (datatypes.Container_Product_Order_Receipt, error)

	GetVolumeSnapshotList(volumeId int) ([]datatypes.Network_Storage, error)
	DeleteSnapshot(snapshotId int) error
	CreateSnapshot(volumeId int, notes string) (datatypes.Network_Storage, error)
	CancelSnapshotSpace(volumeType string, volumeId int, reason string, immediate bool) error
	EnableSnapshot(volumeId int, scheduleType string, retentionCount int, minute int, hour int, dayOfWeek string) error
	DisableSnapshots(volumeId int, scheduleType string) error
	RestoreFromSnapshot(volumeId int, snapshotId int) error
	OrderSnapshotSpace(volumeType string, volumeId int, snapshotSize int, tier float64, iops int, upgrade bool) (datatypes.Container_Product_Order_Receipt, error)
	GetVolumeSnapshotSchedules(volumeId int) (datatypes.Network_Storage, error)

	GetAllDatacenters() ([]string, error)
	GetVolumeCountLimits() ([]datatypes.Container_Network_Storage_DataCenterLimits_VolumeCountLimitContainer, error)
	VolumeRefresh(volumeId int, snapshotId int) error
	VolumeConvert(volumeId int) error
	VolumeSetNote(volumeId int, note string) (bool, error)
}

type storageManager struct {
	StorageService     services.Network_Storage
	PackageService     services.Product_Package
	OrderService       services.Product_Order
	AccountService     services.Account
	BillingService     services.Billing_Item
	LocationService    services.Location_Datacenter
	AllowedHostService services.Network_Storage_Allowed_Host
}

//Used for OrderDuplicateVolume
type DuplicateOrderConfig struct {
	// "block" or "file"
	VolumeType string
	//ID of the volume to duplicate
	OriginalVolumeId int
	//Id of the snapshot to duplicate
	OriginalSnapshotId int
	//Size of duplicate, optional. Defaults to size of original
	DuplicateSize int
	//IOPS of duplicate, optional. Defaults to IOPS of original
	DuplicateIops int
	//IOPS Tier for endurance type volumes, optional. Defaults to Tier of original
	DuplicateTier float64
	//Snapshot size of duplicate
	DuplicateSnapshotSize int
	//Create a dependent duplicate volume.
	DependentDuplicate bool
}

func NewStorageManager(session *session.Session) *storageManager {
	return &storageManager{
		services.GetNetworkStorageService(session),
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetAccountService(session),
		services.GetBillingItemService(session),
		services.GetLocationDatacenterService(session),
		services.GetNetworkStorageAllowedHostService(session),
	}
}

func (s storageManager) GetVolumeSnapshotSchedules(volumeId int) (datatypes.Network_Storage, error) {
	mask := "schedules[type,properties[type]]"
	return s.StorageService.Id(volumeId).Mask(mask).GetObject()
}

//Returns a list of authorized hosts for a specified volume.
//volumeId: ID of volume
func (s storageManager) GetVolumeAccessList(volumeId int) (datatypes.Network_Storage, error) {
	mask := "id,allowedVirtualGuests.allowedHost.credential," +
		"allowedHardware.allowedHost.credential," +
		"allowedSubnets.allowedHost.credential," +
		"allowedIpAddresses.allowedHost.credential"
	return s.StorageService.Id(volumeId).Mask(mask).GetObject()
}

//Returns a specific volume.
//string username: The volume username.
func (s storageManager) GetVolumeByUsername(username string) ([]datatypes.Network_Storage, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("networkStorage.username").Eq(username))
	return s.AccountService.Filter(filters.Build()).GetNetworkStorage()
}

//Authorizes hosts to Block/File Storage Volumes
//volumeId: The Block/File volume to authorize hosts to
//hardwareIds: A List of SoftLayer_Hardware ids
//vsIds: A List of SoftLayer_Virtual_Guest ids
//IPIds: A List of SoftLayer_Network_Subnet_IpAddress ids
//subnetIds: A List of SoftLayer_Network_Subnet ids, block volume does not support subnetIds
func (s storageManager) AuthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error) {
	templates := PopulateHostTemplates(hardwareIds, vsIds, IPIds, subnetIds)
	return s.StorageService.Id(volumeId).AllowAccessFromHostList(templates)
}

//Revokes authorization of hosts to Block/File Storage Volumes
//volumeId: The Block/File volume to authorize hosts to
//hardwareIds: A List of SoftLayer_Hardware ids
//vsIds: A List of SoftLayer_Virtual_Guest ids
//IPIds: A List of SoftLayer_Network_Subnet_IpAddress ids
//subnetIds: A List of SoftLayer_Network_Subnet ids, block volume does not support subnetIds
func (s storageManager) DeauthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error) {
	templates := PopulateHostTemplates(hardwareIds, vsIds, IPIds, subnetIds)
	return s.StorageService.Id(volumeId).RemoveAccessFromHostList(templates)
}

//Sets the password for an access host
//hostId: id of the allowed access host
//password: password to set
func (s storageManager) SetCredentialPassword(hostId int, password string) error {
	_, err := s.AllowedHostService.Id(hostId).SetCredentialPassword(&password)
	return err
}

//Set the LUN ID on a volume
//volumeId: the id of the volume
//lunId: LUN ID to set on the volume
func (s storageManager) SetLunId(volumeId int, lunId int) (datatypes.Network_Storage_Property, error) {
	return s.StorageService.Id(volumeId).CreateOrUpdateLunId(&lunId)
}

//Places an order for a replicant Block/File volume.
//volumeType: block or file
//volumeId: The ID of the primary volume to be replicated
//snapshotSchedule: The primary volume's snapshot schedule to use for replication
//location: The location for the ordered replicant volume
//tier: The tier (IOPS per GB) of endurance volume
//iops: The IOPS of performance volume
//opType: The OS type of block volume
func (s storageManager) OrderReplicantVolume(volumeType string, volumeId int, snapshotSchedule string, location string, tier float64, iops int, osType string) (datatypes.Container_Product_Order_Receipt, error) {
	mask := "billingItem[hourlyFlag,activeChildren],snapshotCapacityGb,schedules,hourlySchedule,dailySchedule,weeklySchedule,storageType.keyName,iops,storageTierLevel"
	if volumeType == VOLUME_TYPE_BLOCK {
		mask = mask + ",osType"
	}
	volume, err := s.GetVolumeDetails(volumeType, volumeId, mask)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if volumeType == VOLUME_TYPE_BLOCK && osType == "" {
		if volume.OsType != nil && volume.OsType.KeyName != nil {
			osType = *volume.OsType.KeyName
		} else {
			return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Cannot find primary volume's os-type automatically; must specify manually."))
		}
	}
	productPackage, err := GetPackage(s.PackageService, SaaS_Category)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	locationId, err := GetLocationId(s.LocationService, location)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	order, err := PrepareSaaSReplicantOrderObject(productPackage, snapshotSchedule, locationId, tier, iops, volume, volumeType)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if volumeType == VOLUME_TYPE_BLOCK {
		order.OsFormatType = &datatypes.Network_Storage_Iscsi_OS_Type{
			KeyName: sl.String(osType),
		}
	}
	order.UseHourlyPricing = volume.BillingItem.HourlyFlag
	if order.UseHourlyPricing == nil {
		order.UseHourlyPricing = sl.Bool(false)
	}

	return s.OrderService.PlaceOrder(&order, sl.Bool(false))
}

//Failover to a volume replicant.
//volumeId: The id of the volume
//replicantId: ID of replicant to failover to
func (s storageManager) FailOverToReplicant(volumeId int, replicantId int) error {
	_, err := s.StorageService.Id(volumeId).FailoverToReplicant(sl.Int(replicantId))
	return err
}

//Failback from a volume replicant.
//volumeId: The id of the volume
func (s storageManager) FailBackFromReplicant(volumeId int) error {
	_, err := s.StorageService.Id(volumeId).FailbackFromReplicant()
	return err
}

//DISASTER Failover to a volume replicant.
//If a volume (with replication) becomes inaccessible due to a disaster event,
//this method can be used to immediately failover to an available replica in another location.
//This method does not allow for fail back via the API.
//To fail back to the original volume after using this method, open a support ticket.
//To test failover, use FailOverToReplicant() instead.
//volumeId: The id of the volume that is inaccessible.
//replicantId: ID of replicant volume to make the new primary
func (s storageManager) DisasterRecoveryFailover(volumeId int, replicantId int) error {
	_, err := s.StorageService.Id(volumeId).DisasterRecoveryFailoverToReplicant(sl.Int(replicantId))
	return err
}

//Acquires list of replicant volumes pertaining to the given volume.
//volumeId: The id of the volume
func (s storageManager) GetReplicationPartners(volumeId int) ([]datatypes.Network_Storage, error) {
	return s.StorageService.Id(volumeId).GetReplicationPartners()
}

//Acquires list of the datacenters to which a volume can be replicated.
//volumeId: The id of the volume
func (s storageManager) GetReplicationLocations(volumeId int) ([]datatypes.Location, error) {
	return s.StorageService.Id(volumeId).GetValidReplicationTargetDatacenterLocations()
}

//Returns a list of block volumes.
//volumeType: block or file
//datacenter: Datacenter short name (e.g.: dal09)
//username: Name of volume.
//storageType: Type of volume: Endurance or Performance
//orderId: ID of order
func (s storageManager) ListVolumes(volumeType string, datacenter string, username string, storageType string, orderId int, mask string) ([]datatypes.Network_Storage, error) {
	filters := filter.New()
	if volumeType == VOLUME_TYPE_BLOCK {
		if mask == "" {
			mask = BLOCK_VOLUME_DEFAULT_MASK
		}
		filters = append(filters, filter.Path("iscsiNetworkStorage.storageType.keyName").NotEq("ISCSI"))
		if datacenter != "" {
			filters = append(filters, filter.Path("iscsiNetworkStorage.serviceResource.datacenter.name").Eq(datacenter))
		}
		if username != "" {
			filters = append(filters, filter.Path("iscsiNetworkStorage.username").Eq(username))
		}
		if storageType != "" {
			keyName := fmt.Sprintf("%s_BLOCK_STORAGE", strings.ToUpper(storageType))
			filters = append(filters, filter.Path("iscsiNetworkStorage.storageType.keyName").Eq(keyName))
		} else {
			filters = append(filters, filter.Path("iscsiNetworkStorage.storageType.keyName").Contains("BLOCK_STORAGE"))
		}
		if orderId != 0 {
			filters = append(filters, filter.Path("iscsiNetworkStorage.billingItem.orderItem.order.id").Eq(orderId))
		}

		//i := 0
		//var resourceList []datatypes.Network_Storage
		//for {
		//	resp, err := s.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetIscsiNetworkStorage()
		//	i++
		//	if err != nil {
		//		return []datatypes.Network_Storage{}, err
		//	}
		//	resourceList = append(resourceList, resp...)
		//	if len(resp) < metadata.LIMIT {
		//		break
		//	}
		//}
		resourceList, err := s.AccountService.Mask(mask).Filter(filters.Build()).GetIscsiNetworkStorage()
		if err != nil {
			return []datatypes.Network_Storage{}, err
		}
		return resourceList, nil

	} else if volumeType == VOLUME_TYPE_FILE {
		if mask == "" {
			mask = FILE_VOLUME_DEFAULT_MASK
		}

		filters = append(filters, filter.Path("nasNetworkStorage.serviceResource.type.type").NotEq("*NAS"))
		filters = append(filters, filter.Path("nasNetworkStorage.storageType.keyName").Contains("FILE_STORAGE"))
		if datacenter != "" {
			filters = append(filters, filter.Path("nasNetworkStorage.serviceResource.datacenter.name").Eq(datacenter))
		}
		if username != "" {
			filters = append(filters, filter.Path("nasNetworkStorage.username").Eq(username))
		}
		if storageType != "" {
			keyName := fmt.Sprintf("%s_FILE_STORAGE", strings.ToUpper(storageType))
			filters = append(filters, filter.Path("nasNetworkStorage.storageType.keyName").Eq(keyName))
		}
		if orderId != 0 {
			filters = append(filters, filter.Path("nasNetworkStorage.billingItem.orderItem.order.id").Eq(orderId))
		}

		//i := 0
		//var resourceList []datatypes.Network_Storage
		//for {
		//	resp, err := s.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetNasNetworkStorage()
		//	i++
		//	if err != nil {
		//		return []datatypes.Network_Storage{}, err
		//	}
		//	resourceList = append(resourceList, resp...)
		//	if len(resp) < metadata.LIMIT {
		//		break
		//	}
		//}
		resourceList, err := s.AccountService.Mask(mask).Filter(filters.Build()).GetNasNetworkStorage()
		if err != nil {
			return []datatypes.Network_Storage{}, err
		}
		return resourceList, nil

	} else {
		return nil, errors.New(T("Invalid volume type"))
	}
}

//Returns details about the specified volume.
//volumeType: block or file
//volumeId: ID of volume
//mask: mask of properties
func (s storageManager) GetVolumeDetails(volumeType string, volumeId int, mask string) (datatypes.Network_Storage, error) {
	if mask == "" {
		if volumeType == VOLUME_TYPE_BLOCK {
			mask = BLOCK_VOLUME_DETAIL_MASK
		} else if volumeType == VOLUME_TYPE_FILE {
			mask = FILE_VOLUME_DETAIL_MASK
		} else {
			return datatypes.Network_Storage{}, errors.New(T("Invalid volume type"))
		}
	}
	volume, err := s.StorageService.Id(volumeId).Mask(mask).GetObject()
	if err != nil {
		return datatypes.Network_Storage{}, err
	}
	return volume, nil
}

func (s storageManager) OrderModifiedVolume(volumeType string, volumeID int, newTier float64, size int, iops int) (datatypes.Container_Product_Order_Receipt, error) {
	mask_items := []string{
		"id",
		"billingItem",
		"storageType[keyName]",
		"capacityGb",
		"provisionedIops",
		"storageTierLevel",
		"staasVersion",
		"hasEncryptionAtRest",
	}
	mask := strings.Join(mask_items, ",")
	volume, err := s.GetVolumeDetails(volumeType, volumeID, mask)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if volume.BillingItem == nil {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("The volume has been cancelled; unable to modify volume."))
	}
	staasVersion, err := strconv.Atoi(utils.StringPointertoString(volume.StaasVersion))
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	if volume.HasEncryptionAtRest != nil && !(staasVersion > 1 && *volume.HasEncryptionAtRest) {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("This volume cannot be modified since it does not support Encryption at Rest."))
	}
	productPackage, err := GetPackage(s.PackageService, SaaS_Category)
	var volumeStorageType string

	if volume.StorageType != nil && volume.StorageType.KeyName != nil {
		volumeStorageType = *volume.StorageType.KeyName
	}

	var prices []datatypes.Product_Item_Price
	var volumeIsPerformance bool
	if strings.Contains(volumeStorageType, "PERFORMANCE") {
		volumeIsPerformance = true
		if size == 0 && iops == 0 {
			fmt.Println(iops)
			return datatypes.Container_Product_Order_Receipt{}, errors.New(T("A size or IOPS value must be given to modify this performance volume."))
		}
		if size == 0 && volume.CapacityGb != nil {

			size = *volume.CapacityGb

		} else if iops == 0 {
			provisionedIops, err := strconv.Atoi(utils.StringPointertoString(volume.ProvisionedIops))
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
			iops = provisionedIops
			if iops <= 0 {
				return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Cannot find volume's provisioned IOPS."))
			}
		}
		prices, err = FindVolumePricesUpgrade(productPackage, volumeType, "performance", size, newTier, iops)
	} else if strings.Contains(volumeStorageType, "ENDURANCE") {

		volumeIsPerformance = false
		if size == 0 && newTier == 0 {
			return datatypes.Container_Product_Order_Receipt{}, errors.New(T("A size or tier value must be given to modify this endurance volume."))
		}
		if size == 0 && volume.CapacityGb != nil {
			size = *volume.CapacityGb
		} else if newTier == 0 {
			tier := volume.StorageTierLevel
			iopsPerGb := 0.25
			var tierString string
			if volume.StorageType != nil {
				tierString = *tier
			}
			switch tierString {
			case "LOW_INTENSITY_TIER":
				iopsPerGb = 0.25
			case "READHEAVY_TIER":
				iopsPerGb = 2
			case "WRITEHEAVY_TIER":
				iopsPerGb = 4
			case "10_IOPS_PER_GB":
				iopsPerGb = 10
			default:
				return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Could not find tier IOPS per GB for this volume"))
			}
			newTier = iopsPerGb
		}
		prices, err = FindVolumePricesUpgrade(productPackage, volumeType, "endurance", size, newTier, iops)
	} else {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Volume does not have a valid storage type (with an appropriate \nkeyName to indicate the volume is a PERFORMANCE or an ENDURANCE volume)."))
	}

	modify_order := datatypes.Container_Product_Order_Network_Storage_AsAService_Upgrade{
		Container_Product_Order_Network_Storage_AsAService: datatypes.Container_Product_Order_Network_Storage_AsAService{
			Container_Product_Order: datatypes.Container_Product_Order{
				ComplexType: sl.String(Saas_Order_Upgrade),
				PackageId:   productPackage.Id,
				Prices:      prices,
			},
			VolumeSize: sl.Int(size),
		},
		Volume: &datatypes.Network_Storage{
			Id: &volumeID,
		},
	}
	if volumeIsPerformance {
		modify_order.Iops = &iops
	}

	return s.OrderService.PlaceOrder(&modify_order, sl.Bool(false))
}

func (s storageManager) OrderVolume(volumeType string, location string, storageType string, osType string, size int, tier float64, iops int, snapshotSize int, billing bool) (datatypes.Container_Product_Order_Receipt, error) {
	locationId, err := GetLocationId(s.LocationService, location)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid datacenter name specified. Please provide the lower case short name (e.g.: dal09)."))
	}
	productPackage, err := GetPackage(s.PackageService, SaaS_Category)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	prices, err := FindVolumePrices(productPackage, volumeType, storageType, size, tier, iops, snapshotSize)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	order := datatypes.Container_Product_Order_Network_Storage_AsAService{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType:      sl.String(Saas_Order),
			PackageId:        productPackage.Id,
			Location:         sl.String(strconv.Itoa(locationId)),
			Quantity:         sl.Int(1),
			Prices:           prices,
			UseHourlyPricing: &billing,
		},
		VolumeSize: sl.Int(size),
	}
	if volumeType == VOLUME_TYPE_BLOCK {
		order.OsFormatType = &datatypes.Network_Storage_Iscsi_OS_Type{
			KeyName: sl.String(osType),
		}
	}
	if storageType == STORAGE_TYPE_PERFORMANCE {
		order.Iops = sl.Int(iops)
	}
	return s.OrderService.PlaceOrder(&order, sl.Bool(false))
}

//Cancels the given block storage volume.
//volumeType: block or file
//volumeId: The volume ID
//reason: The reason for cancellation
//immediate: Cancel immediately or on anniversary date
func (s storageManager) CancelVolume(volumeType string, volumeId int, reason string, immediate bool) error {
	volume, err := s.GetVolumeDetails(volumeType, volumeId, "id,billingItem.id")
	if err != nil {
		return err
	}
	if volume.BillingItem == nil || volume.BillingItem.Id == nil {
		return errors.New(T("No billing item is found to cancel."))
	}
	var billitemId int
	if volume.BillingItem != nil && volume.BillingItem.Id != nil {
		billitemId = *volume.BillingItem.Id
	}
	_, err = s.BillingService.Id(billitemId).CancelItem(sl.Bool(immediate), sl.Bool(true), sl.String(reason), sl.String(""))
	return err
}

//Places an order for a duplicate block/file volume
//config A DuplicateOrderConfig entry.
func (s storageManager) OrderDuplicateVolume(config DuplicateOrderConfig) (datatypes.Container_Product_Order_Receipt, error) {
	mask := "id,billingItem.location,snapshotCapacityGb,storageType.keyName,capacityGb,originalVolumeSize,provisionedIops,storageTierLevel"
	if config.VolumeType == VOLUME_TYPE_BLOCK {
		mask = mask + ",osType.keyName"
	}
	originalVolume, err := s.GetVolumeDetails(config.VolumeType, config.OriginalVolumeId, mask)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	var osType string
	if config.VolumeType == VOLUME_TYPE_BLOCK {
		if originalVolume.OsType == nil || originalVolume.OsType.KeyName == nil {
			return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Cannot find original volume's os-type"))
		}
		osType = *originalVolume.OsType.KeyName
	}

	productPackage, err := GetPackage(s.PackageService, SaaS_Category)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	order, err := PrepareDuplicateOrderObject(productPackage, originalVolume, config)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if config.VolumeType == VOLUME_TYPE_BLOCK {
		order.OsFormatType = &datatypes.Network_Storage_Iscsi_OS_Type{KeyName: sl.String(osType)}
	}
	if config.OriginalSnapshotId != 0 {
		order.DuplicateOriginSnapshotId = sl.Int(config.OriginalSnapshotId)
	}
	if config.DependentDuplicate == true {
		// Needs to be set only if true. If this property is set AT ALL, the API will treat it as true.
		order.IsDependentDuplicateFlag = sl.Bool(true)
	}
	return s.OrderService.PlaceOrder(&order, sl.Bool(false))
}

//Returns a list of snapshots for the specified volume.
//volumeId: ID of volume
func (s storageManager) GetVolumeSnapshotList(volumeId int) ([]datatypes.Network_Storage, error) {
	mask := "id,username,notes,snapshotSizeBytes,storageType.keyName,snapshotCreationTimestamp,hourlySchedule,dailySchedule,weeklySchedule"
	return s.StorageService.Id(volumeId).Mask(mask).GetSnapshots()
}

//Deletes the specified snapshot object.
//snapshotId: The ID of the snapshot object to delete.
func (s storageManager) DeleteSnapshot(snapshotId int) error {
	_, err := s.StorageService.Id(snapshotId).DeleteObject()
	return err
}

//Creates a snapshot on the given block volume.
//volumeId: The id of the volume
//notes: The notes or "name" to assign the snapshot
func (s storageManager) CreateSnapshot(volumeId int, notes string) (datatypes.Network_Storage, error) {
	return s.StorageService.Id(volumeId).CreateSnapshot(sl.String((notes)))
}

//Cancels snapshot space for a given volume.
//volumeId: The volume ID
//reason: The reason for cancellation
//immediate: Cancel immediately or on anniversary date
func (s storageManager) CancelSnapshotSpace(volumeType string, volumeId int, reason string, immediate bool) error {
	volume, err := s.GetVolumeDetails(volumeType, volumeId, "id,billingItem.activeChildren")
	if err != nil {
		return err
	}
	if volume.BillingItem == nil || len(volume.BillingItem.ActiveChildren) == 0 {
		return errors.New(T("No snapshot space found to cancel."))
	}
	children := volume.BillingItem.ActiveChildren
	billingItemId := 0
	for _, child := range children {
		if child.CategoryCode != nil && *child.CategoryCode == Snapshot_Category {
			billingItemId = *child.Id
			break
		}
	}
	if billingItemId == 0 {
		return errors.New(T("No snapshot space found to cancel."))
	}
	_, err = s.BillingService.Id(billingItemId).CancelItem(sl.Bool(immediate), sl.Bool(true), sl.String(reason), sl.String(""))
	return err
}

//Enables snapshots for a specific block volume at a given schedule.
//volumeId: The id of the volume
//scheduleType: 'HOURLY'|'DAILY'|'WEEKLY'
//retentionCount: Number of snapshots to be kept
//minute: Minute when to take snapshot
//hour: Hour when to take snapshot
//dayOfWeek: Day when to take snapshot
func (s storageManager) EnableSnapshot(volumeId int, scheduleType string, retentionCount int, minute int, hour int, dayOfWeek string) error {
	_, err := s.StorageService.Id(volumeId).EnableSnapshots(sl.String(scheduleType), sl.Int(retentionCount), sl.Int(minute), sl.Int(hour), sl.String(dayOfWeek))
	return err
}

//Disables snapshots for a specific block volume at a given schedule.
//volumeId: The id of the volume
//scheduleType: 'HOURLY'|'DAILY'|'WEEKLY'
func (s storageManager) DisableSnapshots(volumeId int, scheduleType string) error {
	_, err := s.StorageService.Id(volumeId).DisableSnapshots(sl.String(scheduleType))
	return err
}

//Restores a specific volume from a snapshot.
//volumeId: The id of the volume
//snapshotId: The id of the restore point
func (s storageManager) RestoreFromSnapshot(volumeId int, snapshotId int) error {
	_, err := s.StorageService.Id(volumeId).RestoreFromSnapshot(sl.Int(snapshotId))
	return err
}

//Orders snapshot space for the given block/file volume.
//volumeType: block or file
//volumeId: The id of the volume
//snapshotSize: The capacity to order, in GB
//tier: The tier level of the endurance volume, in IOPS per GB
//iops: The IOSP of the performance volume
//upgrade: Flag to indicate if this order is an upgrade
func (s storageManager) OrderSnapshotSpace(volumeType string, volumeId int, snapshotSize int, tier float64, iops int, upgrade bool) (datatypes.Container_Product_Order_Receipt, error) {
	productPackage, err := GetPackage(s.PackageService, SaaS_Category)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	mask := "serviceResource.datacenter.id,billingItem,storageTierLevel,storageType.keyName,iops"
	volume, err := s.GetVolumeDetails(volumeType, volumeId, mask)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	if volume.BillingItem == nil || volume.BillingItem.CategoryCode == nil {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Unable to find category code from this volume."))
	}
	var storageType string
	if volume.StorageType != nil && volume.StorageType.KeyName != nil {
		storageType = strings.Split(*volume.StorageType.KeyName, "_")[0]
	}
	if storageType == strings.ToUpper(STORAGE_TYPE_ENDURANCE) {
		if tier == 0 {
			tier, err = FindEnduranceTierIOPSPerGB(volume)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
		}
		iops = 0
	} else if storageType == strings.ToUpper(STORAGE_TYPE_PERFORMANCE) {
		if iops == 0 && volume.Iops != nil {
			iops, err = strconv.Atoi(*volume.Iops)
			if err != nil {
				return datatypes.Container_Product_Order_Receipt{}, err
			}
		}
		tier = 0
	} else {
		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid storage type"))
	}
	spacePrice, err := FindSaasSnapshotSpacePrice(productPackage, snapshotSize, tier, iops)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	var volumeDatacenterID string
	if volume.ServiceResource != nil && volume.ServiceResource.Datacenter != nil && volume.ServiceResource.Datacenter.Id != nil {
		volumeDatacenterID = strconv.Itoa(*volume.ServiceResource.Datacenter.Id)
	}
	if upgrade {
		order := datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade{
			Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace: datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace{
				VolumeId: sl.Int(volumeId),
				Container_Product_Order: datatypes.Container_Product_Order{
					ComplexType: sl.String(("SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade")),
					PackageId:   productPackage.Id,
					Prices:      []datatypes.Product_Item_Price{spacePrice},
					Quantity:    sl.Int(1),
					Location:    sl.String(volumeDatacenterID),
				},
			},
		}
		return s.OrderService.PlaceOrder(&order, sl.Bool(false))
	}
	order := datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace{
		VolumeId: sl.Int(volumeId),
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String(("SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace")),
			PackageId:   productPackage.Id,
			Prices:      []datatypes.Product_Item_Price{spacePrice},
			Quantity:    sl.Int(1),
			Location:    sl.String(volumeDatacenterID),
		},
	}
	return s.OrderService.PlaceOrder(&order, sl.Bool(false))
}

func (s storageManager) GetAllDatacenters() ([]string, error) {
	locations, err := s.LocationService.GetDatacenters()
	if err != nil {
		return nil, err
	}
	var datacenters []string
	for _, location := range locations {
		if location.Name != nil {
			datacenters = append(datacenters, *location.Name)
		}
	}
	sort.Strings(datacenters)
	return datacenters, nil
}

//Retrieves an array of volume count limits per location and globally.
func (s storageManager) GetVolumeCountLimits() ([]datatypes.Container_Network_Storage_DataCenterLimits_VolumeCountLimitContainer, error) {
	volumeLimits, err := s.StorageService.GetVolumeCountLimits()
	return volumeLimits, err
}

//Splits a clone from its parent allowing it to be an independent volume.
//volumeId: The ID of the volume object to convert.
func (s storageManager) VolumeConvert(volumeId int) error {
	_, err := s.StorageService.Id(volumeId).ConvertCloneDependentToIndependent()
	return err
}

//Refreshes a duplicate volume with a snapshot taken from its parent.
//volumeId: The ID of the volume to refresh.
//snapshotId: The Id of the parent volume's snapshot to use as a refresh point.
func (s storageManager) VolumeRefresh(volumeId int, snapshotId int) error {
	_, err := s.StorageService.Id(volumeId).RefreshDuplicate(sl.Int(snapshotId))
	return err
}

//Add a note in a storage volume.
//volumeId: The ID of the volume to add.
//note: The note that will be added.
func (s storageManager) VolumeSetNote(volumeId int, note string) (bool, error) {
	noteTemplate := datatypes.Network_Storage{
		Notes: sl.String(note),
	}
	return s.StorageService.Id(volumeId).EditObject(&noteTemplate)
}

//Enables/Disables snapshot space usage threshold warning for a given volume.
func (s storageManager) SetSnapshotNotification(volumeID int, enabled bool) error {

	return s.StorageService.Id(volumeID).SetSnapshotNotification(&enabled)
}

//returns Enabled/Disabled snapshot space usage threshold warning for a given volume
func (s storageManager) GetSnapshotNotificationStatus(volumeId int) (int, error) {
	status, err := s.StorageService.Id(volumeId).GetSnapshotNotificationStatus()
	if err != nil {
		return -1, err
	}

	if status == "" {
		status = "1"
	}

	result, err := strconv.Atoi(status)
	return result, err
}
