package order

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type PresetListCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
}

func NewPresetListCommand(ui terminal.UI, orderManager managers.OrderManager) (cmd *PresetListCommand) {
	return &PresetListCommand{
		UI:           ui,
		OrderManager: orderManager,
	}
}

func (cmd *PresetListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	packageKeyname := c.Args()[0]

	keyword := c.String("keyword")

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	presets, err := cmd.OrderManager.ListPreset(packageKeyname, keyword)
	if err != nil {
		return cli.NewExitError(T(fmt.Sprintf("Failed to list presets: %s\n", err.Error())), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, presets)
	}

	cmd.Print(presets)
	return nil
}

func (cmd *PresetListCommand) Print(presets []datatypes.Product_Package_Preset) {
	table := cmd.UI.Table([]string{T("category"), T("Key Name"), T("Description")})

	for _, preset := range presets {
		table.Add(utils.FormatStringPointer(preset.Name),
			utils.FormatStringPointer(preset.KeyName),
			utils.FormatStringPointer(preset.Description))
	}
	table.Print()
}
