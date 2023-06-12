package order

import (
	"fmt"

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
	Prices       bool
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
	cobraCmd.Flags().BoolVar(&thisCmd.Prices, "prices", false, T("Lists the prices for each item in this preset"))

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

	if cmd.Prices {
		cmd.PrintPresetPrices(presets)
	} else {
		cmd.Print(presets)
	}

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

func (cmd *PresetListCommand) PrintPresetPrices(presets []datatypes.Product_Package_Preset) {
	table := cmd.UI.Table([]string{T("Key Name"), T("Price Id"), T("Hourly"), T("Monthly"), T("Restriction"), T("Location")})

	for _, preset := range presets {
		locations := "[]"
		if preset.Locations != nil && len(preset.Locations) > 0 {
			locations = ""
			for _, location := range preset.Locations {
				locations = locations + *location.Name + ", "
			}
			locations = "[" + locations[0:len(locations)-2] + "]"
		}

		crMax, crMin, crType, hourly, monthly := "-", "-", "-", "-", "-"
		if len(preset.Prices) > 0 {
			if preset.Prices[0].CapacityRestrictionMaximum != nil {
				crMax = *preset.Prices[0].CapacityRestrictionMaximum
			}

			if preset.Prices[0].CapacityRestrictionMaximum != nil {
				crMin = *preset.Prices[0].CapacityRestrictionMinimum
			}

			if preset.Prices[0].CapacityRestrictionMaximum != nil {
				crType = *preset.Prices[0].CapacityRestrictionType
			}

			if preset.Prices[0].HourlyRecurringFee != nil {
				hourly = fmt.Sprintf("%.2f", *preset.Prices[0].HourlyRecurringFee)
			}

			if preset.Prices[0].HourlyRecurringFee != nil {
				monthly = fmt.Sprintf("%.2f", *preset.Prices[0].RecurringFee)
			}
		}

		table.Add(
			utils.FormatStringPointer(preset.KeyName),
			utils.FormatIntPointer(preset.Id),
			hourly,
			monthly,
			fmt.Sprintf("%s - %s %s", crMin, crMax, crType),
			locations,
		)
	}
	table.Print()
}
