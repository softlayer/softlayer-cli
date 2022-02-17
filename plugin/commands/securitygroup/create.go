package securitygroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewCreateCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	group, err := cmd.NetworkManager.CreateSecurityGroup(c.String("n"), c.String("d"))
	if err != nil {
		return cli.NewExitError(T("Failed to create security group with name {{.Name}}.\n",
			map[string]interface{}{"Name": c.String("n")})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, group)
	}
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(group.Id))
	table.Add(T("Name"), utils.FormatStringPointer(group.Name))
	table.Add(T("Description"), utils.FormatStringPointer(group.Description))
	table.Add(T("Created"), utils.FormatSLTimePointer(group.CreateDate))
	table.Print()
	return nil
}

func SecurityGroupCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "create",
		Description: T("Create a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup create [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the security group"),
			},
			cli.StringFlag{
				Name:  "d,description",
				Usage: T("The description of the security group"),
			},
			metadata.OutputFlag(),
		},
	}
}
