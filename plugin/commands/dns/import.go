package dns

import (

	"os"
	"strings"
	"fmt"
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
			zone, err := GetZone(tld, cmd.DNSManager)
			if err != nil {
				subs := map[string]interface{}{"ZoneName": tld}
				return errors.NewAPIError(T("Failed to create zone: {{.ZoneName}}.\n", subs), err.Error(), 2)
			}
			fmt.Printf("Zone is : %v\n", *zone.Id)
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
		createErr := CreateRecord(&record, rr, cmd.DNSManager)
		if createErr != nil {
			return createErr
		}
		fmt.Printf("Record.Host: %v Record.Data: %v\n", *record.Host, *record.Data)
	} // END for rr, ok
	return nil
}

func GetZone(tld string, manager managers.DNSManager) (zone datatypes.Dns_Domain, err error) {
	zone, createZoneErr := manager.CreateZone(tld)
	// Possibly this zone exists already
	if createZoneErr != nil {
		zoneId, getZoneErr := manager.GetZoneIdFromName(tld)
		// Error creating zone, and getting a zone with this name, failure.
		if getZoneErr != nil {
			return zone, createZoneErr  
		} else {
			zone.Id = sl.Int(zoneId)
		}
	}
	return zone, nil
}

func CreateRecord(record *datatypes.Dns_Domain_ResourceRecord, rr dns.RR, manager managers.DNSManager) error {
	switch rr.(type) {
	case *dns.NS:
	   if *record.Host == "@" {
	       return nil
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
	   srvRecord := datatypes.Dns_Domain_ResourceRecord_SrvType{}
	   priority := int(rr.(*dns.SRV).Priority)
	   srvRecord.DomainId = record.DomainId
	   srvRecord.Type = record.Type
	   srvRecord.Ttl = record.Ttl
	   srvRecord.Host = record.Host
	   srvRecord.Priority = sl.Int(priority)
	   srvRecord.Port = sl.Int(int(rr.(*dns.SRV).Port))
	   srvRecord.Weight = sl.Int(int(rr.(*dns.SRV).Weight))
	   srvRecord.Data = sl.String(rr.(*dns.SRV).Target)
	   _, err := manager.SrvResourceRecordCreate(srvRecord)
	   if err != nil {
	   	return err
	   } else {
	   	return nil
	   }
	case *dns.SOA:
	   return nil
	default:
	   if dns.NumField(rr) < 1 {
	       record.Data = sl.String(rr.String())
	   }
	   record.Data = sl.String(dns.Field(rr, 1))
	}
	_, err := manager.ResourceRecordCreate(*record)
	if err != nil {
		return err
	}
	return nil
}