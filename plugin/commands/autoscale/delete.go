package autoscale

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type DeleteCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
}

func NewDeleteCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager) (cmd *DeleteCommand) {
	return &DeleteCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
	}
}

func (cmd *DeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the AutoScale Group instance: {{.autoScaleGroupId}} and all its members, this action cannot be undone. Continue?", map[string]interface{}{"autoScaleGroupId": autoScaleGroupId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	response, err := cmd.AutoScaleManager.Delete(autoScaleGroupId)
	if err != nil {
		return cli.NewExitError(T("Failed to delete AutoScale Group.\n")+err.Error(), 2)
	}

	if response {
		cmd.UI.Ok()
		cmd.UI.Print(T("Auto Scale Group was deleted successfully"))
	}
	return nil
}

func AutoScaleDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "delete",
		Description: T("Delete this group and destroy all members of it"),
		Usage: T(`${COMMAND_NAME} sl autoscale delete IDENTIFIER

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale delete 123456`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
