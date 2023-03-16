package dns

import (
	"os"
	"strings"
	"github.com/miekg/dns"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

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
	zoneFile, err := os.Open(args[0])
	if err != nil {
	  return errors.NewAPIError(T("Failed to read file: {{.FilePath}}.\n", map[string]interface{}{"FilePath": args[0]}), err.Error(), 2)
	}
	tld := ""

	zone := datatypes.Dns_Domain{}
	zoneParser := dns.NewZoneParser(zoneFile, "", args[0])
	// Go through the zone file line by line
	for rr, ok := zoneParser.Next(); ok; rr, ok = zoneParser.Next() {
		header := rr.Header()
      // If we dont know what the TLD is yet, figure it out.
		if tld == "" {
			tld = strings.TrimSuffix(header.Name, ".")
			zone, err = CreateOrGetZone(tld, cmd.DNSManager)
			if err != nil {
				subs := map[string]interface{}{"ZoneName": tld}
				return errors.NewAPIError(T("Failed to create zone: {{.ZoneName}}.\n", subs), err.Error(), 2)
			}
			cmd.UI.Print("Domain: %s Id: %d", tld, *zone.Id)
		} // END TLD == ""

		record := datatypes.Dns_Domain_ResourceRecord{}
		// Set the DomainId to the Zone we are creating this record under
		record.DomainId = zone.Id
		// Set a default TTL incase the record doesn't have one
		record.Ttl = sl.Int(3600)
		// This is the full record name. `www.example.com`
		fqdn := strings.TrimSuffix(header.Name, ".")
		// Just the domain part. `www`, or if its blank, set to `@`
		domain := strings.TrimSuffix(fqdn, tld)
		if domain == "" {
		   record.Host = sl.String("@")
		} else {
			// The SL API assumes the hostpart will not end with .
		   record.Host = sl.String(strings.TrimSuffix(domain, "."))
		}

		record.Type = sl.String(dns.TypeToString[header.Rrtype])
		record.Ttl = sl.Int(int(header.Ttl))
		record.Data = sl.String("")
		created := false
		if cmd.DryRun {
			record.Data = sl.String(dns.Field(rr, 1))
			created = false
		} else {
			created, err = CreateRecord(&record, rr, cmd.DNSManager)
			if err != nil {
				return err
			}
		}
		if created {
			cmd.UI.Print("Created Record: %v %v %v %v", *record.Host, *record.Ttl, *record.Type, *record.Data)	
		} else {
			cmd.UI.Print("Parsed Record: %v %v %v %v", *record.Host, *record.Ttl, *record.Type, *record.Data)
		}
		
	} // END for rr, ok
	// Detect any errors in parsing the Zone file
	if err := zoneParser.Err(); err != nil {
		return err
	}
	// This zone file was empty, but the dns library doesn't consider that an error.
	if tld == "" {
		return errors.New(T("Unable to parse zone from BIND file."))
	}

	return nil
}

func CreateOrGetZone(tld string, manager managers.DNSManager) (datatypes.Dns_Domain, error) {
	zone, createZoneErr := manager.CreateZone(tld)
	// Possibly this zone exists already
	if createZoneErr != nil {
		zoneId, getZoneErr := manager.GetZoneIdFromName(tld)
		// Error creating zone, and getting a zone with this name, failure.
		if getZoneErr != nil {
			return zone, createZoneErr  
		} else {
			zone.Id = sl.Int(zoneId)
			zone.Name = sl.String(tld)
		}
	}
	return zone, nil
}

// Returns true, nil when a record was created
// false, nil when a record was not created
// false, err when an error happens
func CreateRecord(record *datatypes.Dns_Domain_ResourceRecord, rr dns.RR, manager managers.DNSManager) (bool, error) {
	switch rr.(type) {
	case *dns.NS:
	   if *record.Host == "@" {
	       return false, nil
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
	case *dns.SRV:
		srvrr := rr.(*dns.SRV)
		parts := strings.Split(srvrr.Header().Name, ".")
		if len(parts) < 3 {
			subs := map[string]interface{}{"sub": srvrr.String()}
			return false, errors.New(T("Invalid SRV record: {{.sub}}.", subs))
		}
	   srvRecord := datatypes.Dns_Domain_ResourceRecord_SrvType{}
	   priority := int(srvrr.Priority)
	   srvRecord.DomainId = record.DomainId
	   srvRecord.Type = record.Type
	   srvRecord.Ttl = record.Ttl
	   srvRecord.Host = record.Host
	   srvRecord.Priority = sl.Int(priority)
	   srvRecord.Port = sl.Int(int(srvrr.Port))
	   srvRecord.Weight = sl.Int(int(srvrr.Weight))
	   srvRecord.Data = sl.String(srvrr.Target)
	   srvRecord.Service = sl.String(parts[0])
	   srvRecord.Protocol = sl.String(parts[1])
	   // So we can add this to our output at the end
	   record.Data = sl.String(rr.(*dns.SRV).Target)
	   _, err := manager.SrvResourceRecordCreate(srvRecord)
	   if err != nil {
	   	return false, err
	   } else {
	   	return true, nil
	   }
	case *dns.SOA:
	   return false, nil
	default:
		// Sets data to the end bit of the resource record, basically.
	   record.Data = sl.String(dns.Field(rr, 1))
	}
	_, err := manager.ResourceRecordCreate(*record)
	if err != nil {
		return false, err
	}
	return true, nil
}