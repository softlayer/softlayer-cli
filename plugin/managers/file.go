package managers

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/softlayer/softlayer-go/datatypes"
// 	"github.com/softlayer/softlayer-go/filter"
// 	"github.com/softlayer/softlayer-go/services"
// 	"github.com/softlayer/softlayer-go/session"
// 	"github.com/softlayer/softlayer-go/sl"
// 	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
// )

// //Manages SoftLayer File Storage volumes.
// //See product information here: https://www.ibm.com/cloud-computing/bluemix/file-storage
// type FileStorageManager interface {
// 	ListFileVolumes(datacenter string, username string, storageType string, orderId int) ([]datatypes.Network_Storage, error)
// }

// type fileStorageManager struct {
// 	StorageService  services.Network_Storage
// 	PackageService  services.Product_Package
// 	OrderService    services.Product_Order
// 	AccountService  services.Account
// 	BillingService  services.Billing_Item
// 	LocationService services.Location_Datacenter
// }

// func NewFileStorageManager(session *session.Session) *fileStorageManager {
// 	return &fileStorageManager{
// 		services.GetNetworkStorageService(session),
// 		services.GetProductPackageService(session),
// 		services.GetProductOrderService(session),
// 		services.GetAccountService(session),
// 		services.GetBillingItemService(session),
// 		services.GetLocationDatacenterService(session),
// 	}
// }

// //Returns a list of file volumes.
// //datacenter: Datacenter short name (e.g.: dal09)
// //username: Name of volume.
// //storageType: Type of volume: Endurance or Performance
// //orderId: ID of order
// func (f fileStorageManager) ListFileVolumes(datacenter string, username string, storageType string, orderId int) ([]datatypes.Network_Storage, error) {
// 	filters := filter.New()
// 	filters = append(filters, filter.Path("nasNetworkStorage.serviceResource.type.type").NotEq("*NAS"))
// 	filters = append(filters, filter.Path("nasNetworkStorage.storageType.keyName").Contains("FILE_STORAGE"))
// 	if datacenter != "" {
// 		filters = append(filters, filter.Path("nasNetworkStorage.serviceResource.datacenter.name").Eq(datacenter))
// 	}
// 	if username != "" {
// 		filters = append(filters, filter.Path("nasNetworkStorage.username").Eq(username))
// 	}
// 	if storageType != "" {
// 		keyName := fmt.Sprintf("%s_FILE_STORAGE", strings.ToUpper(storageType))
// 		filters = append(filters, filter.Path("nasNetworkStorage.storageType.keyName").Eq(keyName))
// 	}
// 	if orderId != 0 {
// 		filters = append(filters, filter.Path("nasNetworkStorage.billingItem.orderItem.order.id").Eq(orderId))
// 	}
// 	return f.AccountService.Mask(FILE_VOLUME_DEFAULT_MASK).Filter(filters.Build()).GetNasNetworkStorage()
// }

// //Returns details about the specified volume
// //volumeId: ID of volume
// func (f fileStorageManager) GetFileVolumeDetails(volumeId int, mask string) (datatypes.Network_Storage, error) {
// 	if mask == "" {
// 		mask = FILE_VOLUME_DETAIL_MASK
// 	}
// 	volume, err := f.StorageService.Id(volumeId).Mask(mask).GetObject()
// 	if err != nil {
// 		return datatypes.Network_Storage{}, err
// 	}
// 	return volume, nil
// }

// //Returns a list of authorized hosts for a specified volume
// //volumeId: ID of volume
// func (f fileStorageManager) GetFileVolumeAccessList(volumeId int) (datatypes.Network_Storage, error) {
// 	mask := "id,allowedVirtualGuests.allowedHost.credential," +
// 		"allowedHardware.allowedHost.credential," +
// 		"allowedSubnets.allowedHost.credential," +
// 		"allowedIpAddresses.allowedHost.credential"
// 	return f.StorageService.Id(volumeId).Mask(mask).GetObject()
// }

// //Returns a list of snapshots for the specified volume.
// //volumeId: ID of volume
// func (f fileStorageManager) GetFileVolumeSnapshotList(volumeId int) ([]datatypes.Network_Storage, error) {
// 	mask := "id,username,notes,snapshotSizeBytes,storageType.keyName,snapshotCreationTimestamp,hourlySchedule,dailySchedule,weeklySchedule"
// 	return f.StorageService.Id(volumeId).Mask(mask).GetSnapshots()
// }

// //Authorizes hosts to File Storage Volumes
// //volumeId: The Block volume to authorize hosts to
// //hardwareIds: A List of SoftLayer_Hardware ids
// //vsIds: A List of SoftLayer_Virtual_Guest ids
// //IPIds: A List of SoftLayer_Network_Subnet_IpAddress ids
// //subnetIds: A List of SoftLayer_Network_Subnet ids
// func (f fileStorageManager) AuthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error) {
// 	templates := PopulateHostTemplates(hardwareIds, vsIds, IPIds, subnetIds)
// 	return f.StorageService.Id(volumeId).AllowAccessFromHostList(templates)
// }

// //Revokes authorization of hosts to File Storage Volumes
// //volumeId: The Block volume to authorize hosts to
// //hardwareIds: A List of SoftLayer_Hardware ids
// //vsIds: A List of SoftLayer_Virtual_Guest ids
// //IPIds: A List of SoftLayer_Network_Subnet_IpAddress ids
// //subnetIds: A List of SoftLayer_Network_Subnet ids
// func (f fileStorageManager) DeauthorizeHostToVolume(volumeId int, hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) ([]datatypes.Network_Storage_Allowed_Host, error) {
// 	templates := PopulateHostTemplates(hardwareIds, vsIds, IPIds, subnetIds)
// 	return f.StorageService.Id(volumeId).RemoveAccessFromHostList(templates)
// }

// //Places an order for a replicant file volume.
// //volumeId: The ID of the primary volume to be replicated
// //snapshotSchedule: The primary volume's snapshot schedule to use for replication
// //location: The location for the ordered replicant volume
// //tier: The tier (IOPS per GB) of the primary volume
// //iops
// func (f fileStorageManager) OrderReplicantVolume(volumeId int, snapshotSchedule string, location string, tier float64, iops int) (datatypes.Container_Product_Order_Receipt, error) {
// 	mask := "billingItem.activeChildren,snapshotCapacityGb,schedules,hourlySchedule,dailySchedule,weeklySchedule,storageType.keyName,iops,storageTierLevel"
// 	fileVolume, err := f.GetFileVolumeDetails(volumeId, mask)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	productPackage, err := GetPackage(f.PackageService, SaaS_Category)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	locationId, err := GetLocationId(f.LocationService, location)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	order, err := PrepareSaaSReplicantOrderObject(productPackage, snapshotSchedule, locationId, tier, iops, fileVolume, "file")
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	return f.OrderService.PlaceOrder(&order, sl.Bool(false))
// }

// //Acquires list of replicant volumes pertaining to the given volume.
// //volumeId: The id of the volume
// func (f fileStorageManager) GetReplicationPartners(volumeId int) ([]datatypes.Network_Storage, error) {
// 	return f.StorageService.Id(volumeId).GetReplicationPartners()
// }

// //Acquires list of the datacenters to which a volume can be replicated.
// //volumeId: The id of the volume
// func (f fileStorageManager) GetReplicationLocations(volumeId int) ([]datatypes.Location, error) {
// 	return f.StorageService.Id(volumeId).GetValidReplicationTargetDatacenterLocations()
// }

// //Places an order for a duplicate file volume
// //originalVolumeId: The ID of the origin volume to be duplicated
// //originalSnapshotId: Origin snapshot ID to use for duplication
// //duplicateSize: Size/capacity for the duplicate volume
// //duplicateIops: The IOPS per GB for the duplicate volume
// //duplicateTier: Tier level for the duplicate volume
// //duplicateSnapshotSize: Snapshot space size for the duplicate
// func (f fileStorageManager) OrderDuplicateVolume(originalVolumeId int, originalSnapshotId int, duplicateSize int, duplicateIops int, duplicateTier float64, duplicateSnapshotSize int) (datatypes.Container_Product_Order_Receipt, error) {
// 	fileMask := "id,billingItem.location,snapshotCapacityGb,storageType.keyName,capacityGb,originalVolumeSize,provisionedIops,storageTierLevel"
// 	originalVolume, err := f.GetFileVolumeDetails(originalVolumeId, fileMask)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	productPackage, err := GetPackage(f.PackageService, SaaS_Category)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	order, err := PrepareDuplicateOrderObject(productPackage, originalVolume, duplicateIops, duplicateTier, duplicateSize, duplicateSnapshotSize, "file")
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	if originalSnapshotId != 0 {
// 		order.DuplicateOriginSnapshotId = sl.Int(originalSnapshotId)
// 	}
// 	return f.OrderService.PlaceOrder(&order, sl.Bool(false))
// }

// //Deletes the specified snapshot object.
// //snapshotId: The ID of the snapshot object to delete.
// func (f fileStorageManager) DeleteSnapshot(snapshotId int) error {
// 	_, err := f.StorageService.Id(snapshotId).DeleteObject()
// 	return err
// }

// //Places an order for a file volume
// //storage_type: performance or endurance
// //location: Datacenter in which to order file volume
// //size: Size of the desired volume, in GB
// //iops: Number of IOPs for a "Performance" order
// //tier: Tier level to use for an "Endurance" order
// //snapshotSize: The size of optional snapshot space,
// func (f fileStorageManager) OrderFileVolume(storageType string, location string, size int, tier float64, iops int, snapshotSize int) (datatypes.Container_Product_Order_Receipt, error) {
// 	locationId, err := GetLocationId(f.LocationService, location)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid datacenter name specified. Please provide the lower case short name (e.g.: dal09)."))
// 	}
// 	productPackage, err := GetPackage(f.PackageService, SaaS_Category)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	prices, err := FindVolumePrices(productPackage, "file", storageType, size, tier, iops, snapshotSize)
// 	order := datatypes.Container_Product_Order_Network_Storage_AsAService{
// 		Container_Product_Order: datatypes.Container_Product_Order{
// 			ComplexType: sl.String(Saas_Order),
// 			PackageId:   productPackage.Id,
// 			Location:    sl.String(strconv.Itoa(locationId)),
// 			Quantity:    sl.Int(1),
// 			Prices:      prices,
// 		},
// 		VolumeSize: sl.Int(size),
// 	}
// 	if storageType == "performance" {
// 		order.Iops = sl.Int(iops)
// 	}
// 	return f.OrderService.PlaceOrder(&order, sl.Bool(false))
// }

// func (f fileStorageManager) CreateSnapshot(volumeId int, notes string) (datatypes.Network_Storage, error) {
// 	return f.StorageService.Id(volumeId).CreateSnapshot(sl.String((notes)))
// }

// //Enables snapshots for a specific block volume at a given schedule.
// //volumeId: The id of the volume
// //scheduleType: 'HOURLY'|'DAILY'|'WEEKLY'
// //retentionCount: Number of snapshots to be kept
// //minute: Minute when to take snapshot
// //hour: Hour when to take snapshot
// //dayOfWeek: Day when to take snapshot
// func (f fileStorageManager) EnableSnapshot(volumeId int, scheduleType string, retentionCount int, minute int, hour int, dayOfWeek string) error {
// 	_, err := f.StorageService.Id(volumeId).EnableSnapshots(sl.String(scheduleType), sl.Int(retentionCount), sl.Int(minute), sl.Int(hour), sl.String(dayOfWeek))
// 	return err
// }

// //Disables snapshots for a specific block volume at a given schedule.
// //volumeId: The id of the volume
// //scheduleType: 'HOURLY'|'DAILY'|'WEEKLY'
// func (f fileStorageManager) DisableSnapshots(volumeId int, scheduleType string) error {
// 	_, err := f.StorageService.Id(volumeId).DisableSnapshots(sl.String(scheduleType))
// 	return err
// }

// func (f fileStorageManager) OrderSaasSnapshotSpace(volumeId int, snapshotSize int, tier float64, iops int, upgrade bool) (datatypes.Container_Product_Order_Receipt, error) {
// 	productPackage, err := GetPackage(f.PackageService, SaaS_Category)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	fileMask := "serviceResource.datacenter.id,billingItem,storageTierLevel,storageType.keyName,iops"
// 	fileVolume, err := f.GetFileVolumeDetails(volumeId, fileMask)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}
// 	if fileVolume.BillingItem == nil || fileVolume.BillingItem.CategoryCode == nil {
// 		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Unable to find category code from this volume."))
// 	}
// 	storageType := strings.Split(*fileVolume.StorageType.KeyName, "_")[0]
// 	if storageType == "ENDURANCE" {
// 		if tier == 0 {
// 			tier, err = FindEnduranceTierIOPSPerGB(fileVolume)
// 			if err != nil {
// 				return datatypes.Container_Product_Order_Receipt{}, err
// 			}
// 		}
// 		iops = 0
// 	} else if storageType == "PERFORMANCE" {
// 		if iops == 0 {
// 			iops, err = strconv.Atoi(*fileVolume.Iops)
// 			if err != nil {
// 				return datatypes.Container_Product_Order_Receipt{}, err
// 			}
// 		}
// 		tier = 0
// 	} else {
// 		return datatypes.Container_Product_Order_Receipt{}, errors.New(T("Invalid storage type"))
// 	}
// 	spacePrice, err := FindSaasSnapshotSpacePrice(productPackage, snapshotSize, tier, iops)
// 	if err != nil {
// 		return datatypes.Container_Product_Order_Receipt{}, err
// 	}

// 	if upgrade {
// 		order := datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade{
// 			Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace: datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace{
// 				VolumeId: sl.Int(volumeId),
// 				Container_Product_Order: datatypes.Container_Product_Order{
// 					ComplexType: sl.String(("SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace_Upgrade")),
// 					PackageId:   productPackage.Id,
// 					Prices:      []datatypes.Product_Item_Price{spacePrice},
// 					Quantity:    sl.Int(1),
// 					Location:    sl.String(strconv.Itoa(*fileVolume.ServiceResource.Datacenter.Id)),
// 				},
// 			},
// 		}
// 		return f.OrderService.PlaceOrder(&order, sl.Bool(false))
// 	}
// 	order := datatypes.Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace{
// 		VolumeId: sl.Int(volumeId),
// 		Container_Product_Order: datatypes.Container_Product_Order{
// 			ComplexType: sl.String(("SoftLayer_Container_Product_Order_Network_Storage_Enterprise_SnapshotSpace")),
// 			PackageId:   productPackage.Id,
// 			Prices:      []datatypes.Product_Item_Price{spacePrice},
// 			Quantity:    sl.Int(1),
// 			Location:    sl.String(strconv.Itoa(*fileVolume.ServiceResource.Datacenter.Id)),
// 		},
// 	}
// 	return f.OrderService.PlaceOrder(&order, sl.Bool(false))
// }

// //Cancels snapshot space for a given volume.
// //volumeId: The volume ID
// //reason: The reason for cancellation
// //immediate: Cancel immediately or on anniversary date
// func (f fileStorageManager) CancelSnapshotSpace(volumeId int, reason string, immediate bool) error {
// 	fileVolume, err := f.GetFileVolumeDetails(volumeId, "id,billingItem.activeChildren")
// 	if err != nil {
// 		return err
// 	}
// 	if fileVolume.BillingItem == nil || len(fileVolume.BillingItem.ActiveChildren) == 0 {
// 		return errors.New(T("No snapshot space found to cancel."))
// 	}
// 	children := fileVolume.BillingItem.ActiveChildren
// 	billingItemId := 0
// 	for _, child := range children {
// 		if *child.CategoryCode == Snapshot_Category {
// 			billingItemId = *child.Id
// 			break
// 		}
// 	}
// 	if billingItemId == 0 {
// 		return errors.New(T("No snapshot space found to cancel."))
// 	}
// 	_, err = f.BillingService.Id(billingItemId).CancelItem(sl.Bool(immediate), sl.Bool(true), sl.String(reason), sl.String(""))
// 	return err
// }
