package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewCancelCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will cancel the IPSec: {{.ContextID}} and cannot be undone. Continue?", map[string]interface{}{"ContextID": contextId}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.IPSECManager.CancelTunnelContext(contextId, c.Bool("immediate"), c.String("reason"))
	if err != nil {
		return cli.NewExitError(T("Failed to cancel IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("IPSec {{.ContextID}} is cancelled.", map[string]interface{}{"ContextID": contextId}))
	return nil
}

func IpsecCancelMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "cancel",
		Description: T("Cancel a IPSec VPN tunnel context"),
		Usage:       T(`${COMMAND_NAME} sl ipsec cancel CONTEXT_ID [OPTIONS]`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the IPSec immediately instead of on the billing anniversary"),
			},
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			metadata.ForceFlag(),
		},
	}
}
