package placementgroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type PlacementGroupDeleteCommand struct {
	UI                terminal.UI
	PlaceGroupManager managers.PlaceGroupManager
	VMManager         managers.VirtualServerManager
}

func NewPlacementGroupDeleteCommand(ui terminal.UI, placeGroupManager managers.PlaceGroupManager, vmManager managers.VirtualServerManager) (cmd *PlacementGroupDeleteCommand) {
	return &PlacementGroupDeleteCommand{
		UI:                ui,
		PlaceGroupManager: placeGroupManager,
		VMManager:         vmManager,
	}
}

func (cmd *PlacementGroupDeleteCommand) Run(c *cli.Context) error {
	placementGroupID := c.Int("id")
	if placementGroupID == 0 {
		return errors.NewMissingInputError("--id")
	}

	//if c.IsSet("purge") {
	//	placementGroup, err := cmd.PlaceGroupManager.GetObject(placementGroupID, "")
	//	if err != nil {
	//		return cli.NewExitError(T("Failed to get placement group: {{.ID}}\n", map[string]interface{}{"ID": placementGroupID})+err.Error(), 2)
	//	}
	//	if len(placementGroup.Guests) < 1 {
	//		cmd.UI.Print(T("No virtual server was found in placement group {{.ID}}.", map[string]interface{}{"ID": placementGroupID}))
	//	} else {
	//		if !c.IsSet("f") {
	//			guestTable := cmd.UI.Table([]string{T("ID"), T("FQDN"), T("Primary IP"), T("Backend IP"), T("CPU"), T("Memory"), T("Provisioned")})
	//			for _, guest := range placementGroup.Guests {
	//				guestTable.Add(
	//					utils.FormatIntPointer(guest.Id),
	//					utils.FormatStringPointer(guest.FullyQualifiedDomainName),
	//					utils.FormatStringPointer(guest.PrimaryIpAddress),
	//					utils.FormatStringPointer(guest.PrimaryBackendIpAddress),
	//					utils.FormatIntPointer(guest.MaxCpu),
	//					utils.FormatIntPointer(guest.MaxMemory),
	//					utils.FormatSLTimePointer(guest.ProvisionDate),
	//				)
	//			}
	//			guestTable.Print()
	//			confirm, err := cmd.UI.Confirm(T("This will remove all the above virtual servers! Continue?"))
	//			if err != nil {
	//				return cli.NewExitError(err.Error(), 1)
	//			}
	//			if !confirm {
	//				cmd.UI.Print(T("Aborted."))
	//				return nil
	//			}
	//		}
	//	}
	//	for _, guest := range placementGroup.Guests {
	//		cmd.UI.Print(T("Deleting guest: {{.Name}}.", map[string]interface{}{"Name": guest.FullyQualifiedDomainName}))
	//		cmd.VMManager.CancelInstance(*guest.Id)
	//	}
	//}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will remove placement group: {{.ID}} and cannot be undone. Continue?", map[string]interface{}{"ID": placementGroupID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	_, err := cmd.PlaceGroupManager.Delete(placementGroupID)
	if err != nil {
		return cli.NewExitError(T("Failed to remove placement group: {{.ID}}.\n", map[string]interface{}{"ID": placementGroupID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Placement group {{.ID}} was removed.", map[string]interface{}{"ID": placementGroupID}))
	return nil
}

func PlacementGroupDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "placement-group",
		Name:        "delete",
		Description: T("Delete a placement group"),
		Usage:       "${COMMAND_NAME} sl placement-group delete (--id PLACEMENTGROUP_ID) [-f, --force]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the placement group. [required]"),
			},
			//cli.BoolFlag{   # tmp disable this option. because the placement can't be deleted if the VSI status is delete pending.
			//	Name:  "purge",
			//	Usage: T("Delete all guests in this placement group. The group itself can be deleted once all VMs are fully reclaimed"),
			//},
			metadata.ForceFlag(),
		},
	}
}
