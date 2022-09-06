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

type ResumeCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewResumeCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *ResumeCommand) {
	return &ResumeCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *ResumeCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will resume virtual server instance: {{.VsId}}. Continue?", map[string]interface{}{"VsId": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	err = cmd.VirtualServerManager.ResumeInstance(vsID)
	if err != nil {
		return cli.NewExitError(T("Failed to resume virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Virtual server instance: {{.VsId}} was resumed.", map[string]interface{}{"VsId": vsID}))
	return nil
}

func VSResumeMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "resume",
		Description: T("Resume a paused virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs resume IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs resume 12345678
   This command resumes virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
