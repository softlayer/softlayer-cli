package order

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PresetListCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Keyword      string
}

func NewPresetListCommand(sl *metadata.SoftlayerCommand) (cmd *PresetListCommand) {
	thisCmd := &PresetListCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "preset-list",
		Short: T("List package presets"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Keyword, "keyword", "", T("A word (or string) used to filter presets"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PresetListCommand) Run(args []string) error {
	packageKeyname := args[0]

	keyword := cmd.Keyword

	outputFormat := cmd.GetOutputFlag()

	presets, err := cmd.OrderManager.ListPreset(packageKeyname, keyword)
	if err != nil {
		return errors.NewAPIError(T("Failed to list presets"), err.Error(), 2)
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
