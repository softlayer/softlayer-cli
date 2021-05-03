package security

import (
	"io/ioutil"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type CertAddCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewCertAddCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *CertAddCommand) {
	return &CertAddCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *CertAddCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	crtFile := c.String("crt")
	if crtFile == "" {
		return errors.NewMissingInputError("--crt")
	}
	keyFile := c.String("key")
	if keyFile == "" {
		return errors.NewMissingInputError("--key")
	}
	template := datatypes.Security_Certificate{}
	certificate, err := ioutil.ReadFile(crtFile) // #nosec
	if err != nil {
		return cli.NewExitError(T("Failed to read certificate file: {{.File}}.\n",
			map[string]interface{}{"File": crtFile})+err.Error(), 2)
	}
	template.Certificate = sl.String(string(certificate))
	key, err := ioutil.ReadFile(keyFile) // #nosec
	if err != nil {
		return cli.NewExitError(T("Failed to read private key file: {{.File}}.\n",
			map[string]interface{}{"File": keyFile})+err.Error(), 2)
	}
	template.PrivateKey = sl.String(string(key))
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
	cert, err := cmd.SecurityManager.AddCertificate(template)
	if err != nil {
		return cli.NewExitError(T("Failed to add certificate.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, cert)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("SSL certificate for {{.CommonName}} was added.", map[string]interface{}{"CommonName": utils.StringPointertoString(cert.CommonName)}))
	return nil
}
