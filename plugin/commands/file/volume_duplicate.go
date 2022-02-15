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

type VolumeDuplicateCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	Context        plugin.PluginContext
}

func NewVolumeDuplicateCommand(ui terminal.UI, storageManager managers.StorageManager, context plugin.PluginContext) (cmd *VolumeDuplicateCommand) {
	return &VolumeDuplicateCommand{
		UI:             ui,
		StorageManager: storageManager,
		Context:        context,
	}
}

func FileVolumeDuplicateMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-duplicate",
		Description: T("Order a file volume by duplicating an existing volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-duplicate VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-duplicate 12345678 
   This command shows order a new volume by duplicating the volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "o,origin-snapshot-id",
				Usage: T("ID of an original volume snapshot to use for duplication"),
			},
			cli.IntFlag{
				Name:  "s,duplicate-size",
				Usage: T("Size of duplicate file volume in GB, if no size is specified, the size of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "i,duplicate-iops",
				Usage: T("Performance Storage IOPS, between 100 and 6000 in multiples of 100, if no IOPS value is specified, the IOPS value of the original volume will be used"),
			},
			cli.Float64Flag{
				Name:  "t,duplicate-tier",
				Usage: T("Endurance Storage Tier, if no tier is specified, the tier of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "n,duplicate-snapshot-size",
				Usage: T("The size of snapshot space to order for the duplicate, if no snapshot space size is specified, the snapshot space size of the original volume will be used"),
				Value: -1,
			},
			cli.BoolFlag{
				Name:  "d,dependent-duplicate",
				Usage: T("Whether or not this duplicate will be a dependent duplicate of the origin volume."),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *VolumeDuplicateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

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
	config := managers.DuplicateOrderConfig{
			VolumeType: 			"file",
			OriginalVolumeId: 		volumeID,
			OriginalSnapshotId: 	c.Int("o"),
			DuplicateSize: 			c.Int("s"),
			DuplicateIops: 			c.Int("i"),
			DuplicateTier: 			c.Float64("t"),
			DuplicateSnapshotSize: 	c.Int("n"),
			DependentDuplicate:		c.Bool("d"),
	}
	orderReceipt, err := cmd.StorageManager.OrderDuplicateVolume(config)
	if err != nil {
		return cli.NewExitError(T("Failed to order duplicate volume from {{.VolumeID}}.Please verify your options and try again.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
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
	cmd.UI.Print(T("You may run '{{.CommandName}} sl file volume-list --order {{.OrderID}}' to find this file volume after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))
	return nil
}
