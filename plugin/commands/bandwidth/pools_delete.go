package bandwidth

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DeleteCommand struct {
	*metadata.SoftlayerCommand
	BandwidthManager managers.BandwidthManager
	Command          *cobra.Command
}

func NewDeleteCommand(sl *metadata.SoftlayerCommand) (cmd *DeleteCommand) {
	thisCmd := &DeleteCommand{
		SoftlayerCommand: sl,
		BandwidthManager: managers.NewBandwidthManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "pools-delete " + T("IDENTIFIER"),
		Short: T("Delete bandwidth pool. "),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DeleteCommand) Run(args []string) error {
	bandwidthPoolId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Bandwidth Pool ID")
	}

	err = cmd.BandwidthManager.DeletePool(bandwidthPoolId)
	// if err != nil {
	// 	return slErrors.NewAPIError(T("Failed to delete bandwidth with Id: {{.bandwidthPoolId}}.\n", map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}), err.Error(), 2)

	// }
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return slErrors.NewAPIError(T("Unable to find Bandwidth pool with ID {{.bandwidthPoolId}}.\n", map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}), err.Error(), 0)
		}
		return slErrors.NewAPIError(T("Failed to delete bandwidth with Id: {{.bandwidthPoolId}}.\n", map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}), err.Error(), 2)

	}
	cmd.UI.Ok()
	cmd.UI.Print(T("BandwidthPool associated with Id {{.bandwidthPoolId}} was deleted.", map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}))
	return nil
}
