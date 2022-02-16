package dns

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordListCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewRecordListCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *RecordListCommand) {
	return &RecordListCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *RecordListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	zoneName := c.Args()[0]

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}
	records, err := cmd.DNSManager.ListResourceRecords(zoneID, c.String("type"), c.String("record"), c.String("data"), c.Int("ttl"), "")
	if err != nil {
		return cli.NewExitError(T("Failed to list resource records under zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, records)
	}

	table := cmd.UI.Table([]string{T("ID"), T("host"), T("type"), T("ttl"), T("data")})
	for _, record := range records {
		table.Add(utils.FormatIntPointer(record.Id),
			utils.FormatStringPointer(record.Host),
			strings.ToUpper(utils.FormatStringPointer(record.Type)),
			utils.FormatIntPointer(record.Ttl),
			utils.FormatStringPointer(record.Data))
	}
	table.Print()
	return nil
}

func DnsRecordListMetaData() cli.Command {
	return cli.Command{
		Category:    "dns",
		Name:        "record-list",
		Description: T("List all the resource records in a zone"),
		Usage: T(`${COMMAND_NAME} sl dns record-list ZONE [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl dns record-list ibm.com --record elasticsearch --type A --ttl 900
   This command lists all A records under the zone: ibm.com, and filters by host is elasticsearch and ttl is 900 seconds.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "data",
				Usage: T("Filter by record data, such as an IP address"),
			},
			cli.StringFlag{
				Name:  "record",
				Usage: T("Filter by host record, such as www"),
			},
			cli.IntFlag{
				Name:  "ttl",
				Usage: T("Filter by TTL(Time-To-Live) in seconds, such as 86400"),
			},
			cli.StringFlag{
				Name:  "type",
				Usage: T("Filter by record type, such as A or CNAME"),
			},
			metadata.OutputFlag(),
		},
	}
}
