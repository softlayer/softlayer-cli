package security

import (
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CertListCommand struct {
	UI              terminal.UI
	SecurityManager managers.SecurityManager
}

func NewCertListCommand(ui terminal.UI, securityManager managers.SecurityManager) (cmd *CertListCommand) {
	return &CertListCommand{
		UI:              ui,
		SecurityManager: securityManager,
	}
}

func (cmd *CertListCommand) Run(c *cli.Context) error {
	status := c.String("status")
	if status != "" && status != "all" && status != "valid" && status != "expired" {
		return errors.NewInvalidUsageError(T("[--status] must be either all, valid or expired."))
	}
	sortby := c.String("sortby")
	if sortby != "" && sortby != "id" && sortby != "ID" && sortby != "common_name" && sortby != "days_until_expire" && sortby != "note" {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.",
			map[string]interface{}{"Column": sortby}))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	certs, err := cmd.SecurityManager.ListCertificates(status)
	if err != nil {
		return cli.NewExitError(T("Failed to list SSL certificates on your account.\n")+err.Error(), 2)
	}

	if sortby == "id" || sortby == "ID" {
		sort.Sort(utils.CertById(certs))
	} else if sortby == "common_name" {
		sort.Sort(utils.CertByCommonName(certs))
	} else if sortby == "days_until_expire" {
		sort.Sort(utils.CertByValidityDays(certs))
	} else if sortby == "note" {
		sort.Sort(utils.CertByNotes(certs))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, certs)
	}

	table := cmd.UI.Table([]string{T("ID"), T("common_name"), T("days_until_expire"), T("note")})
	for _, cert := range certs {
		table.Add(utils.FormatIntPointer(cert.Id),
			utils.FormatStringPointer(cert.CommonName),
			utils.FormatIntPointer(cert.ValidityDays),
			utils.FormatStringPointer(cert.Notes))
	}
	table.Print()
	return nil
}
