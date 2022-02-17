package dns

import (
	"errors"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordEditCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewRecordEditCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *RecordEditCommand) {
	return &RecordEditCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *RecordEditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	zone := c.Args()[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zone)
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zone})+err.Error(), 2)
	}
	records, err := cmd.DNSManager.ListResourceRecords(zoneID, "", c.String("by-record"), "", 0, "")
	if err != nil {
		return cli.NewExitError(T("Failed to list resource records under zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zone})+err.Error(), 2)
	}

	if len(records) == 0 {
		cmd.UI.Print(T("No record is found"))
		return nil
	}

	var multiErrors []error
	for _, record := range records {
		if c.IsSet("by-id") {
			if record.Id != nil && c.Int("by-id") != *record.Id {
				continue
			}
		}
		if c.IsSet("data") {
			record.Data = sl.String(c.String("data"))
		}
		if c.IsSet("ttl") {
			record.Ttl = sl.Int(c.Int("ttl"))
		}
		err = cmd.DNSManager.EditResourceRecord(record)
		if err != nil {
			newError := errors.New(T("Failed to update resource record {{.RecordID}} under zone {{.Zone}}.\n{{.ErrorMessage}}",
				map[string]interface{}{"RecordID": utils.IntPointertoInt(record.Id), "Zone": zone, "ErrorMessage": err.Error()}))
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Updated resource record under zone {{.Zone}}: ID={{.ID}}, type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.",
				map[string]interface{}{"Zone": zone, "ID": utils.IntPointertoInt(record.Id), "RecordType": utils.StringPointertoString(record.Type), "Host": utils.StringPointertoString(record.Host), "Data": utils.StringPointertoString(record.Data), "Ttl": utils.IntPointertoInt(record.Ttl)}))
		}
	}
	if len(multiErrors) > 0 {
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
	}
	return nil
}

func DnsRecordEditMetaData() cli.Command {
	return cli.Command{
		Category:    "dns",
		Name:        "record-edit",
		Description: T("Update resource records in a zone"),
		Usage: T(`${COMMAND_NAME} sl dns record-edit ZONE [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl dns record-edit ibm.com --by-id 12345678 --data 127.0.0.2 --ttl 3600
   This command edits records under the zone: ibm.com, whose ID is 12345678, and sets its data to "127.0.0.2" and ttl to 3600.
   ${COMMAND_NAME} sl dns record-edit ibm.com --by-record kibana --ttl 3600
   This command edits records under the zone: ibm.com, whose host is "kibana", and sets their ttl all to 3600.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "by-record",
				Usage: T("Edit by host record, such as www"),
			},
			cli.IntFlag{
				Name:  "by-id",
				Usage: T("Edit a single record by its ID"),
			},
			cli.StringFlag{
				Name:  "data",
				Usage: T("Record data, such as an IP address"),
			},
			cli.IntFlag{
				Name:  "ttl",
				Usage: T("TTL(Time-To-Live) in seconds, such as: 86400"),
			},
		},
	}
}
