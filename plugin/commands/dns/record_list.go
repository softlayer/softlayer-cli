package dns

import (
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type RecordListCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
	Data       string
	Record     string
	Ttl        int
	Type       string
}

func NewRecordListCommand(sl *metadata.SoftlayerCommand) *RecordListCommand {
	thisCmd := &RecordListCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "record-list " + T("ZONE"),
		Short: T("List all the resource records in a zone"),
		Long: T(`${COMMAND_NAME} sl dns record-list ZONE [OPTIONS]
	
EXAMPLE:
	${COMMAND_NAME} sl dns record-list ibm.com --record elasticsearch --type A --ttl 900
	This command lists all A records under the zone: ibm.com, and filters by host is elasticsearch and ttl is 900 seconds.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Data, "data", "", T("Filter by record data, such as an IP address"))
	cobraCmd.Flags().StringVar(&thisCmd.Record, "record", "", T("Filter by host record, such as www"))
	cobraCmd.Flags().IntVar(&thisCmd.Ttl, "ttl", 0, T("TTL(Time-To-Live) in seconds, such as: 86400. The default is: 7200"))
	cobraCmd.Flags().StringVar(&thisCmd.Type, "type", "", T("Filter by record type, such as A or CNAME"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RecordListCommand) Run(args []string) error {
	zoneName := args[0]

	outputFormat := cmd.GetOutputFlag()

	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return errors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}
	records, err := cmd.DNSManager.ListResourceRecords(zoneID, cmd.Type, cmd.Record, cmd.Data, cmd.Ttl, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to list resource records under zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
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
