package dns

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ZonePrintCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewZonePrintCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *ZonePrintCommand) {
	return &ZonePrintCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *ZonePrintCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	zoneName := c.Args()[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}
	fileContent, err := cmd.DNSManager.DumpZone(zoneID)
	if err != nil {
		return cli.NewExitError(T("Failed to dump content for zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}
	//TODO need to test on other platforms about line break
	lines := strings.Split(fileContent, "\\n")
	for _, line := range lines {
		cmd.UI.Print(strings.Trim(line, "\""))
	}
	return nil
}

func DnsZonePrintMetaData() cli.Command {
	return cli.Command{
		Category:    "dns",
		Name:        "zone-print",
		Description: T("Print zone and resource records in BIND format"),
		Usage: T(`${COMMAND_NAME} sl dns zone-print ZONE

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-print ibm.com
   This command prints zone that is named ibm.com, and in BIND format.`),
	}
}
