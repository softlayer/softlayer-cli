package subnet

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

	subnetID, err := utils.ResolveSubnetId(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	if !c.IsSet("tags") && !c.IsSet("note") {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if c.IsSet("tags") {
		tags := c.String("tags")
		response, err := cmd.NetworkManager.SetSubnetTags(subnetID, tags)
		if err != nil {
			return cli.NewExitError(T("Failed to set tags: {{.tags}}.\n", map[string]interface{}{"tags": tags})+err.Error(), 2)
		}
		if response {
			cmd.UI.Ok()
			cmd.UI.Print(T("Set tags successfully"))
		}
	}

	if c.IsSet("note") {
		note := c.String("note")
		response, err := cmd.NetworkManager.SetSubnetNote(subnetID, note)
		if err != nil {
			return cli.NewExitError(T("Failed to set note: {{.note}}.\n", map[string]interface{}{"note": note})+err.Error(), 2)
		}
		if response {
			cmd.UI.Ok()
			cmd.UI.Print(T("Set note successfully"))
		}
	}

	return nil
}

func SubnetEditMetaData() cli.Command {
	return cli.Command{
		Category:    "subnet",
		Name:        "edit",
		Description: T("Edit note and tags of a subnet."),
		Usage: T(`${COMMAND_NAME} sl subnet edit IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet edit 12345678 --note myNote
   ${COMMAND_NAME} sl subnet edit 12345678 --tags tag1`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "note",
				Usage: T("The note "),
			},
			cli.StringFlag{
				Name:  "tags",
				Usage: T("Comma separated list of tags, enclosed in quotes. 'tag1, tag2'"),
			},
		},
	}
}
