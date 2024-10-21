package block

import (
	"fmt"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeModifyCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	NewSize        int
	NewIops        int
	NewTier        float64
	Force          bool
}

func NewVolumeModifyCommand(sl *metadata.SoftlayerStorageCommand) *VolumeModifyCommand {
	thisCmd := &VolumeModifyCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-modify " + T("IDENTIFIER"),
		Short: T("Modify an existing block storage volume"),
		Long: T(`${COMMAND_NAME} sl {{.storageType}} volume-modify VOLUME_ID [OPTIONS]

   EXAMPLE:
	  ${COMMAND_NAME} sl {{.storageType}} volume-modify 12345678 --new-size 1000 --new-iops 4000 
	  This command modify a volume 12345678 with size is 1000GB, IOPS is 4000.
	  ${COMMAND_NAME} sl {{.storageType}} volume-modify 12345678 --new-size 500 --new-tier 4
	  This command modify a volume 12345678 with size is 500GB, tier level is 4 IOPS per GB.`, sl.StorageI18n),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVarP(&thisCmd.NewSize, "new-size", "c", 0, T("New Size of block volume in GB. ***If no size is given, the original size of volume is used.***\n      Potential Sizes: [20, 40, 80, 100, 250, 500, 1000, 2000, 4000, 8000, 12000]\n      Minimum: [the original size of the volume]"))
	cobraCmd.Flags().IntVarP(&thisCmd.NewIops, "new-iops", "i", 0, T("Performance Storage IOPS, between 100 and 6000 in multiples of 100 [only for performance volumes] ***If no IOPS value is specified, the original IOPS value of the volume will be used.***"))
	cobraCmd.Flags().Float64VarP(&thisCmd.NewTier, "new-tier", "t", 0, T("Endurance Storage Tier (IOPS per GB) [only for endurance volumes] ***If no tier is specified, the original tier of the volume will be used.***"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeModifyCommand) Run(args []string) error {

	volumeID, err := cmd.StorageManager.GetVolumeId(args[0], cmd.StorageType)
	if err != nil {
		return err
	}

	newTier := cmd.NewTier
	size := cmd.NewSize
	iops := cmd.NewIops

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

	modifiedVolume, err := cmd.StorageManager.OrderModifiedVolume("block", volumeID, newTier, size, iops)
	if err != nil {
		return err
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, modifiedVolume)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed successfully!.", map[string]interface{}{"OrderID": *modifiedVolume.OrderId}))
	for _, item := range modifiedVolume.PlacedOrder.Items {
		if item.Description != nil {
			cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
			cmd.UI.Print("")
		}
	}
	cmd.UI.Print(T("You may run '{{.CommandName}} sl block volume-list --order {{.OrderID}}' to find this file volume after it is ready.",
		map[string]interface{}{"OrderID": *modifiedVolume.OrderId, "CommandName": "ibmcloud"}))

	return nil
}
