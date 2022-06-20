package licenses

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CancelItemCommand struct {
	UI              terminal.UI
	LicensesManager managers.LicensesManager
}

func NewCancelItemCommand(ui terminal.UI, licensesManager managers.LicensesManager) (cmd *CancelItemCommand) {
	return &CancelItemCommand{
		UI:              ui,
		LicensesManager: licensesManager,
	}
}

func CancelItemMetaData() cli.Command {
	return cli.Command{
		Category:    "licenses",
		Name:        "cancel",
		Description: T("Cancel a license."),
		Usage:       T(`${COMMAND_NAME} sl licenses cancel KEY`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Immediate cancellation."),
			},
		},
	}
}

func (cmd *CancelItemCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	key := c.Args()[0]

	err := cmd.LicensesManager.CancelItem(key, c.Bool("immediate"))
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find license with key: {{.key}}.\n", map[string]interface{}{"key": key})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to cancel license: {{.key}}.\n", map[string]interface{}{"key": key})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("License: {{.key}} was cancelled.", map[string]interface{}{"key": key}))
	return nil
}
