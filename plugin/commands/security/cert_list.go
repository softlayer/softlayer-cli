package security

import (
	"sort"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CertListCommand struct {
	*metadata.SoftlayerCommand
	SecurityManager managers.SecurityManager
	Command         *cobra.Command
	Status          string
	SortBy          string
}

func NewCertListCommand(sl *metadata.SoftlayerCommand) *CertListCommand {
	thisCmd := &CertListCommand{
		SoftlayerCommand: sl,
		SecurityManager:  managers.NewSecurityManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cert-list",
		Short: T("List SSL certificates on your account."),
		Long: T(`${COMMAND_NAME} sl security cert-list [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl security cert-list --status valid --sortby days_until_expire
	This command lists all valid certificates on current account and sort them by validity days.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Status, "status", "", T("Show certificates with this status, default is: all, options are: all,valid,expired"))
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by. Options are: id,common_name,days_until_expire,note"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CertListCommand) Run(args []string) error {
	status := cmd.Status
	if status != "" && status != "all" && status != "valid" && status != "expired" {
		return errors.NewInvalidUsageError(T("[--status] must be either all, valid or expired."))
	}
	sortby := cmd.SortBy
	if sortby != "" && sortby != "id" && sortby != "ID" && sortby != "common_name" && sortby != "days_until_expire" && sortby != "note" {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.",
			map[string]interface{}{"Column": sortby}))
	}

	outputFormat := cmd.GetOutputFlag()

	certs, err := cmd.SecurityManager.ListCertificates(status)
	if err != nil {
		return errors.NewAPIError(T("Failed to list SSL certificates on your account.\n"), err.Error(), 2)
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

	table := cmd.UI.Table([]string{T("ID"), T("common_name"), T("days_until_expire"), T("note")})
	for _, cert := range certs {
		table.Add(utils.FormatIntPointer(cert.Id),
			utils.FormatStringPointer(cert.CommonName),
			utils.FormatIntPointer(cert.ValidityDays),
			utils.FormatStringPointer(cert.Notes))
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
