package virtual

import (
	"github.com/spf13/cobra"
	"strconv"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupDetailsCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewPlacementGroupDetailsCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupDetailsCommand) {
	thisCmd := &PlacementGroupDetailsCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "placementgroup-detail " + T("IDENTIFIER"),
		Short: T("Authorize File, Block and Portable Storage to a Virtual Server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupDetailsCommand) Run(args []string) error {

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Placement Group Virtual server ID")
	}
	placement, err := cmd.VirtualServerManager.GetPlacementGroupDetail(id)

	outputFormat := cmd.GetOutputFlag()
	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Backend Router"), T("Rule"), T("Guests"), T("Created")})

	tableGuest := cmd.UI.Table([]string{T("ID"), T("FQDN"), T("Primary IP"), T("Backend IP"), T("CPU"), T("Memory"),
		T("Provisioned"), T("Transaction")})

	var backendName string
	if placement.BackendRouter == nil {
		backendName = "-"
	} else {
		backendName = utils.FormatStringPointer(placement.BackendRouter.Hostname)
	}
	table.Add(utils.FormatIntPointer(placement.Id),
		utils.FormatStringPointer(placement.Name),
		backendName,
		utils.FormatStringPointer(placement.Rule.Name),
		utils.FormatUIntPointer(placement.GuestCount),
		utils.FormatSLTimePointer(placement.CreateDate))

	for _, guest := range placement.Guests {

		var transaction string
		if guest.ActiveTransaction == nil {
			transaction = "-"
		} else {
			transaction = utils.FormatStringPointer(guest.ActiveTransaction.TransactionStatus.Name)
		}
		tableGuest.Add(utils.FormatIntPointer(guest.Id),
			utils.FormatStringPointer(guest.FullyQualifiedDomainName),
			utils.FormatStringPointer(guest.PrimaryIpAddress),
			utils.FormatStringPointer(guest.PrimaryBackendIpAddress),
			utils.FormatStringPointer(guest.MaxCpuUnits),
			utils.FormatIntPointer(guest.MaxMemory),
			utils.FormatSLTimePointer(guest.ProvisionDate),
			transaction)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	utils.PrintTable(cmd.UI, tableGuest, outputFormat)
	return nil
}
