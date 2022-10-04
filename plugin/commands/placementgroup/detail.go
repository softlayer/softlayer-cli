package placementgroup

import (
	"bytes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupDetailCommand struct {
	*metadata.SoftlayerCommand
	PlaceGroupManager managers.PlaceGroupManager
	Command           *cobra.Command
	Id                int
}

func NewPlacementGroupDetailCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupDetailCommand) {
	thisCmd := &PlacementGroupDetailCommand{
		SoftlayerCommand:  sl,
		PlaceGroupManager: managers.NewPlaceGroupManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail",
		Short: T("View details of a placement group"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVar(&thisCmd.Id, "id", 0, T("ID for the placement group. [required]"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupDetailCommand) Run(args []string) error {
	placementGroupID := cmd.Id
	if placementGroupID == 0 {
		return errors.NewMissingInputError("--id")
	}

	outputFormat := cmd.GetOutputFlag()

	PlaceGroup, err := cmd.PlaceGroupManager.GetObject(placementGroupID, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get placement group: {{.ID}}.", map[string]interface{}{"ID": placementGroupID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, PlaceGroup)
	}

	cmd.UI.Ok()
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	backendRouter := "-"
	rule := "-"
	if PlaceGroup.BackendRouter != nil {
		backendRouter = utils.FormatStringPointer(PlaceGroup.BackendRouter.Hostname)
	}
	if PlaceGroup.Rule != nil {
		rule = utils.FormatStringPointer(PlaceGroup.Rule.Name)
	}
	table.Add(T("ID"), utils.FormatIntPointer(PlaceGroup.Id))
	table.Add(T("Name"), utils.FormatStringPointer(PlaceGroup.Name))
	table.Add(T("Backend Router"), backendRouter)
	table.Add(T("Rule"), rule)
	table.Add(T("Created"), utils.FormatSLTimePointer(PlaceGroup.CreateDate))

	if len(PlaceGroup.Guests) > 0 {
		buf := new(bytes.Buffer)
		guestTable := terminal.NewTable(buf, []string{T("ID"), T("FQDN"), T("Primary IP"), T("Backend IP"), T("CPU"), T("Memory"), T("Provisioned")})
		for _, guest := range PlaceGroup.Guests {
			guestTable.Add(
				utils.FormatIntPointer(guest.Id),
				utils.FormatStringPointer(guest.FullyQualifiedDomainName),
				utils.FormatStringPointer(guest.PrimaryIpAddress),
				utils.FormatStringPointer(guest.PrimaryBackendIpAddress),
				utils.FormatIntPointer(guest.MaxCpu),
				utils.FormatIntPointer(guest.MaxMemory),
				utils.FormatSLTimePointer(guest.ProvisionDate),
			)
		}
		guestTable.Print()
		table.Add(T("Guests: "), buf.String())
	} else {
		table.Add(T("Guests: "), T("No guest was found."))
	}
	table.Print()
	return nil
}
