package dns

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/miekg/dns"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ImportCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewImportCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *ImportCommand) {
	return &ImportCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *ImportCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	bytes, err := ioutil.ReadFile(c.Args()[0])
	if err != nil {
		return cli.NewExitError(T("Failed to read file: {{.FilePath}}.\n", map[string]interface{}{"FilePath": c.Args()[0]})+err.Error(), 2)
	}
	zone, records, err := parseFileContent(string(bytes))
	if err != nil {
		return cli.NewExitError(T("Failed to parse file.\n")+err.Error(), 2)
	}

	if c.IsSet("dry-run") {
		cmd.UI.Ok()
		return nil
	}

	dnsDomain, err := cmd.DNSManager.CreateZone(*zone.Name) //parseFileContent can guarantee zone.name is not nil
	if err != nil {
		return cli.NewExitError(T("Failed to create zone: {{.ZoneName}}.\n", map[string]interface{}{"ZoneName": *zone.Name})+err.Error(), 2)
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
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
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
