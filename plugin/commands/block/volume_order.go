package block

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
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

func (cmd *VolumeOrderCommand) Run(c *cli.Context) error {
	if !c.IsSet("t") {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	storageType := c.String("t")
	if storageType != "performance" && storageType != "endurance" {
		return errors.NewInvalidUsageError(T("-t|--storage-type is required, must be either performance or endurance.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}

	if !c.IsSet("s") {
		return errors.NewInvalidUsageError(T("-s|--size is required, must be a positive integer.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	size := c.Int("s")

	if !c.IsSet("o") {
		return errors.NewInvalidUsageError(T("-o|--os-type is required, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	osType := c.String("o")
	if osType != "HYPER_V" && osType != "LINUX" && osType != "VMWARE" && osType != "WINDOWS_2008" && osType != "WINDOWS_GPT" && osType != "WINDOWS" && osType != "XEN" {
		return errors.NewInvalidUsageError(T("-o|--os-type is required, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	if !c.IsSet("d") {
		return errors.NewInvalidUsageError(T("-d|--datacenter is required.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}
	datacenter := c.String("d")
	var orderReceipt datatypes.Container_Product_Order_Receipt
	var err error

	iops := c.Int("i")
	if storageType == "performance" {
		if iops == 0 {
			return errors.NewInvalidUsageError(T("-i|--iops is required with performance volume.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
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
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("-e|--tier is required with endurance volume in IOPS/GB, options are: 0.25, 2, 4, 10.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
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

	orderReceipt, err = cmd.StorageManager.OrderVolume("block", datacenter, storageType, osType, size, tier, iops, c.Int("n"), billing)
	if err != nil {
		return cli.NewExitError(T("Failed to order block volume.Please verify your options and try again.\n")+err.Error(), 2)
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
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))

	return nil
}
