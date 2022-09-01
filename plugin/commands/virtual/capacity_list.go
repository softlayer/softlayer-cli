package virtual

import (
	"fmt"
	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityListCommand struct {
	*metadata.SoftlayerCommand
	Command        *cobra.Command
	VirtualServerManager managers.VirtualServerManager
}

func NewCapacityListCommand(sl *metadata.SoftlayerCommand) (cmd *CapacityListCommand) {
	thisCmd := &CapacityListCommand{
		SoftlayerCommand: sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capacity-list",
		Short: T("List Reserved Capacity groups."),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CapacityListCommand) Run(args []string) error {
	capacities, err := cmd.VirtualServerManager.CapacityList(utils.EMPTY_STRING)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Failed to get virtual Reserved capacity groups on your account.\n")
	}
	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Capacity"), T("Flavor"), T("Location"), T("Created")})
	for _, capacity := range capacities {
		available := utils.FormatUIntPointer(capacity.AvailableInstanceCount)

		billingDescription := utils.EMPTY_STRING
		if capacity.Instances[0].BillingItem != nil {
			billingDescription = utils.FormatStringPointer(capacity.Instances[0].BillingItem.Description)
		}
		table.Add(utils.FormatIntPointer(capacity.Id),
			utils.FormatStringPointer(capacity.Name),
			fmt.Sprintf("%s%s", available, " available"),
			billingDescription,
			utils.FormatStringPointer(capacity.BackendRouter.Hostname),
			utils.FormatSLTimePointer(capacity.CreateDate))
	}
	table.Print()
	return nil
}

