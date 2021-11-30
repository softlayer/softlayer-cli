package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"strconv"
)

type PlacementGroupDetailsCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPlacementGroupDetailsCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PlacementGroupDetailsCommand) {
	return &PlacementGroupDetailsCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PlacementGroupDetailsCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErrors.NewInvalidUsageError(T("This command requires one argument."))
	}
	id, err := strconv.Atoi(c.Args()[0])
	placement, err := cmd.VirtualServerManager.GetPlacementGroupDetail(id)
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Placement Group Virtual server ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, placement)
	}

	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Backend Router"),
		T("Rule"), T("Guests"), T("Created")})

	tableGuest := cmd.UI.Table([]string{T("ID"), T("FQDN"), T("Primary IP"),
		T("Backend IP"),
		T("CPU"), T("Memory"), T("Provisioned"), T("Transaction")})

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
	table.Print()
	tableGuest.Print()
	return nil
}
