package autoscale

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	AutoScaleManager managers.AutoScaleManager
	Command          *cobra.Command
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		AutoScaleManager: managers.NewAutoScaleManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all Autoscale Groups on your account"),
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

	mask := "mask[id,name,status,minimumMemberCount,maximumMemberCount,virtualGuestMemberCount]"
	scaleGroups, err := cmd.AutoScaleManager.ListScaleGroups(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get scale groups."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Status"), T("Min/Max"), T("Running")})
	for _, scale := range scaleGroups {
		membercount := strconv.Itoa(*scale.MinimumMemberCount) + "/" + strconv.Itoa(*scale.MaximumMemberCount)
		table.Add(
			utils.FormatIntPointer(scale.Id),
			utils.FormatStringPointer(scale.Name),
			utils.FormatStringPointer(scale.Status.Name),
			membercount,
			utils.FormatUIntPointer(scale.VirtualGuestMemberCount),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
