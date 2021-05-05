package security

import (
	"io/ioutil"
	"strconv"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CertEditCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewCertEditCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *CertEditCommand) {
	return &CertEditCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *CertEditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	certID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("SSL certificate ID")
	}
	template := datatypes.Security_Certificate{Id: sl.Int(certID)}
	if c.IsSet("crt") {
		crtFile := c.String("crt")
		certificate, err := ioutil.ReadFile(crtFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read certificate file: {{.File}}.\n",
				map[string]interface{}{"File": crtFile})+err.Error(), 2)
		}
		template.Certificate = sl.String(string(certificate))
	}
	if c.IsSet("key") {
		keyFile := c.String("key")
		key, err := ioutil.ReadFile(keyFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read private key file: {{.File}}.\n",
				map[string]interface{}{"File": keyFile})+err.Error(), 2)
		}
		template.PrivateKey = sl.String(string(key))
	}
	if c.IsSet("icc") {
		iccFile := c.String("icc")
		icc, err := ioutil.ReadFile(iccFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read intermediate certificate file: {{.File}}.\n",
				map[string]interface{}{"File": iccFile})+err.Error(), 2)
		}
		template.IntermediateCertificate = sl.String(string(icc))
	}
	if c.IsSet("csr") {
		csrFile := c.String("csr")
		csr, err := ioutil.ReadFile(csrFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read certificate signing request file: {{.File}}.\n",
				map[string]interface{}{"File": csrFile})+err.Error(), 2)
		}
		template.CertificateSigningRequest = sl.String(string(csr))
	}
	if c.IsSet("notes") {
		template.Notes = sl.String(c.String("notes"))
	}
	err = cmd.SecurityManager.EditCertificate(template)
	if err != nil {
		return cli.NewExitError(T("Failed to edit SSL certificate: {{.ID}}.\n",
			map[string]interface{}{"ID": certID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate {{.ID}} was updated.", map[string]interface{}{"ID": certID}))
	return nil
}
