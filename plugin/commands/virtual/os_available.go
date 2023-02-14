package virtual

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OsAvailableCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewOsAvailableCommand(sl *metadata.SoftlayerCommand) (cmd *OsAvailableCommand) {
	thisCmd := &OsAvailableCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "os-available",
		Short: T("Get all available Operating Systems."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd

	return thisCmd
}

func (cmd *OsAvailableCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	availables, err := cmd.OrderManager.ListItems("PUBLIC_CLOUD_SERVER", "", "os")

	if err != nil {
		return errors.NewAPIError(T("Failed to list available OS's."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("KeyName"), T("Description"), T("Hourly"), T("Monthly"), T("Setup")})
	for _, availableOs := range availables {
		hourly := "-"
		monthly := "-"
		setup := "-"
		if availableOs.Prices != nil && len(availableOs.Prices) > 0 {
			if availableOs.Prices[0].HourlyRecurringFee != nil {
				hourly = fmt.Sprintf("%.2f", *availableOs.Prices[0].HourlyRecurringFee)
			}
			if availableOs.Prices[0].LaborFee != nil {
				monthly = fmt.Sprintf("%.2f", *availableOs.Prices[0].LaborFee)
			}
			if availableOs.Prices[0].SetupFee != nil {
				setup = fmt.Sprintf("%.2f", *availableOs.Prices[0].SetupFee)
			}
		}
		table.Add(
			utils.FormatIntPointer(availableOs.Id),
			utils.FormatStringPointer(availableOs.KeyName),
			utils.FormatStringPointer(availableOs.Description),
			hourly,
			monthly,
			setup,
		)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}
