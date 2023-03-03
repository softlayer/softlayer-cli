package block

import (
	"bytes"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

var (
	volumeSizes = []string{"20", "40", "80", "100", "250", "500", "1000", "2000-3000", "4000-7000", "8000-9000", "10000-12000"}
)

type VolumeOptionsCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	NetworkManager managers.NetworkManager
	Prices         bool
}

func NewVolumeOptionsCommand(sl *metadata.SoftlayerStorageCommand) *VolumeOptionsCommand {
	thisCmd := &VolumeOptionsCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
		NetworkManager:          managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-options",
		Short: T("List all options for ordering a block storage"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} volume-options
	
EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} volume-options
   This command lists all options for creating a block storage volume, including storage type, volume size, OS type, IOPS, tier level, datacenter, and snapshot size.`, sl.StorageI18n),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Prices, "prices", false, T("Show prices in the storage, snapshot and iops range tables."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeOptionsCommand) Run(args []string) error {
	PACKAGE_ID := 759
	CATEGORY_CODE_STORAGE_PAGKAGE := "performance_storage_space"
	CATEGORY_CODE_IOPS := "storage_tier_level"
	CATEGORY_CODE_IOPS_PAGKAGE := "performance_storage_iops"

	datacentersByPackage, err := cmd.StorageManager.GetRegions(PACKAGE_ID)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get all datacenters."), err.Error(), 2)
	}

	storagePackages, err := cmd.StorageManager.ListItems(PACKAGE_ID, CATEGORY_CODE_STORAGE_PAGKAGE, "")
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	osTypes, err := cmd.StorageManager.GetOsType()
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	maskIops := "mask[id,keyName,description,locationConflicts]"
	iopsLevels, err := cmd.StorageManager.ListItems(PACKAGE_ID, CATEGORY_CODE_IOPS, maskIops)
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	iopsPackages, err := cmd.StorageManager.ListItems(PACKAGE_ID, CATEGORY_CODE_IOPS_PAGKAGE, "")
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	// Get datacenters to Iops locations conflicts
	datacenters, err := cmd.StorageManager.GetAllDatacenters()
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	// Get closing pods
	pods, err := cmd.NetworkManager.GetClosingPods()
	if err != nil {
		return slErr.NewAPIError(T("API Error."), err.Error(), 2)
	}

	// Tables
	tableDatacenter := cmd.UI.Table([]string{T("Datacenter"), T("Description"), T("Name"), T("Note")})
	tableStorage := cmd.UI.Table([]string{T("Storage"), T("Description"), T("KeyName")})
	tableSnapshot := cmd.UI.Table([]string{T("Snapshot"), T("Description"), T("KeyName")})
	tableIopsRange := cmd.UI.Table([]string{T("Storage"), T("Range")})
	tableSnapshotRules := cmd.UI.Table([]string{T("Rules Snapshot - Storage Size (GB)"), T("Available Snapshot Size (GB)")})
	tableOsTypes := cmd.UI.Table([]string{T("OS Type"), T("KeyName"), T("Description")})
	tableIops := cmd.UI.Table([]string{T("IOPS"), T("KeyName"), T("Description"), T("Location Conflicts")})

	if cmd.Prices {
		tableStorage = cmd.UI.Table([]string{T("Storage"), T("Description"), T("KeyName"), T("Prices")})
		tableSnapshot = cmd.UI.Table([]string{T("Snapshot"), T("Description"), T("KeyName"), T("Prices")})
		tableIopsRange = cmd.UI.Table([]string{T("Storage"), T("Range"), T("Prices")})
	}

	// Datacenter table
	for _, datacenter := range datacentersByPackage {
		note := utils.GetPodWithClosedAnnouncement(*datacenter.Location.Location.LongName, pods)
		tableDatacenter.Add(
			utils.FormatIntPointer(datacenter.Location.LocationId),
			utils.FormatStringPointer(datacenter.Description),
			utils.FormatStringPointer(datacenter.Location.Location.Name),
			note,
		)
	}
	tableDatacenter.Print()
	println()

	for _, storage := range storagePackages {
		// Adding datas to storage table
		if *storage.Prices[0].CapacityRestrictionType == "IOPS" {
			if cmd.Prices {
				tableStorage.Add(
					utils.FormatIntPointer(storage.Id),
					utils.FormatStringPointer(storage.Description),
					utils.FormatStringPointer(storage.KeyName),
					getPrices(storage.Prices, false),
				)
			} else {
				tableStorage.Add(
					utils.FormatIntPointer(storage.Id),
					utils.FormatStringPointer(storage.Description),
					utils.FormatStringPointer(storage.KeyName),
				)
			}
		}
		// Adding datas to snapshot table
		if strings.Contains(*storage.KeyName, "_STORAGE_SPACE") {
			if cmd.Prices {
				tableSnapshot.Add(
					utils.FormatIntPointer(storage.Id),
					utils.FormatStringPointer(storage.Description),
					utils.FormatStringPointer(storage.KeyName),
					getPrices(storage.Prices, true),
				)
			} else {
				tableSnapshot.Add(
					utils.FormatIntPointer(storage.Id),
					utils.FormatStringPointer(storage.Description),
					utils.FormatStringPointer(storage.KeyName),
				)
			}
		}
	}
	tableStorage.Print()
	println()

	tableSnapshot.Print()
	println()

	// Snapshot rules table
	tableSnapshotRules.Add(volumeSizes[0], "0,5,10,20")
	tableSnapshotRules.Add(volumeSizes[1], "0,5,10,20,40")
	tableSnapshotRules.Add(volumeSizes[2], "0,5,10,20,40,60,80")
	tableSnapshotRules.Add(volumeSizes[3], "0,5,10,20,40,60,80,100")
	tableSnapshotRules.Add(volumeSizes[4], "0,5,10,20,40,60,80,100,150,200,250")
	tableSnapshotRules.Add(volumeSizes[5], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500")
	tableSnapshotRules.Add(volumeSizes[6], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000")
	tableSnapshotRules.Add(volumeSizes[7], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000")
	tableSnapshotRules.Add(volumeSizes[8], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	tableSnapshotRules.Add(volumeSizes[9], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	tableSnapshotRules.Add(volumeSizes[10], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	tableSnapshotRules.Print()
	println()

	// Os type table
	for _, osType := range osTypes {
		tableOsTypes.Add(
			utils.FormatStringPointer(osType.Name),
			utils.FormatStringPointer(osType.KeyName),
			utils.FormatStringPointer(osType.Description),
		)
	}
	tableOsTypes.Print()
	println()

	// Iops tier table
	for _, iopsLevel := range iopsLevels {
		locationConflicts := ""
		if len(iopsLevel.LocationConflicts) != 0 {
			locationConflicts = getRegions(iopsLevel.LocationConflicts, datacenters)
		}
		tableIops.Add(
			utils.FormatIntPointer(iopsLevel.Id),
			utils.FormatStringPointer(iopsLevel.KeyName),
			utils.FormatStringPointer(iopsLevel.Description),
			locationConflicts,
		)
	}
	tableIops.Print()
	println()

	// Iops range table
	println("The storage size affects the selectable IOPS range. View the table below.")
	for _, iops := range iopsPackages {
		if cmd.Prices {
			tableIopsRange.Add(
				utils.FormatStringPointer(iops.Prices[0].CapacityRestrictionMinimum)+" - "+
					utils.FormatStringPointer(iops.Prices[0].CapacityRestrictionMaximum)+" GBs",
				utils.FormatStringPointer(iops.Description),
				getPrices(iops.Prices, false),
			)
		} else {
			tableIopsRange.Add(
				utils.FormatStringPointer(iops.Prices[0].CapacityRestrictionMinimum)+" - "+
					utils.FormatStringPointer(iops.Prices[0].CapacityRestrictionMaximum)+" GBs",
				utils.FormatStringPointer(iops.Description),
			)
		}
	}
	tableIopsRange.Print()
	println(T("Note: IOPS above 6,000 available only in: https://cloud.ibm.com/docs/BlockStorage?topic=BlockStorage-selectDC"))
	return nil
}

func getRegions(regionConflicts []datatypes.Product_Item_Resource_Conflict, datacenters []datatypes.Location) string {
	listLocationConflicts := []string{}
	for _, regionConflict := range regionConflicts {
		for _, datacenter := range datacenters {
			if *regionConflict.ResourceTableId == *datacenter.Id {
				listLocationConflicts = append(listLocationConflicts, *datacenter.Name)
			}
		}
	}
	return strings.Join(listLocationConflicts, ",")
}

func getPrices(prices []datatypes.Product_Item_Price, tierLevel bool) string {
	buf := new(bytes.Buffer)
	tablePrices := terminal.NewTable(buf, []string{
		T("Id"),
		T("Hourly/Monthly"),
		T("Datacenters"),
	})
	if tierLevel {
		tablePrices = terminal.NewTable(buf, []string{
			T("Id"),
			T("Tier"),
			T("Hourly/Monthly"),
			T("Datacenters"),
		})
	}
	for _, price := range prices {
		datacenters := []string{}
		if price.PricingLocationGroup != nil {
			for _, location := range price.PricingLocationGroup.Locations {
				datacenters = append(datacenters, utils.FormatStringPointer(location.Name))
			}
		} else {
			datacenters = append(datacenters, "-")
		}
		if tierLevel {
			tier := "-"
			if *price.CapacityRestrictionType == "STORAGE_TIER_LEVEL" {
				switch *price.CapacityRestrictionMinimum {
				case "100":
					tier = "0.25"
				case "200":
					tier = "2"
				case "300":
					tier = "4"
				case "1000":
					tier = "10"
				}
			}
			tablePrices.Add(
				utils.FormatIntPointer(price.Id),
				tier,
				utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee)+"/"+utils.FormatSLFloatPointerToFloat(price.RecurringFee),
				strings.Join(datacenters, ","),
			)
		} else {
			tablePrices.Add(
				utils.FormatIntPointer(price.Id),
				utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee)+"/"+utils.FormatSLFloatPointerToFloat(price.RecurringFee),
				strings.Join(datacenters, ","),
			)
		}
	}
	tablePrices.Print()
	return buf.String()
}
