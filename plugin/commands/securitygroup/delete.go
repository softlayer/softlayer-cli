package securitygroup

import (
	"strconv"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewDeleteCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *DeleteCommand) {
	return &DeleteCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *DeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will delete security group {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": groupID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.NetworkManager.DeleteSecurityGroup(groupID)
	if err != nil {
		return cli.NewExitError(T("Failed to delete security group {{.ID}}.\n", map[string]interface{}{"ID": groupID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Security group {{.ID}} is deleted.", map[string]interface{}{"ID": groupID}))
	return nil
}

func SecurityGroupDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "delete",
		Description: T("Delete the given security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup delete SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
