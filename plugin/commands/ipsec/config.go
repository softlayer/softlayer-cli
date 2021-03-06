package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ConfigCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewConfigCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *ConfigCommand) {
	return &ConfigCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *ConfigCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	err = cmd.IPSECManager.ApplyConfiguration(contextId)
	if err != nil {
		return cli.NewExitError(T("Failed to enqueue configuration request for IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Configuration request received for IPSec {{.ContextID}}.", map[string]interface{}{"ContextID": contextId}))
	return nil
}

func IpsecConfigMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "config",
		Description: T("Request configuration of a tunnel context"),
		Usage: T(`${COMMAND_NAME} sl ipsec config CONTEXT_ID [OPTIONS]

  Request configuration of a tunnel context.

  This action will update the advancedConfigurationFlag on the context
  instance and further modifications against the context will be prevented
  until all changes can be propagated to network devices.`),
		Flags: []cli.Flag{},
	}
}
