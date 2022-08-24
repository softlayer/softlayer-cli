package block

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeOrderCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	StorageType    string
	Size           int
	Iops           int
	Tier           float64
	OsType         string
	Datacenter     string
	SnapshotSize   int
	Billing        string
	Force          bool
}

func NewVolumeOrderCommand(sl *metadata.SoftlayerStorageCommand) *VolumeOrderCommand {
	thisCmd := &VolumeOrderCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-order",
		Short: T("Order a block storage volume"),
		Long: T(`${COMMAND_NAME} sl block volume-order [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block volume-order --storage-type performance --size 1000 --iops 4000 --os-type LINUX -d dal09
   This command orders a performance volume with size is 1000GB, IOPS is 4000, OS type is LINUX, located at dal09.
   ${COMMAND_NAME} sl block volume-order --storage-type endurance --size 500 --tier 4 --os-type XEN -d dal09 --snapshot-size 500
   This command orders a endurance volume with size is 500GB, tier level is 4 IOPS per GB, OS type is XEN, located at dal09, and additional snapshot space size is 500GB.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.StorageType, "storage-type", "t", "", T("Type of storage volume [required], options are: performance,endurance"))
	cobraCmd.Flags().IntVarP(&thisCmd.Size, "size", "s", 0, T("Size of storage volume in GB [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Iops, "iops", "i", 0, T("Performance Storage IOPs, between 100 and 6000 in multiples of 100 [required for storage-type performance]"))
	cobraCmd.Flags().Float64VarP(&thisCmd.Tier, "tier", "e", 0, T("Endurance Storage Tier (IOP per GB) [required for storage-type endurance], options are: 0.25,2,4,10"))
	cobraCmd.Flags().StringVarP(&thisCmd.OsType, "os-type", "o", "", T("Operating System [required], options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Datacenter short name [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.SnapshotSize, "snapshot-size", "n", 0, T("Optional parameter for ordering snapshot space along with endurance block storage; specifies the size (in GB) of snapshot space to order"))
	cobraCmd.Flags().StringVarP(&thisCmd.Billing, "billing", "b", "", T("Optional parameter for Billing rate (default to monthly), options are: hourly, monthly"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeOrderCommand) Run(args []string) error {
	subs := map[string]interface{}{"CommandName": "ibmcloud"}
	if cmd.StorageType == "" {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}
	storageType := cmd.StorageType
	if storageType != "performance" && storageType != "endurance" {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}

	if cmd.Size == 0 {
		return errors.NewInvalidUsageError(T("-s|--size is required, must be a positive integer.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}
	size := cmd.Size

	if cmd.OsType == "" {
		return errors.NewInvalidUsageError(T("-o|--os-type is required, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}
	osType := cmd.OsType
	if osType != "HYPER_V" && osType != "LINUX" && osType != "VMWARE" && osType != "WINDOWS_2008" && osType != "WINDOWS_GPT" && osType != "WINDOWS" && osType != "XEN" {
		return errors.NewInvalidUsageError(T("-o|--os-type is required, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}
	if cmd.Datacenter == "" {
		return errors.NewInvalidUsageError(T("-d|--datacenter is required.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
	}
	datacenter := cmd.Datacenter
	var orderReceipt datatypes.Container_Product_Order_Receipt
	var err error

	iops := cmd.Iops
	if storageType == "performance" {
		if iops == 0 {
			return errors.NewInvalidUsageError(T("-i|--iops is required with performance volume.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
		}
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
		}
	} else {
		if iops != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops can only be specified with performance volume."))
		}
	}

	tier := cmd.Tier
	if storageType == "endurance" {
		if tier == 0 {
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
		}
		if tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl block volume-options' to check available options.", subs))
		}
	} else {
		if tier != 0 {
			return errors.NewInvalidUsageError(T("-e|--tier can only be specified with endurance volume."))
		}
	}

	outputFormat := cmd.GetOutputFlag()
	if !cmd.Force && outputFormat != "JSON" {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	billingFlag := cmd.Billing

	billing := false
	if billingFlag != "" {
		billingFlag = strings.ToLower(billingFlag)
		if billingFlag != "hourly" && billingFlag != "monthly" {
			return errors.NewInvalidUsageError(T("-b|--billing can only be either hourly or monthly.\nRun '{{.CommandName}} sl file volume-options' to check available options.", subs))
		}
		billing = (billingFlag == "hourly")
	}

	orderReceipt, err = cmd.StorageManager.OrderVolume("block", datacenter, storageType, osType, size, tier, iops, cmd.SnapshotSize, billing)
	if err != nil {
		return errors.NewAPIError(T("Failed to order block volume.Please verify your options and try again.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	for _, item := range orderReceipt.PlacedOrder.Items {
		if item.Description != nil {
			cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
			cmd.UI.Print("")
		}
	}
	cmd.UI.Print(T("You may run '{{.CommandName}} sl block volume-list --order {{.OrderID}}' to find this block volume after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": "ibmcloud"}))

	return nil
}
