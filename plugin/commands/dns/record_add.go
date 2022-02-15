package dns

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordAddCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewRecordAddCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *RecordAddCommand) {
	return &RecordAddCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *RecordAddCommand) Run(c *cli.Context) error {
	if c.NArg() != 4 {
		return errors.NewInvalidUsageError(T("This command requires four arguments."))
	}
	zone := c.Args()[0]
	host := c.Args()[1]
	recordType := c.Args()[2]
	data := c.Args()[3]
	ttl := 7200
	if c.IsSet("ttl") {
		ttl = c.Int("ttl")
	}
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zone)
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zone})+err.Error(), 2)
	}
	record, err := cmd.DNSManager.CreateResourceRecord(zoneID, host, recordType, data, ttl)
	if err != nil {
		return cli.NewExitError(T("Failed to create resource record under zone {{.Zone}}: type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.\n",
			map[string]interface{}{"Zone": zone, "RecordType": recordType, "Host": host, "Data": data, "Ttl": ttl})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, record)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Created resource record under zone {{.Zone}}: ID={{.ID}}, type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.",
		map[string]interface{}{"Zone": zone, "ID": utils.IntPointertoInt(record.Id), "RecordType": utils.StringPointertoString(record.Type), "Host": utils.StringPointertoString(record.Host), "Data": utils.StringPointertoString(record.Data), "Ttl": utils.IntPointertoInt(record.Ttl)}))
	return nil
}

func DnsRecordAddMetaData() cli.Command {
	return cli.Command{
		Category:    "dns",
		Name:        "record-add",
		Description: T("Add resource record in a zone"),
		Usage: T(`${COMMAND_NAME} sl dns record-add ZONE RECORD TYPE DATA [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns record-add ibm.com ftp A 127.0.0.1 --ttl 86400
   This command adds an A record to zone: ibm.com, its host is "ftp", data is "127.0.0.1" and ttl is 86400 seconds.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "ttl",
				Usage: T("TTL(Time-To-Live) in seconds, such as: 86400. The default is: 7200"),
			},
			metadata.OutputFlag(),
		},
	}
}
