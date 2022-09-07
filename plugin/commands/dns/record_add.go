package dns

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordAddCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
	Ttl        int
}

func NewRecordAddCommand(sl *metadata.SoftlayerCommand) *RecordAddCommand {
	thisCmd := &RecordAddCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "record-add " + T("ZONE") + " " + T("RECORD") + " " + T("TYPE") + " " + T("DATA"),
		Short: T("Add resource record in a zone."),
		Long: T(`${COMMAND_NAME} sl dns record-add ZONE RECORD TYPE DATA [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl dns record-add ibm.com ftp A 127.0.0.1 --ttl 86400
	This command adds an A record to zone: ibm.com, its host is "ftp", data is "127.0.0.1" and ttl is 86400 seconds.`),
		Args: metadata.FourArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Ttl, "ttl", 0, T("TTL(Time-To-Live) in seconds, such as: 86400. The default is: 7200"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RecordAddCommand) Run(args []string) error {
	zone := args[0]
	host := args[1]
	recordType := args[2]
	data := args[3]
	ttl := 7200
	if cmd.Ttl != 0 {
		ttl = cmd.Ttl
	}

	outputFormat := cmd.GetOutputFlag()

	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zone)
	if err != nil {
		return errors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zone}), err.Error(), 2)
	}
	record, err := cmd.DNSManager.CreateResourceRecord(zoneID, host, recordType, data, ttl)
	if err != nil {
		return errors.NewAPIError(T("Failed to create resource record under zone {{.Zone}}: type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.\n",
			map[string]interface{}{"Zone": zone, "RecordType": recordType, "Host": host, "Data": data, "Ttl": ttl}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, record)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Created resource record under zone {{.Zone}}: ID={{.ID}}, type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.",
		map[string]interface{}{"Zone": zone, "ID": utils.IntPointertoInt(record.Id), "RecordType": utils.StringPointertoString(record.Type), "Host": utils.StringPointertoString(record.Host), "Data": utils.StringPointertoString(record.Data), "Ttl": utils.IntPointertoInt(record.Ttl)}))
	return nil
}
