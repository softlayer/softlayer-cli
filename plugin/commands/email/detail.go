package email

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI           terminal.UI
	EmailManager managers.EmailManager
}

func NewDetailCommand(ui terminal.UI, emailManager managers.EmailManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:           ui,
		EmailManager: emailManager,
	}
}

func DetailMetaData() cli.Command {
	return cli.Command{
		Category:    "email",
		Name:        "detail",
		Description: T("Display details for a specified email."),
		Usage:       T(`${COMMAND_NAME} sl email detail IDENTIFIER [OPTIONS]`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	emailID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Email ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}
	mask := "mask[emailAddress,type,billingItem,vendor]"
	email, err := cmd.EmailManager.GetInstance(emailID, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get the email {{.emailID}}. ", map[string]interface{}{"emailID": emailID})+err.Error(), 2)
	}
	table := cmd.UI.Table([]string{
		T("Name"),
		T("Value"),
	})

	bufStatistics := new(bytes.Buffer)
	statisticsTable := terminal.NewTable(bufStatistics, []string{
		T("Delivered"),
		T("Requests"),
		T("Bounces"),
		T("Opens"),
		T("Clicks"),
		T("Spam reports"),
	})

	table.Add("Id", utils.FormatIntPointerName(email.Id))
	table.Add("Username", utils.FormatStringPointer(email.Username))
	table.Add("Email address", utils.FormatStringPointer(email.EmailAddress))
	table.Add("Create date", utils.FormatSLTimePointer(email.CreateDate))
	table.Add("Category code", utils.FormatStringPointer(email.BillingItem.CategoryCode))
	table.Add("Description", utils.FormatStringPointer(email.BillingItem.Description))
	table.Add("Type description", utils.FormatStringPointer(email.Type.Description))
	table.Add("Type", utils.FormatStringPointer(email.Type.KeyName))
	table.Add("Vendor", utils.FormatStringPointer(email.Vendor.KeyName))

	statistics, err := cmd.EmailManager.GetStatistics(*email.Id)
	if err != nil {
		return cli.NewExitError(T("Failed to get Statistics.")+err.Error(), 2)
	}
	for _, statistic := range statistics {
		PrintStatistics(statistic, statisticsTable, cmd.UI, outputFormat)
	}

	table.Add(
		"Statistics",
		bufStatistics.String(),
	)
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
