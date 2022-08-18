package licenses

import (
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LicensesOptionsCommand struct {
	*metadata.SoftlayerCommand
	Command         *cobra.Command
	LicensesManager managers.LicensesManager
}

func NewLicensesOptionsCommand(sl *metadata.SoftlayerCommand) *LicensesOptionsCommand {
	thisCmd := &LicensesOptionsCommand{
		SoftlayerCommand: sl,
		LicensesManager:  managers.NewLicensesManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create-options",
		Short: T("Server order options for a given chassis."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *LicensesOptionsCommand) Run(args []string) error {
	table := cmd.UI.Table([]string{T("Id"), T("Description"), T("KeyName"), T("Capacity"), T("RecurringFee")})
	licenses, err := cmd.LicensesManager.CreateLicensesOptions()
	if err != nil {
		return slErr.NewAPIError(T("Failed to licenses create options."), err.Error(), 2)
	}

	for _, license := range licenses {
		for _, item := range license.Items {
			table.Add(utils.FormatIntPointerName(item.Id),
				utils.FormatStringPointer(item.Description),
				utils.FormatStringPointer(item.KeyName),
				utils.FormatSLFloatPointerToFloat(item.Capacity),
				utils.FormatSLFloatPointerToFloat(item.Prices[0].RecurringFee))
		}
	}
	table.Print()
	return nil
}
