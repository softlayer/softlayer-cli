package licenses

import (
	"strings"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CancelItemCommand struct {
	*metadata.SoftlayerCommand
	Command         *cobra.Command
	LicensesManager managers.LicensesManager
	Immediate       bool
}

func NewCancelItemCommand(sl *metadata.SoftlayerCommand) *CancelItemCommand {
	thisCmd := &CancelItemCommand{
		SoftlayerCommand: sl,
		LicensesManager:  managers.NewLicensesManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("IDENTIFIER"),
		Short: T("Cancel a license."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Immediate, "immediate", false, T("Immediate cancellation."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelItemCommand) Run(args []string) error {

	key := args[0]

	err := cmd.LicensesManager.CancelItem(key, cmd.Immediate)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return slErr.NewAPIError(T("Unable to find license with key: {{.key}}.", map[string]interface{}{"key": key}), err.Error(), 0)
		}
		return slErr.NewAPIError(T("Failed to cancel license: {{.key}}.", map[string]interface{}{"key": key}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("License: {{.key}} was cancelled.", map[string]interface{}{"key": key}))
	return nil
}
