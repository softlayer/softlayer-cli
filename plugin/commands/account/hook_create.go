package account

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type HookCreateCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
	Name           string
	Uri            string
}

func NewHookCreateCommand(sl *metadata.SoftlayerCommand) *HookCreateCommand {
	thisCmd := &HookCreateCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "hook-create",
		Short: T("Order/create a provisioning script."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "N", "", T("The name of the hook.  [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Uri, "uri", "U", "", T("The endpoint that the script will be downloaded.  [required]"))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("name")
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("uri")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *HookCreateCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	hookTemplate := datatypes.Provisioning_Hook{
		Name: sl.String(cmd.Name),
		Uri:  sl.String(cmd.Uri),
	}

	provisioningHook, err := cmd.AccountManager.CreateProvisioningScript(hookTemplate)
	if err != nil {
		return errors.NewAPIError(T("Failed to create Provisioning Hook."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("Id"), utils.FormatIntPointer(provisioningHook.Id))
	table.Add(T("Name"), utils.FormatStringPointer(provisioningHook.Name))
	table.Add(T("Created"), utils.FormatSLTimePointer(provisioningHook.CreateDate))
	table.Add(T("Uri"), utils.FormatStringPointer(provisioningHook.Uri))

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
