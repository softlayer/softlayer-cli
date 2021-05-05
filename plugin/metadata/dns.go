package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_DNS_NAME  = "dns"
	CMD_DNS_NAME = "dns"

	CMD_DNS_IMPORT_NAME        = "import"
	CMD_DNS_RECORD_ADD_NAME    = "record-add"
	CMD_DNS_RECORD_EDIT_NAME   = "record-edit"
	CMD_DNS_RECORD_LIST_NAME   = "record-list"
	CMD_DNS_RECORD_REMOVE_NAME = "record-remove"
	CMD_DNS_ZONE_CREATE_NAME   = "zone-create"
	CMD_DNS_ZONE_DELETE_NAME   = "zone-delete"
	CMD_DNS_ZONE_LIST_NAME     = "zone-list"
	CMD_DNS_ZONE_PRINT_NAME    = "zone-print"
)

func DnsNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_DNS_NAME,
		Description: T("Classic infrastructure Domain Name System"),
	}
}

func DnsMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_DNS_NAME,
		Description: T("Classic infrastructure Domain Name System"),
		Usage:       "${COMMAND_NAME} sl dns",
		Subcommands: []cli.Command{
			DnsImportMetaData(),
			DnsRecordAddMetaData(),
			DnsRecordEditMetaData(),
			DnsRecordListMetaData(),
			DnsRecordRemoveMetaData(),
			DnsZoneCreateMetaData(),
			DnsZoneDeleteMetaData(),
			DnsZoneListMetaData(),
			DnsZonePrintMetaData(),
		},
	}
}
func DnsImportMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_IMPORT_NAME,
		Description: T("Import a zone based off a BIND zone file"),
		Usage: T(`${COMMAND_NAME} sl dns import ZONEFILE [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns import ~/ibm.com.txt
   This command imports zone and its resource records from file: ~/ibm.com.txt.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: T("Don't actually create records"),
			},
		},
	}
}

func DnsRecordAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_RECORD_ADD_NAME,
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
			OutputFlag(),
		},
	}
}

func DnsRecordEditMetaData() cli.Command {

	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_RECORD_EDIT_NAME,
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

func DnsRecordListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_RECORD_LIST_NAME,
		Description: T("List all the resource records in a zone"),
		Usage: T(`${COMMAND_NAME} sl dns record-list ZONE [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl dns record-list ibm.com --record elasticsearch --type A --ttl 900
   This command lists all A records under the zone: ibm.com, and filters by host is elasticsearch and ttl is 900 seconds.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "data",
				Usage: T("Filter by record data, such as an IP address"),
			},
			cli.StringFlag{
				Name:  "record",
				Usage: T("Filter by host record, such as www"),
			},
			cli.IntFlag{
				Name:  "ttl",
				Usage: T("Filter by TTL(Time-To-Live) in seconds, such as 86400"),
			},
			cli.StringFlag{
				Name:  "type",
				Usage: T("Filter by record type, such as A or CNAME"),
			},
			OutputFlag(),
		},
	}
}
func DnsRecordRemoveMetaData() cli.Command {

	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_RECORD_REMOVE_NAME,
		Description: T("Remove resource record from a zone"),
		Usage: T(`${COMMAND_NAME} sl dns record-remove RECORD_ID

	
EXAMPLE:
   ${COMMAND_NAME} sl dns record-remove 12345678
   This command removes resource record with ID 12345678.`),
	}
}

func DnsZoneCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_ZONE_CREATE_NAME,
		Description: T("Create a zone"),
		Usage: T(`${COMMAND_NAME} sl dns zone-create ZONE [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-create ibm.com 
   This command creates a zone that is named ibm.com.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func DnsZoneDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_ZONE_DELETE_NAME,
		Description: T("Delete a zone"),
		Usage: T(`${COMMAND_NAME} sl dns zone-delete ZONE

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-delete ibm.com 
   This command deletes a zone that is named ibm.com.`),
	}
}

func DnsZoneListMetaData() cli.Command {

	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_ZONE_LIST_NAME,
		Description: T("List all zones on your account"),
		Usage: T(`${COMMAND_NAME} sl dns zone-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-list
   This command lists all zones under current account.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func DnsZonePrintMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_DNS_NAME,
		Name:        CMD_DNS_ZONE_PRINT_NAME,
		Description: T("Print zone and resource records in BIND format"),
		Usage: T(`${COMMAND_NAME} sl dns zone-print ZONE

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-print ibm.com
   This command prints zone that is named ibm.com, and in BIND format.`),
	}
}
