package autoscale

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	AutoScaleManager managers.AutoScaleManager
	Command          *cobra.Command
	ForceFlag        bool
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		AutoScaleManager: managers.NewAutoScaleManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "delete " + T("IDENTIFIER"),
		Short: T("Delete this group and destroy all members of it"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {

	autoScaleGroupId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will cancel the AutoScale Group instance: {{.autoScaleGroupId}} and all its members, this action cannot be undone. Continue?", map[string]interface{}{"autoScaleGroupId": autoScaleGroupId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	response, err := cmd.AutoScaleManager.Delete(autoScaleGroupId)
	if err != nil {
		return errors.NewAPIError(T("Failed to delete AutoScale Group."), err.Error(), 2)
	}

	if response {
		cmd.UI.Ok()
		cmd.UI.Print(T("Auto Scale Group was deleted successfully"))
	}
	return nil
}
