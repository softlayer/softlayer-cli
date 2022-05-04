package email

import (
	"bytes"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI           terminal.UI
	EmailManager managers.EmailManager
}

func NewListCommand(ui terminal.UI, emailManager managers.EmailManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:           ui,
		EmailManager: emailManager,
	}
}

func ListMetaData() cli.Command {
	return cli.Command{
		Category:    "email",
		Name:        "list",
		Description: T("Lists Email Delivery Service."),
		Usage:       T(`${COMMAND_NAME} sl email list`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[vendor,type]"
	emailList, err := cmd.EmailManager.GetNetworkMessageDeliveryAccounts(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get Network Message Delivery Accounts.")+err.Error(), 2)
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

	for _, email := range emailList {
		emailTable.Add(
			utils.FormatIntPointer(email.Id),
			utils.FormatStringPointer(email.Username),
			utils.FormatStringPointer(email.Type.Description),
			utils.FormatStringPointer(email.Vendor.KeyName),
		)

		accountOverview, err := cmd.EmailManager.GetAccountOverview(*email.Id)
		if err != nil {
			return cli.NewExitError(T("Failed to get Account Overview.")+err.Error(), 2)
		}
		PrintAccountOverview(accountOverview, bufOverview, cmd.UI, outputFormat)

		statistics, err := cmd.EmailManager.GetStatistics(*email.Id)
		if err != nil {
			return cli.NewExitError(T("Failed to get Statistics.")+err.Error(), 2)
		}
		for _, statistic := range statistics {
			PrintStatistics(statistic, bufStatistics, cmd.UI, outputFormat)
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

func PrintAccountOverview(accountOverview datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Account_Overview, bufOverview *bytes.Buffer, ui terminal.UI, outputFormat string) {
	overviewTable := terminal.NewTable(bufOverview, []string{
		T("Credit allowed"),
		T("Credits remain"),
		T("Credits overage"),
		T("Credits used"),
		T("Package"),
		T("Reputation"),
		T("Requests"),
	})

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

func PrintStatistics(statistic datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Statistics, bufStatistics *bytes.Buffer, ui terminal.UI, outputFormat string) {
	statisticsTable := terminal.NewTable(bufStatistics, []string{
		T("Delivered"),
		T("Requests"),
		T("Bounces"),
		T("Opens"),
		T("Clicks"),
		T("Spam reports"),
	})

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
