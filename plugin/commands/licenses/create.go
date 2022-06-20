package licenses

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI              terminal.UI
	LicensesManager managers.LicensesManager
}

func NewCreateCommand(ui terminal.UI, licensesManager managers.LicensesManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:              ui,
		LicensesManager: licensesManager,
	}
}

func CreateMetaData() cli.Command {
	return cli.Command{
		Category:    "licenses",
		Name:        "create",
		Description: T("Order/Create License."),
		Usage:       "${COMMAND_NAME} sl licenses create",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "key",
				Usage: T("The VMware License Key. To get the required package you can use the command sl licenses create-options Package. E.g VMWARE_VSAN_ENTERPRISE_TIER_III_65_124_TB_6_X_2  [required]"),
			},
			cli.StringFlag{
				Name:  "datacenter",
				Usage: T("Datacenter shortname  [required]"),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("datacenter") || !c.IsSet("key") {
		return slErr.NewInvalidUsageError(T("This command requires two arguments."))
	}
	
	table := cmd.UI.Table([]string{T("Name"), T("Value")})

	orderLicense, err := cmd.LicensesManager.CreateLicense(c.String("datacenter"), c.String("key"))
	if err != nil {
		return cli.NewExitError(T("Failed to create the license.\n")+err.Error(), 2)
	}
	table.Add("Id", utils.FormatIntPointer(orderLicense.OrderId))
	table.Add("Created", utils.FormatSLTimePointer(orderLicense.OrderDate))

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
