package dns

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/miekg/dns"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
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
	bytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to read file: {{.FilePath}}.\n", map[string]interface{}{"FilePath": args[0]}), err.Error(), 2)
	}
	zone, records, err := parseFileContent(string(bytes))
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
		rr, err := cmd.DNSManager.CreateResourceRecord(*dnsDomain.Id, *record.Host, *record.Type, *record.Data, *record.Ttl)
		if err != nil {
			newError := errors.New(T("Failed to create resource record under zone {{.Zone}}: type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.\n{{.ErrorMessage}}",
				map[string]interface{}{"Zone": zone, "RecordType": *rr.Type, "Host": *rr.Host, "Data": *rr.Data, "Ttl": *rr.Ttl, "ErrorMessage": err.Error()}))
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Print(T("Created resource record under zone {{.Zone}}: ID={{.ID}}, type={{.RecordType}}, record={{.Host}}, data={{.Data}}, ttl={{.Ttl}}.",
				map[string]interface{}{"Zone": zone, "ID": *rr.Id, "RecordType": *rr.Type, "Host": *rr.Host, "Data": *rr.Data, "Ttl": *rr.Ttl}))
		}
	}

	if len(multiErrors) > 0 {
		return errors.CollapseErrors(multiErrors)
	}
	return nil
}

func parseFileContent(content string) (datatypes.Dns_Domain, []datatypes.Dns_Domain_ResourceRecord, error) {
	zone := datatypes.Dns_Domain{}
	records := []datatypes.Dns_Domain_ResourceRecord{}
	lines := strings.Split(content, "\n")
	if strings.HasPrefix(lines[0], "$ORIGIN") {
		zone.Name = sl.String(strings.TrimRight(strings.Replace(lines[0], "$ORIGIN", "", -1), "."))
	} else {
		return datatypes.Dns_Domain{}, nil, errors.New(T("Unable to parse zone from BIND file."))
	}
	for x := range dns.ParseZone(strings.NewReader(content), "", "") {
		if x.Error != nil {
			return datatypes.Dns_Domain{}, nil, x.Error
		}
		record := datatypes.Dns_Domain_ResourceRecord{Domain: &zone}
		//fmt.Println("RR", x.RR.String())
		arrs := strings.Split(x.RR.String(), "\t")
		fqdn := strings.TrimSuffix(strings.TrimSpace(arrs[0]), "\t")
		domain := strings.TrimSpace(*zone.Name)
		index := strings.Index(fqdn, domain)
		host := ""
		if index > 0 {
			host = strings.TrimRight(fqdn[0:index], ".")
		}
		if host == "" {
			record.Host = sl.String("@")
		} else {
			record.Host = sl.String(host)
		}
		record.Type = sl.String(dns.Type(x.RR.Header().Rrtype).String())
		record.Ttl = sl.Int(int(x.RR.Header().Ttl))
		if *record.Type == "SOA" {
			continue //skip SOA record in records
		} else if *record.Type == "NS" {
			record.Data = sl.String(arrs[4])
			if *record.Host == "@" {
				continue // skip the 2 default ns records
			}
		} else if *record.Type == "MX" {
			ss := strings.Split(arrs[4], " ")
			priority := 10
			priority, _ = strconv.Atoi(ss[0])
			record.Priority = sl.Int(priority)
			record.Data = sl.String(ss[1])
		} else if *record.Type == "TXT" {
			record.Data = sl.String(strings.Trim(arrs[4], "\""))
		} else {
			record.Data = sl.String(arrs[4])
		}
		//fmt.Println("zone:", *record.Domain.Name, "\thost:", *record.Host, "\ttype:", *record.Type, "\tttl:", *record.Ttl, "\tdata:", *record.Data)
		//fmt.Println()
		records = append(records, record)
	}
	zone.ResourceRecords = records
	return zone, records, nil
}
