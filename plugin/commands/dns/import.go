package dns

import (
	"os"
	"strings"
	"bytes"

	"github.com/miekg/dns"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ImportCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
	DryRun     bool
}

func NewImportCommand(sl *metadata.SoftlayerCommand) *ImportCommand {
	thisCmd := &ImportCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "import " + T("ZONEFILE"),
		Short: T("Import a zone based off a BIND zone file"),
		Long: T(`${COMMAND_NAME} sl dns import ZONEFILE [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns import ~/ibm.com.txt
   This command imports zone and its resource records from file: ~/ibm.com.txt.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.DryRun, "dry-run", false, T("Don't actually create records"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ImportCommand) Run(args []string) error {
	bytes, err := os.ReadFile(args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to read file: {{.FilePath}}.\n", map[string]interface{}{"FilePath": args[0]}), err.Error(), 2)
	}
	zone, records, err := parseFileContent(bytes, args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to parse file.\n"), err.Error(), 2)
	}

	if cmd.DryRun {
		cmd.UI.Ok()
		return nil
	}

	dnsDomain, err := cmd.DNSManager.CreateZone(*zone.Name) //parseFileContent can guarantee zone.name is not nil
	if err != nil {

		return errors.NewAPIError(T("Failed to create zone: {{.ZoneName}}.\n", map[string]interface{}{"ZoneName": *zone.Name}), err.Error(), 2)
	}

	cmd.UI.Print(T("Zone {{.Zone}} was created.", map[string]interface{}{"Zone": utils.StringPointertoString(dnsDomain.Name)}))

	var multiErrors []error
	for _, record := range records {
		record.DomainId = dnsDomain.Id
		rr, err := cmd.DNSManager.ResourceRecordCreate(record)
		if err != nil {
			newError := errors.New(T("Failed to create resource record under zone {{.Zone}}: type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.\n{{.ErrorMessage}}",
				map[string]interface{}{"Zone": zone.Name, "RecordType": *rr.Type, "Host": *rr.Host, "Data": *rr.Data, "Ttl": *rr.Ttl, "ErrorMessage": err.Error()}))
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Print(T("Created resource record under zone {{.Zone}}: ID={{.ID}}, type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.",
				map[string]interface{}{"Zone": zone.Name, "ID": *rr.Id, "RecordType": *rr.Type, "Host": *rr.Host, "Data": *rr.Data, "Ttl": *rr.Ttl}))
		}
	}

	if len(multiErrors) > 0 {
		return errors.CollapseErrors(multiErrors)
	}
	return nil
}

func parseFileContent(content []byte, filename string) (datatypes.Dns_Domain, []datatypes.Dns_Domain_ResourceRecord, error) {
	zone := datatypes.Dns_Domain{}
	records := []datatypes.Dns_Domain_ResourceRecord{}
	// Top Level Domain: AKA $ORIGIN
	tld := ""
	zoneParser := dns.NewZoneParser(bytes.NewReader(content), "", filename)

	for rr, ok := zoneParser.Next(); ok; rr, ok = zoneParser.Next() {
		header := rr.Header()
		// Not really a good way to do this outside of this loop that I could find.
		if tld == "" {
			tld = strings.TrimSuffix(header.Name, ".")
			zone.Name = sl.String(tld)

		}
		record := datatypes.Dns_Domain_ResourceRecord{Domain: &zone}
		// This is the full record name. `www.example.com`
		fqdn := strings.TrimSuffix(header.Name, ".")
		// Just the domain part. `www`
		domain := strings.TrimSuffix(fqdn, tld)
		if domain == "" {
			record.Host = sl.String("@")
		} else {
			record.Host = sl.String(strings.TrimSuffix(domain, "."))
		}

		record.Type = sl.String(dns.TypeToString[header.Rrtype])
		record.Ttl = sl.Int(int(header.Ttl))

		switch rr.(type) {
		case *dns.NS:
			if *record.Host == "@" {
				continue
			}
			record.Data = sl.String(rr.(*dns.NS).Ns)
		case *dns.CNAME:
			record.Data = sl.String(rr.(*dns.CNAME).Target)
		case *dns.MX:
			priority := int(rr.(*dns.MX).Preference)
			record.Data = sl.String(rr.(*dns.MX).Mx)
			record.MxPriority = sl.Int(priority)
		case *dns.TXT:
			data := strings.Join(rr.(*dns.TXT).Txt, " ")
			record.Data = sl.String(strings.Trim(data, "\""))
		case *dns.SOA:
			continue
		default:
			if dns.NumField(rr) < 1 {
				record.Data = sl.String(rr.String())
			}
			record.Data = sl.String(dns.Field(rr, 1))

		}
		records = append(records, record)

	}

	// Detect any errors in parsing the Zone file
	if err := zoneParser.Err(); err != nil {
		return datatypes.Dns_Domain{}, nil, err
	}
	// This zone file was empty, but the dns library doesn't consider that an error.
	if tld == "" && len(records) == 0 {
		return datatypes.Dns_Domain{}, nil, errors.New(T("Unable to parse zone from BIND file."))
	}
	zone.ResourceRecords = records
	return zone, records, nil
}
