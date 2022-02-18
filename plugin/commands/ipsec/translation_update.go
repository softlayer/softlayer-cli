package ipsec

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

type UpdateTranslationCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewUpdateTranslationCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *UpdateTranslationCommand) {
	return &UpdateTranslationCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *UpdateTranslationCommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	args1 := c.Args()[1]
	translationId, err := strconv.Atoi(args1)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Translation ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.IPSECManager.UpdateTranslation(contextId, translationId, c.String("s"), c.String("r"), c.String("n"))
	if err != nil {
		return cli.NewExitError(T("Failed to update translation with ID {{.TransID}} in IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Updated translation with ID {{.TransID}} in IPSec {{.ID}}.",
		map[string]interface{}{"TransID": translationId, "ID": contextId}))
	return nil
}

func IpsecTransUpdataMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "translation-update",
		Description: T("Update an address translation for an IPSec"),
		Usage: T(`${COMMAND_NAME} sl ipsec translation-update CONTEXT_ID TRANSLATION_ID [OPTIONS]

  Update an address translation for an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,static-ip",
				Usage: T("Static IP address[required]"),
			},
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("Remote IP address[required]"),
			},
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Note"),
			},
			metadata.OutputFlag(),
		},
	}
}
