package dns

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ZoneDeleteCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewZoneDeleteCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *ZoneDeleteCommand) {
	return &ZoneDeleteCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *ZoneDeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	zoneName := c.Args()[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}
	err = cmd.DNSManager.DeleteZone(zoneID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find zone with ID: {{.ZoneID}}.\n", map[string]interface{}{"ZoneID": zoneID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to delete zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Zone {{.Zone}} was deleted.", map[string]interface{}{"Zone": zoneName}))
	return nil
}

func DnsZoneDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "dns",
		Name:        "zone-delete",
		Description: T("Delete a zone"),
		Usage: T(`${COMMAND_NAME} sl dns zone-delete ZONE

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-delete ibm.com 
   This command deletes a zone that is named ibm.com.`),
	}
}
