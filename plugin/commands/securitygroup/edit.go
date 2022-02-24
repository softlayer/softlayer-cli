package securitygroup

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewEditCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	if c.String("n") == "" && c.String("d") == "" {
		return errors.NewInvalidUsageError(T("Either -n, --name or -d, --description is required to edit security group."))
	}
	err = cmd.NetworkManager.EditSecurityGroup(groupID, c.String("n"), c.String("d"))
	if err != nil {
		return cli.NewExitError(T("Failed to edit security group {{.ID}}.\n", map[string]interface{}{"ID": groupID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Security group {{.ID}} is updated.", map[string]interface{}{"ID": groupID}))
	return nil
}

func SecurityGroupEditMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "edit",
		Description: T("Edit details of a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup edit SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the security group"),
			},
			cli.StringFlag{
				Name:  "d,description",
				Usage: T("The description of the security group"),
			},
		},
	}
}
