package dns

import (
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordEditCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
	ByRecord   string
	ById       int
	Data       string
	Ttl        int
}

func NewRecordEditCommand(sl *metadata.SoftlayerCommand) *RecordEditCommand {
	thisCmd := &RecordEditCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "record-edit " + T("ZONE"),
		Short: T("Update resource records in a zone"),
		Long: T(`${COMMAND_NAME} sl dns record-edit ZONE [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl dns record-edit ibm.com --by-id 12345678 --data 127.0.0.2 --ttl 3600
   This command edits records under the zone: ibm.com, whose ID is 12345678, and sets its data to "127.0.0.2" and ttl to 3600.
   ${COMMAND_NAME} sl dns record-edit ibm.com --by-record kibana --ttl 3600
   This command edits records under the zone: ibm.com, whose host is "kibana", and sets their ttl all to 3600.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.ByRecord, "by-record", "", T("Edit by host record, such as www"))
	cobraCmd.Flags().IntVar(&thisCmd.ById, "by-id", 0, T("Edit a single record by its ID"))
	cobraCmd.Flags().StringVar(&thisCmd.Data, "data", "", T("Record data, such as an IP address"))
	cobraCmd.Flags().IntVar(&thisCmd.Ttl, "ttl", 0, T("TTL(Time-To-Live) in seconds, such as: 86400. The default is: 7200"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RecordEditCommand) Run(args []string) error {

	zone := args[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zone)
	if err != nil {
		return errors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zone}), err.Error(), 2)
	}
	records, err := cmd.DNSManager.ListResourceRecords(zoneID, "", cmd.ByRecord, "", 0, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to list resource records under zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zone}), err.Error(), 2)
	}

	if len(records) == 0 {
		cmd.UI.Print(T("No record is found"))
		return nil
	}

	var multiErrors []error
	for _, record := range records {
		if cmd.ById != 0 {
			if record.Id != nil && cmd.ById != *record.Id {
				continue
			}
		}
		if cmd.Data != "" {
			record.Data = sl.String(cmd.Data)
		}
		if cmd.Ttl != 0 {
			record.Ttl = sl.Int(cmd.Ttl)
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
