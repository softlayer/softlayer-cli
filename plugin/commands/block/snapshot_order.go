package block

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SnapshotOrderCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	Context        plugin.PluginContext
}

func NewSnapshotOrderCommand(ui terminal.UI, storageManager managers.StorageManager, context plugin.PluginContext) (cmd *SnapshotOrderCommand) {
	return &SnapshotOrderCommand{
		UI:             ui,
		StorageManager: storageManager,
		Context:        context,
	}
}

func BlockSnapshotOrderMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "snapshot-order",
		Description: T("Order snapshot space for a block storage volume"),
		Usage: T(`${COMMAND_NAME} sl block snapshot-order VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block snapshot-order 12345678 -s 1000 -t 4 
   This command orders snapshot space for volume with ID 12345678, the size is 1000GB, the tier level is 4 IOPS per GB.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "s,size",
				Usage: T("Size of snapshot space to create in GB  [required]"),
			},
			cli.Float64Flag{
				Name:  "t,tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) of the block volume for which space is ordered [optional], options are: 0.25,2,4,10"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100"),
			},
			cli.BoolFlag{
				Name:  "u,upgrade",
				Usage: T("Flag to indicate that the order is an upgrade"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *SnapshotOrderCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("s") {
		return errors.NewInvalidUsageError(T("[-s|--size] is required.\nRun '{{.CommandName}} sl block volume-options' to get available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	size := c.Int("s")

	tier := float64(0.0)
	if c.IsSet("t") {
		tier := c.Float64("t")
		if tier != 0 && tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("[-t|--tier] is optional, options are: 0.25,2,4,10."))
		}
	}
	iops := 0
	if c.IsSet("i") {
		iops = c.Int("i")
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))

		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))

		}
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
	orderReceipt, err := cmd.StorageManager.OrderSnapshotSpace("block", volumeID, size, tier, iops, c.Bool("u"))
	if err != nil {
		return cli.NewExitError(T("Failed to order snapshot space for volume {{.VolumeID}}.Please verify your options and try again.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
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
