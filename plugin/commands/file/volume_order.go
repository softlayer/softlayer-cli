package file

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type VolumeOrderCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	Context        plugin.PluginContext
}

func NewVolumeOrderCommand(ui terminal.UI, storageManager managers.StorageManager, context plugin.PluginContext) (cmd *VolumeOrderCommand) {
	return &VolumeOrderCommand{
		UI:             ui,
		StorageManager: storageManager,
		Context:        context,
	}
}

func FileVolumeOrderMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-order",
		Description: T("Order a file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-order [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-order --storage-type performance --size 1000 --iops 4000  -d dal09
   This command orders a performance volume with size is 1000GB, IOPS is 4000, located at dal09.
   ${COMMAND_NAME} sl file volume-order --storage-type endurance --size 500 --tier 4 -d dal09 --snapshot-size 500
   This command orders a endurance volume with size is 500GB, tier level is 4 IOPS per GB,located at dal09, and additional snapshot space size is 500GB.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "t,storage-type",
				Usage: T("Type of storage volume [required], options are: performance,endurance"),
			},
			cli.IntFlag{
				Name:  "s,size",
				Usage: T("Size of storage volume in GB [required]"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100 [required for storage-type performance]"),
			},
			cli.Float64Flag{
				Name:  "e,tier",
				Usage: T("Endurance Storage Tier (IOP per GB) [required for storage-type endurance], options are: 0.25,2,4,10"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter short name [required]"),
			},
			cli.IntFlag{
				Name:  "n,snapshot-size",
				Usage: T("Optional parameter for ordering snapshot space along with the volume"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Optional parameter for Billing rate (default to monthly), options are: hourly, monthly"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *VolumeOrderCommand) Run(c *cli.Context) error {
	if !c.IsSet("t") {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	storageType := c.String("t")
	if storageType != "performance" && storageType != "endurance" {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}

	if !c.IsSet("s") {
		return errors.NewInvalidUsageError(T("-s|--size is required, must be a positive integer.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	size := c.Int("s")

	datacenter := c.String("d")
	if !c.IsSet("d") || datacenter == "" {
		return errors.NewInvalidUsageError(T("-d|--datacenter is required.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	var orderReceipt datatypes.Container_Product_Order_Receipt
	var err error

	iops := c.Int("i")
	if storageType == "performance" {
		if iops == 0 {
			return errors.NewInvalidUsageError(T("-i|--iops is required with performance volume.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
	} else {
		if iops != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops can only be specified with performance volume."))
		}
	}

	tier := c.Float64("e")
	if storageType == "endurance" {
		if tier == 0 {
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
	} else {
		if tier != 0 {
			return errors.NewInvalidUsageError(T("-e|--tier can only be specified with endurance volume."))
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

	billingFlag := c.String("b")

	billing := false
	if billingFlag != "" {
		billingFlag = strings.ToLower(billingFlag)
		if billingFlag != "hourly" && billingFlag != "monthly" {
			return errors.NewInvalidUsageError(T("-b|--billing can only be either hourly or monthly.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		billing = (billingFlag == "hourly")
	}

	orderReceipt, err = cmd.StorageManager.OrderVolume("file", datacenter, storageType, "", size, tier, iops, c.Int("n"), billing)
	if err != nil {
		return cli.NewExitError(T("Failed to order file volume.Please verify your options and try again.\n")+err.Error(), 2)
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
