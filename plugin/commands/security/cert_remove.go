package security

import (
	"strconv"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CertRemoveCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewCertRemoveCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *CertRemoveCommand) {
	return &CertRemoveCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *CertRemoveCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	certID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will remove SSL certificate: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": certID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.SecurityManager.RemoveCertificate(certID)
	if err != nil {
		return cli.NewExitError(T("Failed to remove SSL certificate: {{.ID}}.\n", map[string]interface{}{"ID": certID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate {{.ID}} was removed.", map[string]interface{}{"ID": certID}))
	return nil
}

func SecuritySSLCertRemove() cli.Command {
	return cli.Command{
		Category:    "security",
		Name:        "cert-remove",
		Description: T("Remove SSL certificate"),
		Usage: T(`${COMMAND_NAME} sl security cert-remove IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-remove 12345678 
   This command removes certificate with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
