package email

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	EmailManager managers.EmailManager
	Command      *cobra.Command
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		EmailManager:     managers.NewEmailManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Display details for a specified email."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {

	emailID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Email ID")
	}

	outputFormat := cmd.GetOutputFlag()

	mask := "mask[emailAddress,type,billingItem,vendor]"
	email, err := cmd.EmailManager.GetInstance(emailID, mask)
	if err != nil {
		return slErr.NewAPIError(T("Failed to get the email {{.emailID}}. ", map[string]interface{}{"emailID": emailID}), err.Error(), 2)
	}
	table := cmd.UI.Table([]string{
		T("Name"),
		T("Value"),
	})

	//Commented these lines until we fix EmailManager.GetStatistics() method
	/*
		bufStatistics := new(bytes.Buffer)
		statisticsTable := terminal.NewTable(bufStatistics, []string{
			T("Delivered"),
			T("Requests"),
			T("Bounces"),
			T("Opens"),
			T("Clicks"),
			T("Spam reports"),
		})
	*/

	table.Add("Id", utils.FormatIntPointerName(email.Id))
	table.Add("Username", utils.FormatStringPointer(email.Username))
	table.Add("Email address", utils.FormatStringPointer(email.EmailAddress))
	table.Add("Create date", utils.FormatSLTimePointer(email.CreateDate))
	table.Add("Category code", utils.FormatStringPointer(email.BillingItem.CategoryCode))
	table.Add("Description", utils.FormatStringPointer(email.BillingItem.Description))
	table.Add("Type description", utils.FormatStringPointer(email.Type.Description))
	table.Add("Type", utils.FormatStringPointer(email.Type.KeyName))
	table.Add("Vendor", utils.FormatStringPointer(email.Vendor.KeyName))

	//Commented these lines until we fix EmailManager.GetStatistics() method
	/*
		statistics, err := cmd.EmailManager.GetStatistics(*email.Id)
		if err != nil {
			return slErr.NewAPIError(T("Failed to get Statistics."), err.Error(), 2)
		}
		for _, statistic := range statistics {
			PrintStatistics(statistic, statisticsTable, cmd.UI, outputFormat)
		}

		table.Add(
			"Statistics",
			bufStatistics.String(),
		)
	*/
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
