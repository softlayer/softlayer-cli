package block

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeDuplicateCommand struct {
	*metadata.SoftlayerCommand
	Command               *cobra.Command
	StorageManager        managers.StorageManager
	OriginSnapshotId      int
	DuplicateSize         int
	DuplicateIops         int
	DuplicateTier         float64
	DuplicateSnapshotSize int
	DependentDuplicate    bool
	Force                 bool
}

func NewVolumeDuplicateCommand(sl *metadata.SoftlayerCommand) *VolumeDuplicateCommand {
	thisCmd := &VolumeDuplicateCommand{
		SoftlayerCommand: sl,
		StorageManager:   managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-duplicate " + T("IDENTIFIER"),
		Short: T("Order a block volume by duplicating an existing volume"),
		Long: T(`${COMMAND_NAME} sl block volume-duplicate VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block volume-duplicate 12345678 
   This command shows order a new volume by duplicating the volume with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVarP(&thisCmd.OriginSnapshotId, "origin-snapshot-id", "o", 0, T("ID of an origin volume snapshot to use for duplication"))
	cobraCmd.Flags().IntVarP(&thisCmd.DuplicateSize, "duplicate-size", "s", 0, T("Size of duplicate block volume in GB, if no size is specified, the size of the original volume will be used"))
	cobraCmd.Flags().IntVarP(&thisCmd.DuplicateIops, "duplicate-iops", "i", 0, T("Performance Storage IOPS, between 100 and 6000 in multiples of 100, if no IOPS value is specified, the IOPS value of the original volume will be used"))
	cobraCmd.Flags().Float64VarP(&thisCmd.DuplicateTier, "duplicate-tier", "t", 0, T("Endurance Storage Tier, if no tier is specified, the tier of the original volume will be used"))
	cobraCmd.Flags().IntVarP(&thisCmd.DuplicateSnapshotSize, "duplicate-snapshot-size", "n", -1, T("The size of snapshot space to order for the duplicate, if no snapshot space size is specified, the snapshot space size of the origin volume will be used"))
	cobraCmd.Flags().BoolVarP(&thisCmd.DependentDuplicate, "dependent-duplicate", "d", false, T("Whether or not this duplicate will be a dependent duplicate of the origin volume."))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeDuplicateCommand) Run(args []string) error {

	volumeID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
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
	config := managers.DuplicateOrderConfig{
		VolumeType:            "block",
		OriginalVolumeId:      volumeID,
		OriginalSnapshotId:    cmd.OriginSnapshotId,
		DuplicateSize:         cmd.DuplicateSize,
		DuplicateIops:         cmd.DuplicateIops,
		DuplicateTier:         cmd.DuplicateTier,
		DuplicateSnapshotSize: cmd.DuplicateSnapshotSize,
		DependentDuplicate:    cmd.DependentDuplicate,
	}
	orderReceipt, err := cmd.StorageManager.OrderDuplicateVolume(config)
	if err != nil {
		return slErr.NewAPIError(T("Failed to order duplicate volume from {{.VolumeID}}.Please verify your options and try again.\n", map[string]interface{}{"VolumeID": volumeID}), err.Error(), 2)
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
