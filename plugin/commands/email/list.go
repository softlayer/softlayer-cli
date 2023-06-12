package email

import (

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

	mask := "mask[vendor,type,billingItem[id,description]]"
	emailList, err := cmd.EmailManager.GetNetworkMessageDeliveryAccounts(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Network Message Delivery Accounts."), err.Error(), 2)
	}

	emailTable := cmd.UI.Table([]string{
		T("Id"),
		T("Username"),
		T("Plan"),
		T("Created"),
		T("Modified"),
		T("SMTP"),
	})

	for _, email := range emailList {
		emailTable.Add(
			utils.FormatIntPointer(email.Id),
			utils.FormatStringPointer(email.Username),
			utils.FormatStringPointer(email.BillingItem.Description),
			utils.FormatSLTimePointer(email.CreateDate),
			utils.FormatSLTimePointer(email.ModifyDate),
			utils.FormatStringPointer(email.SmtpAccess),		
		)
	}

	utils.PrintTable(cmd.UI, emailTable, outputFormat)
	return nil
}
