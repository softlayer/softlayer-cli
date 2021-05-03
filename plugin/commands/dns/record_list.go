package dns

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
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
