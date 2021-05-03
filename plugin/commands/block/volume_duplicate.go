package block

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
			VolumeType: 			"block",
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
	cmd.UI.Print(T("You may run '{{.CommandName}} sl block volume-list --order {{.OrderID}}' to find this block volume after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))
	return nil
}
