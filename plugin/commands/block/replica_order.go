package block

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReplicaOrderCommand struct {
	*metadata.SoftlayerCommand
	Command          *cobra.Command
	StorageManager   managers.StorageManager
	SnapshotSchedule string
	Datacenter       string
	Tier             float64
	Iops             int
	OsType           string
	Force            bool
}

func NewReplicaOrderCommand(sl *metadata.SoftlayerCommand) *ReplicaOrderCommand {
	thisCmd := &ReplicaOrderCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "replica-order " + T("IDENTIFIER"),
		Short: T("Order a block storage replica volume"),
		Long: T(`${COMMAND_NAME} sl block replica-order VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-order 12345678 -s DAILY -d dal09 --tier 4 --os-type LINUX
   This command orders a replica for volume with ID 12345678, which performs DAILY replication, is located at dal09, tier level is 4, OS type is Linux.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.SnapshotSchedule, "snapshot-schedule", "s", "", T("Snapshot schedule to use for replication. Options are: HOURLY,DAILY,WEEKLY [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Short name of the datacenter for the replica. For example, dal09 [required]"))
	cobraCmd.Flags().Float64VarP(&thisCmd.Tier, "tier", "t", 0, T("Endurance Storage Tier (IOPS per GB) of the primary volume for which a replica is ordered [optional], options are: 0.25,2,4,10,if no tier is specified, the tier of the original volume will be used"))
	cobraCmd.Flags().IntVarP(&thisCmd.Iops, "iops", "i", 0, ("Performance Storage IOPs, between 100 and 6000 in multiples of 100,if no IOPS value is specified, the IOPS value of the original volume will be used"))
	cobraCmd.Flags().StringVarP(&thisCmd.OsType, "os-type", "o", "", T("Operating System Type (eg. LINUX) of the primary volume for which a replica is ordered [optional], options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ReplicaOrderCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	snapshotSchedule := cmd.SnapshotSchedule
	if snapshotSchedule == "" || (snapshotSchedule != "HOURLY" && snapshotSchedule != "DAILY" && snapshotSchedule != "WEEKLY") {
		return errors.NewInvalidUsageError(T("[-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY."))
	}

	datacenter := cmd.Datacenter
	if datacenter == "" {
		// Need a better way to get command_name
		return errors.NewInvalidUsageError(T("[-d|--datacenter] is required.\n Run '{{.CommandName}} sl block volume-options' to get available options.",
			map[string]interface{}{"CommandName": "${COMMAND_NAME}"}))
	}

	tier := cmd.Tier
	if tier > 0 {
		if tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("[-t|--tier] is optional, options are: 0.25,2,4,10."))
		}
	}
	iops := cmd.Iops
	if iops != 0 {
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": "${COMMAND_NAME}"}))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": "${COMMAND_NAME}"}))
		}
	}

	osType := cmd.OsType
	if osType != "" {
		if osType != "HYPER_V" && osType != "LINUX" && osType != "VMWARE" && osType != "WINDOWS_2008" && osType != "WINDOWS_GPT" && osType != "WINDOWS" && osType != "XEN" {
			return errors.NewInvalidUsageError(T("-o|--os-type is optional, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN."))
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
	orderReceipt, err := cmd.StorageManager.OrderReplicantVolume("block", volumeID, snapshotSchedule, datacenter, tier, iops, osType)
	if err != nil {
		return errors.NewAPIError(T("Failed to order replicant for volume {{.VolumeID}}.Please verify your options and try again.\n", map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	for _, item := range orderReceipt.PlacedOrder.Items {
		cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
		cmd.UI.Print("")
	}
	return nil
}
