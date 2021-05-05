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
