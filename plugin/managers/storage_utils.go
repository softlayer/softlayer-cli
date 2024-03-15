package managers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func HasCategory(categories []datatypes.Product_Item_Category, categoryCode string) bool {
	for _, category := range categories {
		if category.CategoryCode != nil && *category.CategoryCode == categoryCode {
			return true
		}
	}
	return false
}

// Find the snapshot schedule ID for the given volume and keyname
// volume: The volume for which the snapshot ID is desired
// scheduleKeyName: The keyname of the snapshot schedule
func FindSnapshotScheduleId(volume datatypes.Network_Storage, scheduleKeyName string) (int, error) {
	for _, s := range volume.Schedules {
		if s.Type != nil && s.Type.Keyname != nil {
			if *s.Type.Keyname == scheduleKeyName {
				return *s.Id, nil
			}
		}
	}
	return 0, errors.New(T("The given snapshot schedule name was not found for the given storage volume."))
}

// Populate the given host_templates array with the IDs provided
// hardwareIds: A List of SoftLayer_Hardware ids
// vsIds A List of SoftLayer_Virtual_Guest ids
// IPIds: A List of SoftLayer_Network_Subnet_IpAddress ids
// subnetIds: A List of SoftLayer_Network_Subnet ids
func PopulateHostTemplates(hardwareIds []int, vsIds []int, IPIds []int, subnetIds []int) []datatypes.Container_Network_Storage_Host {
	templates := []datatypes.Container_Network_Storage_Host{}
	for _, hardwareId := range hardwareIds {
		t := datatypes.Container_Network_Storage_Host{
			ObjectType: sl.String("SoftLayer_Hardware"),
			Id:         sl.Int(hardwareId),
		}
		templates = append(templates, t)
	}
	for _, vsId := range vsIds {
		t := datatypes.Container_Network_Storage_Host{
			ObjectType: sl.String("SoftLayer_Virtual_Guest"),
			Id:         sl.Int(vsId),
		}
		templates = append(templates, t)
	}
	for _, IPId := range IPIds {
		t := datatypes.Container_Network_Storage_Host{
			ObjectType: sl.String("SoftLayer_Network_Subnet_IpAddress"),
			Id:         sl.Int(IPId),
		}
		templates = append(templates, t)
	}
	for _, subnetId := range subnetIds {
		t := datatypes.Container_Network_Storage_Host{
			ObjectType: sl.String("SoftLayer_Network_Subnet"),
			Id:         sl.Int(subnetId),
		}
		templates = append(templates, t)
	}
	return templates
}

func ValidateDuplicateSize(originVolume datatypes.Network_Storage, duplicateSize int, volumeType string) (int, error) {
	if originVolume.CapacityGb == nil {
		return 0, errors.New(T("Cannot find origin volume's size"))
	}
	if duplicateSize == 0 && originVolume.CapacityGb != nil {
		duplicateSize = *originVolume.CapacityGb
	}
	if originVolume.CapacityGb != nil && duplicateSize < *originVolume.CapacityGb {
		return 0, errors.New(T("The requested duplicate volume size is too small. Duplicate volumes must be at least as large as their origin volumes."))
	}
	//if volumeType == VOLUME_TYPE_BLOCK {
	var baseVolumeSize int
	var err error
	if originVolume.OriginalVolumeSize != nil {
		baseVolumeSize, err = strconv.Atoi(*originVolume.OriginalVolumeSize)
		if err != nil {
			return 0, err
		}
	} else if originVolume.CapacityGb != nil {
		baseVolumeSize = *originVolume.CapacityGb
	}
	if duplicateSize > baseVolumeSize*10 {
		return 0, errors.New(T("The requested duplicate volume size is too large. The maximum size for duplicate block volumes is 10 times the size of the origin volume or, if the origin volume was also a duplicate, 10 times the size of the initial origin volume (i.e. the origin volume from which the first duplicate was created in the chain of duplicates). Requested: {{.DuplicateSize}} GB. Base origin size: {{.BaseSize}} GB.", map[string]interface{}{"DuplicateSize": duplicateSize, "BaseSize": baseVolumeSize}))
	}
	//}
	return duplicateSize, nil
}

func ValidateDuplicatePerformanceIops(originalVolume datatypes.Network_Storage, duplicateIops int, duplicateSize int) (int, error) {
	if originalVolume.ProvisionedIops == nil {
		return 0, errors.New(T("Cannot find origin volume's provisioned IOPS"))
	}
	if duplicateIops == 0 {
		var err error
		duplicateIops, err = strconv.Atoi(*originalVolume.ProvisionedIops)
		if err != nil {
			return 0, err
		}
	} else {
		originalProvisionedIops, err := strconv.Atoi(*originalVolume.ProvisionedIops)
		if err != nil {
			return 0, err
		}
		if originalVolume.CapacityGb != nil && *originalVolume.CapacityGb != 0 {
			originIopsPerGb := float32(originalProvisionedIops) / float32(*originalVolume.CapacityGb)
			duplicateIopsPerGb := float32(duplicateIops) / float32(duplicateSize)
			if originIopsPerGb < 0.3 && duplicateIopsPerGb > 0.3 {
				return 0, errors.New(T("Origin volume performance is < 0.3 IOPS/GB, duplicate volume performance must also be < 0.3 IOPS/GB. {{.DuplicateIopsPerGb}} IOPS/GB ({{.DuplicateIops}}/{{.DuplicateSize}}) requested.", map[string]interface{}{"DuplicateIopsPerGb": duplicateIopsPerGb, "DuplicateIops": duplicateIops, "DuplicateSize": duplicateSize}))
			} else if originIopsPerGb >= 0.3 && duplicateIopsPerGb < 0.3 {
				return 0, errors.New(T("Origin volume performance is >= 0.3 IOPS/GB, duplicate volume performance must also be >= 0.3 IOPS/GB. {{.DuplicateIopsPerGb}} IOPS/GB ({{.DuplicateIops}}/{{.DuplicateSize}}) requested.", map[string]interface{}{"DuplicateIopsPerGb": duplicateIopsPerGb, "DuplicateIops": duplicateIops, "DuplicateSize": duplicateSize}))
			}
		}
	}
	return duplicateIops, nil
}

func ValidateDuplicateEnduranceTier(originalVolume datatypes.Network_Storage, duplicateTier float64) (float64, error) {
	originalTier, err := FindEnduranceTierIOPSPerGB(originalVolume)
	if err != nil {
		return 0, errors.New(T("Cannot find original volume's tier level"))
	}
	if duplicateTier == 0 {
		duplicateTier = originalTier
	} else {
		if originalTier == 0.25 && duplicateTier != 0.25 {
			return 0, errors.New(T("Origin volume performance tier is 0.25 IOPS/GB, duplicate volume performance tier must also be 0.25 IOPS/GB. {{.DuplicateTier}} IOPS/GB requested.", map[string]interface{}{"DuplicateTier": duplicateTier}))
		}
		if originalTier != 0.25 && duplicateTier == 0.25 {
			return 0, errors.New(T("Origin volume performance tier is above 0.25 IOPS/GB, duplicate volume performance tier must also be above 0.25 IOPS/GB. {{.DuplicateTier}} IOPS/GB requested.", map[string]interface{}{"DuplicateTier": duplicateTier}))
		}
	}
	return duplicateTier, nil
}

// Find the tier for the given endurance volume (IOPS per GB)
// volume: The volume for which the tier level is desired
func FindEnduranceTierIOPSPerGB(volume datatypes.Network_Storage) (float64, error) {
	if volume.StorageTierLevel != nil {
		tier := TIER_PER_IOPS[*volume.StorageTierLevel]
		if tier > 0 {
			return tier, nil
		}
	}
	return 0.0, errors.New(T("Could not find tier IOPS per GB for this volume."))
}

// Find a price in the SaaS package with the specified category
// productPackage: The product package (performance_storage_iscsi,storage_service_enterprise,storage_endurance,storage_as_a_service,storage_block,storage_file )
// priceCategory: The price category to search for
func FindPriceByCategory(productPackage datatypes.Product_Package, priceCategory string) (datatypes.Product_Item_Price, error) {
	for _, item := range productPackage.Items {
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, priceCategory) {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price with the category: {{.PriceCategory}}", map[string]interface{}{"PriceCategory": priceCategory}))
}

// Find the SaaS endurance storage space price for the size and tier
// productPackage: The Storage As A Service product package
// size: The volume size for which a price is desired
// tier: The endurance tier for which a price is desired
func FindSaasEnduranceSpacePrice(productPackage datatypes.Product_Package, size int, tier float64) (datatypes.Product_Item_Price, error) {
	var keyName string
	if tier == 0.25 {
		keyName = fmt.Sprintf("STORAGE_SPACE_FOR_%.2f_IOPS_PER_GB", tier)
		keyName = strings.Replace(keyName, ".", "_", -1)
	} else {
		keyName = fmt.Sprintf("STORAGE_SPACE_FOR_%d_IOPS_PER_GB", int(tier))
	}
	for _, item := range productPackage.Items {
		if item.KeyName == nil || !strings.Contains(*item.KeyName, keyName) {
			continue
		}
		if item.CapacityMaximum == nil || item.CapacityMinimum == nil {
			continue
		}
		maxCapacity, err := strconv.Atoi(*item.CapacityMaximum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		minCapacity, err := strconv.Atoi(*item.CapacityMinimum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		if size < minCapacity || size > maxCapacity {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Space_Category) {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for endurance storage space, size={{.Size}} tier={{.Tier}}", map[string]interface{}{"Size": size, "Tier": tier}))
}

// Find the SaaS storage tier level price for the specified tier level
// productPackage: The Storage As A Service product package
// tier: The endurance tier for which a price is desired
func FindSaasEnduranceTierPrice(productPacakge datatypes.Product_Package, tier float64) (datatypes.Product_Item_Price, error) {
	targetCapacity := ENDURANCE_TIERS[tier]
	for _, item := range productPacakge.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != Tier_Category {
			continue
		}
		if item.Capacity == nil || int(*item.Capacity) != targetCapacity {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Tier_Category) {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for endurance tier level, tier={{.Tier}}", map[string]interface{}{"Tier": tier}))
}

// Find the SaaS performance storage space price for the given size
// productPackage: The Storage As A Service product package
// size: The volume size for which a price is desired
func FindSaasPerformanceSpacePrice(productPacakge datatypes.Product_Package, size int) (datatypes.Product_Item_Price, error) {
	for _, item := range productPacakge.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != Space_Category {
			continue
		}
		if item.CapacityMaximum == nil || item.CapacityMinimum == nil {
			continue
		}
		maxCapacity, err := strconv.Atoi(*item.CapacityMaximum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		minCapacity, err := strconv.Atoi(*item.CapacityMinimum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		if size < minCapacity || size > maxCapacity {
			continue
		}
		keyName := fmt.Sprintf("%d_%d_GBS", minCapacity, maxCapacity)
		if item.KeyName == nil || !strings.Contains(*item.KeyName, keyName) {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Space_Category) {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for performance storage space, size={{.Size}}", map[string]interface{}{"Size": size}))
}

// Find the SaaS IOPS price for the specified size and iops
// productPackage: The Storage As A Service product package
// size: The volume size for which a price is desired
// iops: The number of IOPS for which a price is desired
func FindSaasPerformanceIopsPrice(productPacakge datatypes.Product_Package, size int, iops int) (datatypes.Product_Item_Price, error) {
	for _, item := range productPacakge.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != Iops_Category {
			continue
		}
		if item.CapacityMaximum == nil || item.CapacityMinimum == nil {
			continue
		}
		maxCapacity, err := strconv.Atoi(*item.CapacityMaximum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		minCapacity, err := strconv.Atoi(*item.CapacityMinimum)
		if err != nil {
			return datatypes.Product_Item_Price{}, err
		}
		if iops < minCapacity || iops > maxCapacity {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Iops_Category) {
				continue
			}
			restrictMinCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMinimum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			restrictMaxCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMaximum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			if price.CapacityRestrictionType == nil ||
				*price.CapacityRestrictionType != "STORAGE_SPACE" ||
				size < restrictMinCapacity ||
				size > restrictMaxCapacity {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for iops for the given volume, size={{.Size}},,Iops={{.IOPS}}", map[string]interface{}{"Size": size, "IOPS": iops}))
}

// Find the price in the SaaS package for the desired snapshot space size
// productPackage: The Storage As A Service product package
// size: The volume size for which a price is desired
// tier: The tier of the volume for which space is being ordered
// iops: The number of IOPS for which a price is desired
func FindSaasSnapshotSpacePrice(productPacakge datatypes.Product_Package, size int, tier float64, iops int) (datatypes.Product_Item_Price, error) {
	var targetValue int
	var targetRestrictionType string
	if tier != 0 && iops == 0 {
		targetValue = ENDURANCE_TIERS[tier]
		targetRestrictionType = "STORAGE_TIER_LEVEL"
	} else if tier == 0 && iops != 0 {
		targetValue = iops
		targetRestrictionType = "IOPS"
	} else {
		return datatypes.Product_Item_Price{}, errors.New(T("Specify either tier or iops, unable to specify both"))
	}
	for _, item := range productPacakge.Items {
		if item.Capacity == nil || int(*item.Capacity) != size {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Snapshot_Category) {
				continue
			}
			restrictMinCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMinimum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			restrictMaxCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMaximum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			if price.CapacityRestrictionType == nil ||
				targetRestrictionType != *price.CapacityRestrictionType ||
				targetValue < restrictMinCapacity ||
				targetValue > restrictMaxCapacity {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for snapshot space,size={{.Size}},tier={{.Tier}},Iops={{.IOPS}}", map[string]interface{}{"Size": size, "Tier": tier, "IOPS": iops}))
}

// Find the price in the SaaS package for the desired replication  price
// productPackage: The Storage As A Service product package
// tier: The tier of the volume for which space is being ordered
// iops: The number of IOPS for which a price is desired
func FindSaasReplicationPrice(productPacakge datatypes.Product_Package, tier float64, iops int) (datatypes.Product_Item_Price, error) {
	var targetValue int
	var targetRestrictionType string
	if tier != 0 && iops == 0 {
		targetValue = ENDURANCE_TIERS[tier]
		targetRestrictionType = "STORAGE_TIER_LEVEL"
	} else if tier == 0 && iops != 0 {
		targetValue = iops
		targetRestrictionType = "IOPS"
	} else {
		return datatypes.Product_Item_Price{}, errors.New(T("Specify either tier or iops, unable to specify both"))
	}
	for _, item := range productPacakge.Items {
		if item.ItemCategory == nil || item.ItemCategory.CategoryCode == nil || *item.ItemCategory.CategoryCode != Replication_Category {
			continue
		}
		for _, price := range item.Prices {
			if price.LocationGroupId != nil {
				continue
			}
			if !HasCategory(price.Categories, Replication_Category) {
				continue
			}
			restrictMinCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMinimum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			restrictMaxCapacity, err := strconv.Atoi(utils.StringPointertoString(price.CapacityRestrictionMaximum))
			if err != nil {
				return datatypes.Product_Item_Price{}, err
			}
			if price.CapacityRestrictionType == nil ||
				targetRestrictionType != *price.CapacityRestrictionType ||
				targetValue < restrictMinCapacity ||
				targetValue > restrictMaxCapacity {
				continue
			}
			return price, nil
		}
	}
	return datatypes.Product_Item_Price{}, errors.New(T("Could not find price for replication,tier={{.Tier}},Iops={{.IOPS}}", map[string]interface{}{"Tier": tier, "IOPS": iops}))
}

func PrepareSaaSReplicantOrderObject(productPackage datatypes.Product_Package, snapshotSchedule string, locationId int, tier float64, iops int, volume datatypes.Network_Storage, volumeType string) (datatypes.Container_Product_Order_Network_Storage_AsAService, error) {
	var err error
	var volumeId, size int
	if volume.Id != nil {
		volumeId = *volume.Id
	}
	if volume.CapacityGb != nil {
		size = *volume.CapacityGb
	}
	volumeSnapshotCapacity := 0
	if volume.SnapshotCapacityGb != nil {
		volumeSnapshotCapacity, err = strconv.Atoi(*volume.SnapshotCapacityGb)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
	} else {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Snapshot capacity not found for the given primary volume."))
	}

	snapshotScheduleId, err := FindSnapshotScheduleId(volume, "SNAPSHOT_"+snapshotSchedule)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}

	if volume.BillingItem != nil && volume.BillingItem.CancellationDate != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("This volume is set for cancellation; unable to order replicant volume."))
	}

	if volume.BillingItem != nil && volume.BillingItem.ActiveChildren != nil {
		for _, child := range volume.BillingItem.ActiveChildren {
			if child.CategoryCode != nil && *child.CategoryCode == Snapshot_Category && child.CancellationDate != nil {
				return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("The snapshot space for this volume is set for cancellation; unable to order replicant volume."))
			}
		}
	}
	var storageType string
	if volume.StorageType != nil && volume.StorageType.KeyName != nil {
		storageType = strings.Split(*volume.StorageType.KeyName, "_")[0]
	}
	if storageType == strings.ToUpper(STORAGE_TYPE_ENDURANCE) {
		if tier == 0 {
			tier, err = FindEnduranceTierIOPSPerGB(volume)
			if err != nil {
				return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
			}
		}
		iops = 0
	} else if storageType == strings.ToUpper(STORAGE_TYPE_PERFORMANCE) {
		if iops == 0 && volume.Iops != nil {
			iops, err = strconv.Atoi(*volume.Iops)
			if err != nil {
				return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
			}
		}
		tier = 0
	} else {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Invalid storage type"))
	}
	var prices []datatypes.Product_Item_Price
	servicePrice, err := FindPriceByCategory(productPackage, SaaS_Category)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}
	volumePrice, err := FindPriceByCategory(productPackage, "storage_"+volumeType)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}
	prices = append(prices, servicePrice)
	prices = append(prices, volumePrice)
	// PERFORMANCE Prices
	if tier == 0 && iops != 0 {
		spacePrice, err := FindSaasPerformanceSpacePrice(productPackage, size)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		iopsPrice, err := FindSaasPerformanceIopsPrice(productPackage, size, iops)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		prices = append(prices, spacePrice)
		prices = append(prices, iopsPrice)
		// ENDURANCE Prices
	} else if tier != 0 && iops == 0 {
		spacePrice, err := FindSaasEnduranceSpacePrice(productPackage, size, tier)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		tierPrice, err := FindSaasEnduranceTierPrice(productPackage, tier)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		prices = append(prices, spacePrice)
		prices = append(prices, tierPrice)
		// BAD INPUT
	} else {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Specify either iops or tier, cannot specify both."))
	}
	snapshotPrice, err := FindSaasSnapshotSpacePrice(productPackage, volumeSnapshotCapacity, tier, iops)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}
	replicantPrice, err := FindSaasReplicationPrice(productPackage, tier, iops)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}
	prices = append(prices, snapshotPrice)
	prices = append(prices, replicantPrice)

	replicantOrder := datatypes.Container_Product_Order_Network_Storage_AsAService{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String(Saas_Order),
			PackageId:   productPackage.Id,
			Location:    sl.String(strconv.Itoa(locationId)),
			Quantity:    sl.Int(1),
			Prices:      prices,
		},
		OriginVolumeId:         sl.Int(volumeId),
		OriginVolumeScheduleId: sl.Int(snapshotScheduleId),
		VolumeSize:             sl.Int(size),
	}
	if storageType == strings.ToUpper(STORAGE_TYPE_PERFORMANCE) {
		replicantOrder.Iops = sl.Int(iops)
	}
	return replicantOrder, nil
}

// Returns a product packaged based on type of storage.
// categoryCode: Category code of product package.
func GetPackage(packageService services.Product_Package, categoryCode string) (datatypes.Product_Package, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("categories.categoryCode").Eq(categoryCode))
	filters = append(filters, filter.Path("statusCode").Eq("ACTIVE"))
	packages, err := packageService.Mask("id,name,items[prices[categories],attributes]").Filter(filters.Build()).GetAllObjects()
	if err != nil {
		return datatypes.Product_Package{}, err
	}
	if len(packages) == 0 {
		return datatypes.Product_Package{}, errors.New(T("No packages were found for {{.CategoryCode}}.", map[string]interface{}{"CategoryCode": categoryCode}))
	}
	if len(packages) > 1 {
		return datatypes.Product_Package{}, errors.New(T("More than one packages were found for {{.CategoryCode}}.", map[string]interface{}{"CategoryCode": categoryCode}))
	}
	return packages[0], nil
}

// Returns location id of datacenter for ProductOrder::placeOrder().
// location: Datacenter short name
func GetLocationId(locationService services.Location_Datacenter, location string) (int, error) {
	filter := filter.New(filter.Path("name").Eq(location))
	datacenters, err := locationService.Mask("longName,id,name").Filter(filter.Build()).GetDatacenters()
	if err != nil {
		return 0, err
	}
	for _, datacenter := range datacenters {
		if datacenter.Name != nil && *datacenter.Name == location {
			return *datacenter.Id, nil
		}
	}
	return 0, errors.New(T("Invalid datacenter name specified."))
}

// Prepare the duplicate order to submit to SoftLayer_Product.PlaceOrder()
// originalVolume: The origin volume which is being duplicated
// iops: The IOPS per GB for the duplicant volume (performance)
// tier: The tier level for the duplicant volume (endurance)
// duplicateSize: The requested size for the duplicate volume
// duplicateSnapshotSize: The size for the duplicate snapshot space. -1 will use originalVolumes snapshotSize.
// volumeType: The type of the origin volume ('file' or 'block')
func PrepareDuplicateOrderObject(productPackage datatypes.Product_Package, originalVolume datatypes.Network_Storage, config DuplicateOrderConfig) (datatypes.Container_Product_Order_Network_Storage_AsAService, error) {
	// iops int, tier float64, duplicateSize int, duplicateSnapshotSize int, volumeType string
	iops := config.DuplicateIops
	tier := config.DuplicateTier
	duplicateSize := config.DuplicateSize
	duplicateSnapshotSize := config.DuplicateSnapshotSize
	volumeType := config.VolumeType
	//Verify that the origin volume has not been cancelled
	if originalVolume.BillingItem == nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("The original volume has been cancelled, unable to order duplicate volume"))
	}
	//Verify that the origin volume has snapshot space (needed for duplication)
	if originalVolume.SnapshotCapacityGb == nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Snapshot space not found for original volume, origin snapshot space is needed for duplication"))
	}
	originalSnapshotSize, err := strconv.Atoi(*originalVolume.SnapshotCapacityGb)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Snapshot space not found for original volume, origin snapshot space is needed for duplication"))
	}
	//Obtain the datacenter location ID for the duplicate
	if originalVolume.BillingItem.Location == nil || originalVolume.BillingItem.Location.Id == nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Cannot find original volume's location"))
	}
	locationId := *originalVolume.BillingItem.Location.Id
	if duplicateSnapshotSize == -1 {
		duplicateSnapshotSize = originalSnapshotSize
	}
	duplicateSize, err = ValidateDuplicateSize(originalVolume, duplicateSize, volumeType)
	if err != nil {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
	}

	var volumeIsPerformance bool
	var prices []datatypes.Product_Item_Price
	var originalStorageType string
	if originalVolume.StorageType != nil && originalVolume.StorageType.KeyName != nil {
		originalStorageType = strings.Split(*originalVolume.StorageType.KeyName, "_")[0]
	}
	if originalStorageType == strings.ToUpper(STORAGE_TYPE_PERFORMANCE) {
		volumeIsPerformance = true
		iops, err = ValidateDuplicatePerformanceIops(originalVolume, iops, duplicateSize)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		prices, err = FindVolumePrices(productPackage, volumeType, STORAGE_TYPE_PERFORMANCE, duplicateSize, 0, iops, duplicateSnapshotSize)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
	} else if originalStorageType == strings.ToUpper(STORAGE_TYPE_ENDURANCE) {
		volumeIsPerformance = false
		tier, err := ValidateDuplicateEnduranceTier(originalVolume, tier)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
		prices, err = FindVolumePrices(productPackage, volumeType, STORAGE_TYPE_ENDURANCE, duplicateSize, tier, 0, duplicateSnapshotSize)
		if err != nil {
			return datatypes.Container_Product_Order_Network_Storage_AsAService{}, err
		}
	} else {
		return datatypes.Container_Product_Order_Network_Storage_AsAService{}, errors.New(T("Origin volume does not have a valid storage type (with an appropriate keyName to indicate the volume is a PERFORMANCE or ENDURANCE volume)"))
	}
	duplicateOrder := datatypes.Container_Product_Order_Network_Storage_AsAService{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String(Saas_Order),
			PackageId:   productPackage.Id,
			Location:    sl.String(strconv.Itoa(locationId)),
			Quantity:    sl.Int(1),
			Prices:      prices,
		},
		VolumeSize:              sl.Int(duplicateSize),
		DuplicateOriginVolumeId: originalVolume.Id,
	}
	if volumeIsPerformance {
		duplicateOrder.Iops = sl.Int(iops)
	}
	return duplicateOrder, nil
}

// volumeType: block or file
// storageType: performance or endurance
func FindVolumePrices(productPackage datatypes.Product_Package, volumeType string, storageType string, size int, tier float64, iops int, snapshotSize int) ([]datatypes.Product_Item_Price, error) {
	var prices []datatypes.Product_Item_Price
	servicePrice, err := FindPriceByCategory(productPackage, SaaS_Category)
	if err != nil {
		return nil, err
	}
	prices = append(prices, servicePrice)
	volumeTypePrice, err := FindPriceByCategory(productPackage, "storage_"+volumeType)
	if err != nil {
		return nil, err
	}
	prices = append(prices, volumeTypePrice)
	if storageType == STORAGE_TYPE_PERFORMANCE {
		spacePrice, err := FindSaasPerformanceSpacePrice(productPackage, size)
		if err != nil {
			return nil, err
		}
		prices = append(prices, spacePrice)
		iopsPrice, err := FindSaasPerformanceIopsPrice(productPackage, size, iops)
		if err != nil {
			return nil, err
		}
		prices = append(prices, iopsPrice)
	} else if storageType == STORAGE_TYPE_ENDURANCE {
		spacePrice, err := FindSaasEnduranceSpacePrice(productPackage, size, tier)
		if err != nil {
			return nil, err
		}
		prices = append(prices, spacePrice)
		tierPrice, err := FindSaasEnduranceTierPrice(productPackage, tier)
		if err != nil {
			return nil, err
		}
		prices = append(prices, tierPrice)
	} else {
		return nil, errors.New(T("Invalid storage type {{.StorageType}}", map[string]interface{}{"StorageType": storageType}))
	}
	if snapshotSize > 0 {
		snapshotPrice, err := FindSaasSnapshotSpacePrice(productPackage, snapshotSize, tier, iops)
		if err != nil {
			return nil, err
		}
		prices = append(prices, snapshotPrice)
	}
	return prices, nil
}

// volumeType: block or file
// storageType: performance or endurance
func FindVolumePricesUpgrade(productPackage datatypes.Product_Package, volumeType string, storageType string, size int, tier float64, iops int) ([]datatypes.Product_Item_Price, error) {
	var prices []datatypes.Product_Item_Price
	servicePrice, err := FindPriceByCategory(productPackage, SaaS_Category)
	if err != nil {
		return nil, err
	}
	prices = append(prices, servicePrice)

	if storageType == STORAGE_TYPE_PERFORMANCE {
		spacePrice, err := FindSaasPerformanceSpacePrice(productPackage, size)
		if err != nil {
			return nil, err
		}
		prices = append(prices, spacePrice)
		iopsPrice, err := FindSaasPerformanceIopsPrice(productPackage, size, iops)
		if err != nil {
			return nil, err
		}
		prices = append(prices, iopsPrice)
	} else if storageType == STORAGE_TYPE_ENDURANCE {
		spacePrice, err := FindSaasEnduranceSpacePrice(productPackage, size, tier)
		if err != nil {
			return nil, err
		}
		prices = append(prices, spacePrice)
		tierPrice, err := FindSaasEnduranceTierPrice(productPackage, tier)
		if err != nil {
			return nil, err
		}
		prices = append(prices, tierPrice)
	} else {
		return nil, errors.New(T("Invalid storage type {{.StorageType}}", map[string]interface{}{"StorageType": storageType}))
	}
	return prices, nil
}
