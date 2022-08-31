package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialsCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCredentialsCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CredentialsCommand) {
	return &CredentialsCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CredentialsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError("This command requires one argument.")
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, "operatingSystem[passwords[username,password]]")
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		if virtualGuest.OperatingSystem != nil {
			return utils.PrintPrettyJSON(cmd.UI, virtualGuest.OperatingSystem.Passwords)
		}
		return utils.PrintPrettyJSON(cmd.UI, virtualGuest.OperatingSystem)
	}

	table := cmd.UI.Table([]string{T("username"), T("password")})
	if virtualGuest.OperatingSystem != nil && len(virtualGuest.OperatingSystem.Passwords) > 0 {
		for _, pwd := range virtualGuest.OperatingSystem.Passwords {
			table.Add(utils.FormatStringPointer(pwd.Username), utils.FormatStringPointer(pwd.Password))
		}
	}
	table.Print()
	return nil
}

func VSCredentialsMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "credentials",
		Description: T("List virtual server instance credentials"),
		Usage: T(`${COMMAND_NAME} sl vs credentials IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs credentials 12345678
   This command lists all username and password pairs of virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
