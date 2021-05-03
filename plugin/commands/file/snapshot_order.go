package file

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
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

func (cmd *SnapshotOrderCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}
	if !c.IsSet("s") {
		return errors.NewInvalidUsageError(T("[-s|--size] is required.\nRun '{{.CommandName}} sl file volume-options' to get available options.",
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
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl file volume-options' to check available options.",
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
	orderReceipt, err := cmd.StorageManager.OrderSnapshotSpace("file", volumeID, size, tier, iops, c.Bool("u"))
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
		if item.Description != nil {
			cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
			cmd.UI.Print("")
		}
	}
	return nil
}
