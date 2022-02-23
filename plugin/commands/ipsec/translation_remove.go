package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type RemoveTranslationCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewRemoveTranslationCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *RemoveTranslationCommand) {
	return &RemoveTranslationCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *RemoveTranslationCommand) Run(c *cli.Context) error {
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
	_, err = cmd.IPSECManager.GetTranslation(contextId, translationId)
	if err != nil {
		return cli.NewExitError(T("Failed to get translation with ID {{.TransID}} from IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId})+err.Error(), 2)
	}
	err = cmd.IPSECManager.RemoveTranslation(contextId, translationId)
	if err != nil {
		return cli.NewExitError(T("Failed to remove translation with ID {{.TransID}} from IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Removed translation with ID {{.TransID}} from IPSec {{.ID}}.",
		map[string]interface{}{"TransID": translationId, "ID": contextId}))
	return nil
}

func IpsecTransRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "translation-remove",
		Description: T("Remove a translation entry from an IPSec"),
		Usage: T(`${COMMAND_NAME} sl ipsec translation-remove CONTEXT_ID TRANSLATION_ID 

  Remove a translation entry from an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{},
	}
}
