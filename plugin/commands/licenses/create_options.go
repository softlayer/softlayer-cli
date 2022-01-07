package licenses

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LicensesOptionsCommand struct {
	UI             terminal.UI
	LicensesManager managers.LicensesManager
}

func NewLicensesOptionsCommand(ui terminal.UI, LicensesManager managers.LicensesManager) (cmd *LicensesOptionsCommand) {
	return &LicensesOptionsCommand{
		UI:             ui,
		LicensesManager: LicensesManager,
	}
}

func (cmd *LicensesOptionsCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{T("Id"), T("Description"), T("KeyName"), T("Capacity"), T("RecurringFee")})
	licenses, err := cmd.LicensesManager.CreateLicensesOptions()
	if err != nil {
		return cli.NewExitError(T("Failed to licenses create options.\n")+err.Error(), 2)
	}

	for _, license := range licenses {
		table.Add(utils.FormatIntPointerName(license.Id),
			utils.FormatStringPointer(license.Description),
			utils.FormatStringPointer(license.KeyName),
			utils.FormatSLFloatPointerToFloat(license.Capacity),
			utils.FormatSLFloatPointerToFloat(license.Prices[0].RecurringFee))

	}
	table.Print()
	return nil
}
