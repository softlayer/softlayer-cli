package managers

import (
	"errors"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	RECORD_DEFAULT_MASK = "id,host,type,ttl,data"
	RECORD_DETAIL_MASK  = "id,expire,domainId,host,minimum,refresh,retry,mxPriority,ttl,type,data,responsiblePerson"
)

//Manage SoftLayer DNS.
//See product information here: http://www.softlayer.com/DOMAIN-SERVICES
type DNSManager interface {
	GetZoneIdFromName(zoneName string) (int, error)
	ListZones() ([]datatypes.Dns_Domain, error)
	GetZone(zoneId int, getRecords bool) (datatypes.Dns_Domain, error)
	CreateZone(zoneName string) (datatypes.Dns_Domain, error)
	DeleteZone(zoneId int) error
	DumpZone(zoneId int) (string, error)
	CreateResourceRecord(zoneId int, host string, recordType string, data string, ttl int) (datatypes.Dns_Domain_ResourceRecord, error)
	DeleteResourceRecord(recordId int) error
	GetResourceRecord(recordId int) (datatypes.Dns_Domain_ResourceRecord, error)
	ListResourceRecords(zoneId int, recordType string, host string, data string, ttl int, mask string) ([]datatypes.Dns_Domain_ResourceRecord, error)
	EditResourceRecord(record datatypes.Dns_Domain_ResourceRecord) error
	SyncARecord(vs datatypes.Virtual_Guest, zoneId int, ttl int) error
	SyncAAAARecord(vs datatypes.Virtual_Guest, zoneId int, ttl int) error
	SyncPTRRecord(vs datatypes.Virtual_Guest, ttl int) error
}

type DNSmanager struct {
	DomainService         services.Dns_Domain
	ResourceRecordService services.Dns_Domain_ResourceRecord
	AccountService        services.Account
	VirtualServerService  services.Virtual_Guest
}

func NewDNSManager(session *session.Session) *DNSmanager {
	return &DNSmanager{
		services.GetDnsDomainService(session),
		services.GetDnsDomainResourceRecordService(session),
		services.GetAccountService(session),
		services.GetVirtualGuestService(session),
	}
}

//Return zoen ID based on a zone name.
func (dns DNSmanager) GetZoneIdFromName(zoneName string) (int, error) {
	domains, err := dns.AccountService.Mask("id").Filter(filter.New(filter.Path("domains.name").Eq(zoneName)).Build()).GetDomains()
	if err != nil {
		return 0, err
	}
	if len(domains) != 1 {
		return 0, errors.New(T("Failed to find ID for domain: {{.Domain}} on your account.", map[string]interface{}{"Domain": zoneName}))
	}
	return *domains[0].Id, nil
}

//Retrieve a list of all DNS zones in current account.
func (dns DNSmanager) ListZones() ([]datatypes.Dns_Domain, error) {
	return dns.AccountService.GetDomains()
}

//Get a zone and its resource records.
//zone: the zone name
//getRecords: whether to return resource records
func (dns DNSmanager) GetZone(zoneId int, getRecords bool) (datatypes.Dns_Domain, error) {
	if getRecords == true {
		return dns.DomainService.Mask("resourceRecords").GetObject()
	}
	return dns.DomainService.Id(zoneId).GetObject()
}

//Create a zone from the specified parameters, with default SOA and NS records
//zone: the zone name to create
func (dns DNSmanager) CreateZone(zoneName string) (datatypes.Dns_Domain, error) {
	dnsDomain := datatypes.Dns_Domain{}
	dnsDomain.Name = sl.String(zoneName)
	dnsDomain.ResourceRecords = []datatypes.Dns_Domain_ResourceRecord{
		datatypes.Dns_Domain_ResourceRecord{
			Type: sl.String("SOA"),
			Data: sl.String("ns1.softlayer.com. support.softlayer.com."),
			Ttl:  sl.Int(86400),
		},
		datatypes.Dns_Domain_ResourceRecord{
			Type: sl.String("NS"),
			Data: sl.String("ns1.softlayer.com"),
			Ttl:  sl.Int(86400),
		},
		datatypes.Dns_Domain_ResourceRecord{
			Type: sl.String("NS"),
			Data: sl.String("ns2.softlayer.com"),
			Ttl:  sl.Int(86400),
		},
	}
	return dns.DomainService.CreateObject(&dnsDomain)
}

//Delete a zone by its ID
//zoneId: the zone ID to delete
func (dns DNSmanager) DeleteZone(zoneId int) error {
	_, err := dns.DomainService.Id(zoneId).DeleteObject()
	return err
}

//Create a resource record on a domain
//zoneId: the zone's ID
//host: the name of the record to add
//recordType: the type of record (A, AAAA, CNAME, MX, TXT, etc.)
//data: the record's value
//ttl: the TTL or time-to-live value (default: 60)
func (dns DNSmanager) CreateResourceRecord(zoneId int, host string, recordType string, data string, ttl int) (datatypes.Dns_Domain_ResourceRecord, error) {
	if ttl == 0 {
		ttl = 60
	}
	record := datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(zoneId),
		Host:     sl.String(host),
		Type:     sl.String(recordType),
		Data:     sl.String(data),
		Ttl:      sl.Int(ttl),
	}
	return dns.ResourceRecordService.CreateObject(&record)
}

//Delete a resource record by its ID
//recordId: the record's ID
func (dns DNSmanager) DeleteResourceRecord(recordId int) error {
	_, err := dns.ResourceRecordService.Id(recordId).DeleteObject()
	return err
}

//Get a DNS resource record
//recordId: the record's ID
func (dns DNSmanager) GetResourceRecord(recordId int) (datatypes.Dns_Domain_ResourceRecord, error) {
	return dns.ResourceRecordService.Id(recordId).GetObject()
}

//List resource records within a zone
//zoneId: the zone name in which to search.
//recordType: the type of record
//host: record's host
//data: the records data
//ttl: time in seconds
//mask: mask of properties
func (dns DNSmanager) ListResourceRecords(zoneId int, recordType string, host string, data string, ttl int, mask string) ([]datatypes.Dns_Domain_ResourceRecord, error) {
	filters := filter.New()
	if data != "" {
		filters = append(filters, utils.QueryFilter(data, "resourceRecords.data"))
	}
	if host != "" {
		filters = append(filters, utils.QueryFilter(host, "resourceRecords.host"))
	}
	if ttl != 0 {
		filters = append(filters, filter.Path("resourceRecords.ttl").Eq(ttl))
	}
	if recordType != "" {
		filters = append(filters, filter.Path("resourceRecords.type").Eq(recordType))
	}
	if mask == "" {
		mask = RECORD_DEFAULT_MASK
	}
	return dns.DomainService.Id(zoneId).Filter(filters.Build()).GetResourceRecords()
}

//Update an existing record with specified options
//record: the record to update, it must include id
func (dns DNSmanager) EditResourceRecord(record datatypes.Dns_Domain_ResourceRecord) error {
	_, err := dns.ResourceRecordService.Id(*record.Id).EditObject(&record)
	return err
}

//Retrieve a zone dump in BIND format
//zoneId: The zone ID to dump
func (dns DNSmanager) DumpZone(zoneId int) (string, error) {
	return dns.DomainService.Id(zoneId).GetZoneFileContents()
}

//Sync A record.
//vs: virtual server the A record is belong to
//zoneId: The zone ID to be added/updated
//ttl: Time-To-Live to be updated
func (dns DNSmanager) SyncARecord(vs datatypes.Virtual_Guest, zoneId int, ttl int) error {
	if vs.PrimaryIpAddress == nil {
		return errors.New(T("No primary IP address associated with virtual server instance: {{.VsId}}.", map[string]interface{}{"VsId": *vs.Id}))
	}
	var hostname, ipAddress string
	if vs.Hostname != nil {
		hostname = *vs.Hostname
	}
	if vs.PrimaryIpAddress != nil {
		ipAddress = *vs.PrimaryIpAddress
	}
	records, err := dns.ListResourceRecords(zoneId, "a", hostname, ipAddress, 0, "")
	if err != nil {
		return err
	}
	if len(records) == 0 {
		_, err := dns.CreateResourceRecord(zoneId, hostname, "a", ipAddress, ttl)
		if err != nil {
			return err
		}
	} else if len(records) != 1 {
		return errors.New(T("Aborting A record sync, found {{.Num}} A records exists!", map[string]interface{}{"Num": len(records)}))
	} else {
		record := records[0]
		record.Host = &hostname
		record.Data = &ipAddress
		record.Ttl = &ttl
		err := dns.EditResourceRecord(record)
		if err != nil {
			return err
		}
	}
	return nil
}

//Sync AAAA record.
//vs: virtual server the A record is belong to
//zoneId: The zone ID to be added/updated
//ttl: Time-To-Live to be updated
func (dns DNSmanager) SyncAAAARecord(vs datatypes.Virtual_Guest, zoneId int, ttl int) error {
	if vs.PrimaryNetworkComponent == nil ||
		vs.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord == nil ||
		vs.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress == nil {
		return errors.New(T("No IP V6 address associated with virtual server instance: {{.VsId}}.", map[string]interface{}{"VsId": *vs.Id}))
	}
	var hostname, ipV6Address string
	if vs.Hostname != nil {
		hostname = *vs.Hostname
	}
	if vs.PrimaryNetworkComponent != nil && vs.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil && vs.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress != nil {
		ipV6Address = *vs.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress
	}
	records, err := dns.ListResourceRecords(zoneId, "aaaa", hostname, ipV6Address, 0, "")
	if err != nil {
		return err
	}
	if len(records) == 0 {
		_, err := dns.CreateResourceRecord(zoneId, hostname, "aaaa", ipV6Address, ttl)
		if err != nil {
			return err
		}
	} else if len(records) != 1 {
		return errors.New(T("Aborting AAAA record sync, found {{.Num}} AAAA records exists!", map[string]interface{}{"Num": len(records)}))
	} else {
		record := records[0]
		record.Host = &hostname
		record.Data = &ipV6Address
		record.Ttl = &ttl
		err := dns.EditResourceRecord(record)
		if err != nil {
			return err
		}
	}
	return nil
}

//Sync PTR record.
//vs: virtual server the PTR record is belong to
//ttl: Time-To-Live to be updated
func (dns DNSmanager) SyncPTRRecord(vs datatypes.Virtual_Guest, ttl int) error {
	var hostRecords []string
	if vs.PrimaryIpAddress != nil {
		hostRecords = strings.Split(*vs.PrimaryIpAddress, ".")
	}
	hostRec := hostRecords[len(hostRecords)-1]
	ptrDomains, err := dns.VirtualServerService.Id(*vs.Id).GetReverseDomainRecords()
	if err != nil {
		return err
	}
	if len(ptrDomains) < 1 {
		return errors.New(T("Domain record not found"))
	}
	editPtr := datatypes.Dns_Domain_ResourceRecord{}
	for _, ptr := range ptrDomains[0].ResourceRecords {
		if ptr.Host != nil && *ptr.Host == hostRec {
			ptr.Ttl = &ttl
			editPtr = ptr
			break
		}
	}
	if editPtr.Id == nil {
		if vs.FullyQualifiedDomainName == nil {
			return errors.New(T("Virtual guest domain name not found."))
		}
		_, err := dns.CreateResourceRecord(*ptrDomains[0].Id, hostRec, "ptr", *vs.FullyQualifiedDomainName, ttl)
		if err != nil {
			return err
		}
	} else {
		editPtr.Data = vs.FullyQualifiedDomainName
		editPtr.IsGatewayAddress = nil
		err := dns.EditResourceRecord(editPtr)
		if err != nil {
			return err
		}
	}
	return nil
}
