package dedicatedhost

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CancelCommand struct {
	UI                   terminal.UI
	DedicatedHostManager managers.DedicatedHostManager
}

func NewCancelCommand(ui terminal.UI, dedicatedHostManager managers.DedicatedHostManager) (cmd *CancelCommand) {
	return &CancelCommand{
		UI:                   ui,
		DedicatedHostManager: dedicatedHostManager,
	}
}

func (cmd *CancelCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	HostID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Host ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This will cancel all virtual server instances in the dedicatedhost: {{.HostID}} and cannot be undone. Continue?", map[string]interface{}{"HostID": HostID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	listGuest, err := cmd.DedicatedHostManager.CancelGuests(HostID)
	if err != nil {
		return cli.NewExitError(T("Failed to cancel all guests in the dedicatedhost: {{.HostID}}.\n", map[string]interface{}{"HostID": HostID})+err.Error(), 2)
	}

	if len(listGuest) > 0 {
		table := cmd.UI.Table([]string{T("Id"), T("Server Name"), T("Status")})
		for _, guest := range listGuest {
			table.Add(utils.FormatIntPointer(&guest.Id), utils.FormatStringPointer(&guest.Fqdn), utils.FormatStringPointer(&guest.Status))
		}
		table.Print()
		return nil
	}

	return cli.NewExitError(T("There is not any guest into the dedicated host {{.ID}}.", map[string]interface{}{"ID": HostID}), 2)
}

func DedicatedhostCancelGuestsMetaData() cli.Command {
	return cli.Command{
		Category:    "dedicatedhost",
		Name:        "cancel-guests",
		Description: T("Cancel all virtual guests of the dedicated host immediately."),
		Usage:       "${COMMAND_NAME} sl dedicatedhost cancel-guests IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			metadata.ForceFlag(),
		},
	}
}
