package email

import (
	"bytes"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	EmailManager managers.EmailManager
	Command      *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		EmailManager:     managers.NewEmailManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("Lists Email Delivery Service."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := "mask[vendor,type]"
	emailList, err := cmd.EmailManager.GetNetworkMessageDeliveryAccounts(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Network Message Delivery Accounts."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{
		T("Name"),
		T("Value"),
	})

	bufEmail := new(bytes.Buffer)
	bufOverview := new(bytes.Buffer)
	bufStatistics := new(bytes.Buffer)

	emailTable := terminal.NewTable(bufEmail, []string{
		T("Id"),
		T("Username"),
		T("Description"),
		T("Vendor"),
	})

	// overviewTable := terminal.NewTable(bufOverview, []string{
	// 	T("Credit allowed"),
	// 	T("Credits remain"),
	// 	T("Credits overage"),
	// 	T("Credits used"),
	// 	T("Package"),
	// 	T("Reputation"),
	// 	T("Requests"),
	// })

	statisticsTable := terminal.NewTable(bufStatistics, []string{
		T("Delivered"),
		T("Requests"),
		T("Bounces"),
		T("Opens"),
		T("Clicks"),
		T("Spam reports"),
	})

	for _, email := range emailList {
		emailTable.Add(
			utils.FormatIntPointer(email.Id),
			utils.FormatStringPointer(email.Username),
			utils.FormatStringPointer(email.Type.Description),
			utils.FormatStringPointer(email.Vendor.KeyName),
		)

		// accountOverview, err := cmd.EmailManager.GetAccountOverview(*email.Id)
		// if err != nil {
		// 	return errors.NewAPIError(T("Failed to get Account Overview."), err.Error(), 2)
		// }
		// PrintAccountOverview(accountOverview, overviewTable, cmd.UI, outputFormat)

		statistics, err := cmd.EmailManager.GetStatistics(*email.Id)
		if err != nil {
			return errors.NewAPIError(T("Failed to get Statistics."), err.Error(), 2)
		}
		for _, statistic := range statistics {
			PrintStatistics(statistic, statisticsTable, cmd.UI, outputFormat)
		}
	}

	utils.PrintTable(cmd.UI, emailTable, outputFormat)

	table.Add(
		"Email information",
		bufEmail.String(),
	)
	table.Add(
		"Email overview",
		bufOverview.String(),
	)
	table.Add(
		"Statistics",
		bufStatistics.String(),
	)

	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}
/*

func PrintAccountOverview(accountOverview datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Account, overviewTable terminal.Table, ui terminal.UI, outputFormat string) {

	overviewTable.Add(
		utils.FormatIntPointer(accountOverview.CreditsAllowed),
		utils.FormatIntPointer(accountOverview.CreditsRemain),
		utils.FormatIntPointer(accountOverview.CreditsOverage),
		utils.FormatIntPointer(accountOverview.CreditsUsed),
		utils.FormatStringPointer(accountOverview.Package),
		utils.FormatIntPointer(accountOverview.Reputation),
		utils.FormatIntPointer(accountOverview.Requests),
	)
	utils.PrintTable(ui, overviewTable, outputFormat)
}

*/
func PrintStatistics(statistic datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Statistics, statisticsTable terminal.Table, ui terminal.UI, outputFormat string) {
	statisticsTable.Add(
		utils.FormatIntPointer(statistic.Delivered),
		utils.FormatIntPointer(statistic.Requests),
		utils.FormatIntPointer(statistic.Bounces),
		utils.FormatIntPointer(statistic.Opens),
		utils.FormatIntPointer(statistic.Clicks),
		utils.FormatIntPointer(statistic.SpamReports),
	)
	utils.PrintTable(ui, statisticsTable, outputFormat)
}
