package file

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeModifyCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	Context        plugin.PluginContext
}

func NewVolumeModifyCommand(ui terminal.UI, storageManager managers.StorageManager, context plugin.PluginContext) (cmd *VolumeModifyCommand) {
	return &VolumeModifyCommand{
		UI:             ui,
		StorageManager: storageManager,
		Context:        context,
	}
}

func FileVolumeModifyMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-modify",
		Description: T("Modify an existing file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-modify VOLUME_ID [OPTIONS]

   EXAMPLE:
	  ${COMMAND_NAME} sl file volume-modify 12345678 --new-size 1000 --new-iops 4000 
	  This command modify a volume 12345678 with size is 1000GB, IOPS is 4000.
	  ${COMMAND_NAME} sl file volume-modify 12345678 --new-size 500 --new-tier 4
	  This command modify a volume 12345678 with size is 500GB, tier level is 4 IOPS per GB.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,new-size",
				Usage: T("New Size of file volume in GB. ***If no size is given, the original size of volume is used.***\n      Potential Sizes: [20, 40, 80, 100, 250, 500, 1000, 2000, 4000, 8000, 12000]\n      Minimum: [the original size of the volume]"),
			},
			cli.IntFlag{
				Name:  "i,new-iops",
				Usage: T("Performance Storage IOPS, between 100 and 6000 in multiples of 100 [only for performance volumes] ***If no IOPS value is specified, the original IOPS value of the volume will be used.***\n      Requirements: [If original IOPS/GB for the volume is less than 0.3, new IOPS/GB must also be less than 0.3. If original IOPS/GB for the volume is greater than or equal to 0.3, new IOPS/GB for the volume must also be greater than or equal to 0.3.]"),
			},
			cli.Float64Flag{
				Name:  "t, new-tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) [only for endurance volumes] ***If no tier is specified, the original tier of the volume will be used.***\n      Requirements: [If original IOPS/GB for the volume is 0.25, new IOPS/GB for the volume must also be 0.25. If original IOPS/GB for the volume is greater than 0.25, new IOPS/GB for the volume must also be greater than 0.25.]"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *VolumeModifyCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	newTier := c.Float64("t")
	size := c.Int("c")
	iops := c.Int("i")

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("f") && outputFormat != "JSON" {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	modifiedVolume, err := cmd.StorageManager.OrderModifiedVolume("file", volumeID, newTier, size, iops)
	if err != nil {
		return err
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, modifiedVolume)
	}

	cmd.UI.Print(T("Order {{.OrderID}} was placed successfully!.", map[string]interface{}{"OrderID": *modifiedVolume.OrderId}))
	for _, item := range modifiedVolume.PlacedOrder.Items {
		if item.Description != nil {
			cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
			cmd.UI.Print("")
		}
	}
	cmd.UI.Print(T("You may run '{{.CommandName}} sl file volume-list --order {{.OrderID}}' to find this file volume after it is ready.",
		map[string]interface{}{"OrderID": *modifiedVolume.OrderId, "CommandName": cmd.Context.CLIName()}))

	return nil
}
