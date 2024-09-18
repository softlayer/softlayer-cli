package block

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotOrderCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Size           int
	Tier           float64
	Iops           int
	Upgrade        bool
	Force          bool
}

var TIERS = []float64{0.25, 2, 4, 10}

func NewSnapshotOrderCommand(sl *metadata.SoftlayerStorageCommand) *SnapshotOrderCommand {
	thisCmd := &SnapshotOrderCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "snapshot-order " + T("IDENTIFIER"),
		Short: T("Order snapshot space for a block storage volume"),
		Long: T(`See https://cloud.ibm.com/docs/BlockStorage?topic=BlockStorage-getting-started for sizing options.
${COMMAND_NAME} sl block volume-options' to get available options.

EXAMPLE:
   ${COMMAND_NAME} sl {{.storageType}} snapshot-order 12345678 -s 1000 -t 4 
   This command orders snapshot space for volume with ID 12345678, the size is 1000GB, the tier level is 4 IOPS per GB.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVarP(&thisCmd.Size, "size", "s", 0, T("Size of snapshot space to create in GB"))
	cobraCmd.Flags().Float64VarP(&thisCmd.Tier, "tier", "t", 0,
		T("Endurance Storage Tier (IOPS per GB) of the block volume for which space is ordered [optional], options are: 0.25,2,4,10"))
	cobraCmd.Flags().IntVarP(&thisCmd.Iops, "iops", "i", 0,
		T("Performance Storage IOPs, between 100 and 6000 in multiples of 100"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Upgrade, "upgrade", "u", false,
		T("Flag to indicate that the order is an upgrade"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	cobraCmd.MarkFlagRequired("size") // #nosec G104 -- Doesn't matter if this errors
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SnapshotOrderCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	if cmd.Tier > 0 {
		if !slices.Contains(TIERS, cmd.Tier) {
			return slErr.NewInvalidUsageError(T("[-t|--tier] is optional, options are: 0.25,2,4,10."))
		}
	}

	if cmd.Iops > 0 {
		if cmd.Iops%100 != 0 {
			return slErr.NewInvalidUsageError(T("-i|--iops must be a multiple of 100."))
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
	orderReceipt, err := cmd.StorageManager.OrderSnapshotSpace(
		cmd.GetStorageType(), volumeID, cmd.Size, cmd.Tier, cmd.Iops, cmd.Upgrade)
	if err != nil {
		return slErr.NewAPIError(
			T("Failed to order snapshot space for volume {{.VolumeID}}.Please verify your options and try again.\n",
				map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
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
