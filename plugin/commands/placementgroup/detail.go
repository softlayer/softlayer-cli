package placementgroup

import (
	"bytes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupDetailCommand struct {
	UI                terminal.UI
	PlaceGroupManager managers.PlaceGroupManager
}

func NewPlacementGroupDetailCommand(ui terminal.UI, placeGroupManager managers.PlaceGroupManager) (cmd *PlacementGroupDetailCommand) {
	return &PlacementGroupDetailCommand{
		UI:                ui,
		PlaceGroupManager: placeGroupManager,
	}
}

func (cmd *PlacementGroupDetailCommand) Run(c *cli.Context) error {
	placementGroupID := c.Int("id")
	if placementGroupID == 0 {
		return errors.NewMissingInputError("--id")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	PlaceGroup, err := cmd.PlaceGroupManager.GetObject(placementGroupID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get placement group: {{.ID}}.\n", map[string]interface{}{"ID": placementGroupID})+err.Error(), 2)
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

func PlacementGroupDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "placement-group",
		Name:        "detail",
		Description: T("View details of a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group detail (--id PLACEMENTGROUP_ID) [--output FORMAT]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the placement group. [required]"),
			},
			metadata.OutputFlag(),
		},
	}
}
