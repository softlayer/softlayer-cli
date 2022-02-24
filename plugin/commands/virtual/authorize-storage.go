package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AuthorizeStorageCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewAuthorizeStorageCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *AuthorizeStorageCommand) {
	return &AuthorizeStorageCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *AuthorizeStorageCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if c.IsSet("u") {
		storageResult, err := cmd.VirtualServerManager.AuthorizeStorage(vsID, c.String("username-storage"))
		if err != nil {
			return cli.NewExitError(T("Failed to authorize storage to the virtual server instance: {{.Storage}}.\n{{.Error}}",
				map[string]interface{}{"Storage": c.String("username-storage"), "Error": err.Error()}), 2)
		}

		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, storageResult)
		}

		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully authorized storage: {{.Storage}} was added.",
			map[string]interface{}{"Storage": c.String("username-storage")}))
	}

	if c.IsSet("portable-id") {
		portableResult, err := cmd.VirtualServerManager.AttachPortableStorage(vsID, c.Int("portable-id"))
		if err != nil {
			return cli.NewExitError(T("Failed to authorize portable storage to the virtual server instance: {{.PortableID}}.\n{{.Error}}",
				map[string]interface{}{"PortableID": c.Int("portable-id"), "Error": err.Error()}), 2)
		}

		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, portableResult)
		}

		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully authorized storage: {{.PortableID}} was added.",
			map[string]interface{}{"PortableID": c.Int("portable-id")}))
		table := cmd.UI.Table([]string{T("id"), T("CreateDate")})
		table.Add(utils.FormatIntPointer(portableResult.Id), utils.FormatSLTimePointer(portableResult.CreateDate))
		table.Print()
	}

	return nil
}


func VSAuthorizeStorageMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "authorize-storage",
		Description: T("Authorize File, Block and Portable Storage to a Virtual Server"),
		Usage: T(`${COMMAND_NAME} sl vs authorize-storage [OPTIONS] IDENTIFIER

EXAMPLE:
   ${COMMAND_NAME} sl vs authorize-storage --username-storage SL01SL30-37 1234567
   Authorize File, Block and Portable Storage to a Virtual Server.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u, username-storage",
				Usage: T("The storage username to be added to the virtual server."),
			},
			cli.IntFlag{
				Name:  "p, portable-id",
				Usage: T("The portable storage id to be added to the virtual server"),
			},
			metadata.OutputFlag(),
		},
	}
}