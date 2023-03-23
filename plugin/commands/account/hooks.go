package account

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type HooksCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewHooksCommand(sl *metadata.SoftlayerCommand) *HooksCommand {
	thisCmd := &HooksCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "hooks",
		Short: T("Show all Provisioning Scripts."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *HooksCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	provisioningHooks, err := cmd.AccountManager.GetPostProvisioningHooks("")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Provisioning Hooks."), err.Error(), 2)
	}
	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Uri")})
	for _, hook := range provisioningHooks {
		table.Add(
			utils.FormatIntPointer(hook.Id),
			utils.FormatStringPointer(hook.Name),
			utils.FormatStringPointer(hook.Uri),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
