package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AddTranslationCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewAddTranslationCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *AddTranslationCommand) {
	return &AddTranslationCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *AddTranslationCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	staticIp := c.String("s")
	if staticIp == "" {
		return errors.NewMissingInputError("-s|--static-ip")
	}
	remoteIp := c.String("r")
	if remoteIp == "" {
		return errors.NewMissingInputError("-r|--remote-ip")
	}
	_, err = cmd.IPSECManager.GetTunnelContext(contextId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId})+err.Error(), 2)
	}
	translation, err := cmd.IPSECManager.CreateTranslation(contextId, staticIp, remoteIp, c.String("n"))
	if err != nil {
		return cli.NewExitError(T("Failed to create translation for IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, translation)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Created translation from {{.StaticIP}} to {{.RemoteIP}} #{{.ID}}.",
		map[string]interface{}{"StaticIP": staticIp, "RemoteIP": remoteIp, "ID": *translation.Id}))
	return nil
}
