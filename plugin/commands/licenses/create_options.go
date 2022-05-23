package licenses

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type LicensesOptionsCommand struct {
	UI              terminal.UI
	LicensesManager managers.LicensesManager
}

func NewLicensesOptionsCommand(ui terminal.UI, licensesManager managers.LicensesManager) (cmd *LicensesOptionsCommand) {
	return &LicensesOptionsCommand{
		UI:              ui,
		LicensesManager: licensesManager,
	}
}

func LicensesCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "licenses",
		Name:        "create-options",
		Description: T("Server order options for a given chassis"),
		Usage:       "${COMMAND_NAME} sl licenses create-options",
		Flags:       []cli.Flag{},
	}
}

func (cmd *LicensesOptionsCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{T("Id"), T("Description"), T("KeyName"), T("Capacity"), T("RecurringFee")})
	licenses, err := cmd.LicensesManager.CreateLicensesOptions()
	if err != nil {
		return cli.NewExitError(T("Failed to licenses create options.\n")+err.Error(), 2)
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
