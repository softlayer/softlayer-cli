package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AuthorizeStorageCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewAuthorizeStorageCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *AuthorizeStorageCommand) {
	return &AuthorizeStorageCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *AuthorizeStorageCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	storageResult, err := cmd.HardwareManager.AuthorizeStorage(hardwareId, c.String("username-storage"))
	if err != nil {
		return cli.NewExitError(T("Failed to authorize storage to the hardware server instance: {{.Storage}}.\n{{.Error}}",
			map[string]interface{}{"Storage": c.String("username-storage"), "Error": err.Error()}), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, storageResult)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully authorized storage: {{.Storage}} was added.",
		map[string]interface{}{"Storage": c.String("username-storage")}))

	return nil
}

func HardwareAuthorizeStorageMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "authorize-storage",
		Description: T("Authorize File and Block Storage to a Hardware Server"),
		Usage: T(`${COMMAND_NAME} sl hardware authorize-storage [OPTIONS] IDENTIFIER
	
EXAMPLE:
   ${COMMAND_NAME} sl hardware authorize-storage --username-storage SL01SL30-37 1234567
   Authorize File and Block Storage to a Hardware Server.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u, username-storage",
				Usage: T("The storage username to be added to the hardware server."),
			},
			metadata.OutputFlag(),
		},
	}
}
