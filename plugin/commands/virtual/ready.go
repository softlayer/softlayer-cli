package virtual

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReadyCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewReadyCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *ReadyCommand) {
	return &ReadyCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *ReadyCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	until := time.Now().Add(time.Duration(c.Int("wait")) * time.Second)
	ready, message, err := cmd.VirtualServerManager.InstanceIsReady(vsID, until)
	if err != nil {
		return cli.NewExitError(T("Failed to check virtual server instance {{.VsID}} is ready.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}
	if ready {
		cmd.UI.Print(T("Virtual server instance: {{.VsId}} is ready.", map[string]interface{}{"VsId": vsID}))
	} else {
		cmd.UI.Print(T("Not ready: {{.Message}}", map[string]interface{}{"Message": message}))
	}
	return nil
}

func VSReadyMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "ready",
		Description: T("Check if a virtual server instance is ready for use"),
		Usage: T(`${COMMAND_NAME} sl vs ready IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs ready 12345678 --wait 30
   This command checks virtual server instance with ID 12345678 status to see if it is ready for use continuously and waits up to 30 seconds.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "wait",
				Usage: T("Wait until the virtual server is finished provisioning for up to X seconds before returning"),
			},
		},
	}
}
