package managers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("DNSManager", func() {
	var (
		fakeSLSession *session.Session
		dnsManager    managers.DNSManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		dnsManager = managers.NewDNSManager(fakeSLSession)
	})

	Describe("Create a zone", func() {
		Context("Create a zone from a name", func() {
			It("It returns a domain", func() {
				zone, err := dnsManager.CreateZone("wilma.org")
				Expect(*zone.Name).To(Equal("wilma.org"))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("List zones", func() {
		Context("List all the zones under current account", func() {
			It("It returns a list of domains", func() {
				zones, err := dnsManager.ListZones()
				Expect(err).ToNot(HaveOccurred())
				for _, zone := range zones {
					Expect(zone.Id).NotTo(Equal(nil))
					Expect(zone.Name).NotTo(Equal(nil))
					Expect(zone.Serial).NotTo(Equal(nil))
					Expect(zone.UpdateDate).NotTo(Equal(nil))
				}
			})
		})
	})

	Describe("Get ZoneId From Name", func() {
		Context("Get only one ZoneId From Name", func() {
			It("It return zone id", func() {
				zoneId, err := dnsManager.GetZoneIdFromName("ibmcloud.ibmcsf.net")
				Expect(err).ToNot(HaveOccurred())
				Expect(zoneId).To(Equal(1745158))
			})
		})
	})

	Describe("Get a zone with its resource records", func() {
		Context("Get a zone with its resource records", func() {
			It("It return zone and its resource records", func() {
				zone, err := dnsManager.GetZone(1745153, true)
				Expect(err).ToNot(HaveOccurred())
				Expect(*zone.Id).To(Equal(1745153))
				for _, record := range zone.ResourceRecords {
					Expect(*record.Id).NotTo(Equal(nil))
					Expect(*record.DomainId).To(Equal(1745153))
				}
			})
		})
	})

	Describe("Delete a domain", func() {
		Context("Delete a domain given its ID", func() {
			It("It returns nil", func() {
				err := dnsManager.DeleteZone(1745153)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Create a resource record", func() {
		Context("Create a A record under a domain", func() {
			It("It returns the created resource record", func() {
				record, err := dnsManager.CreateResourceRecord(1745153, "zookeeper1", "a", "10.106.76.197", 900)
				Expect(err).ToNot(HaveOccurred())
				Expect(*record.DomainId).To(Equal(1745153))
				Expect(*record.Host).To(Equal("zookeeper1"))
				Expect(*record.Data).To(Equal("10.106.76.197"))
				Expect(*record.Type).To(Equal("a"))
				Expect(*record.Ttl).To(Equal(900))
			})
		})
	})

	Describe("ResourceRecordCreate", func() {
		Context("Create a A record under a domain", func() {
			It("Happy path", func() {
				newRecord := datatypes.Dns_Domain_ResourceRecord{
					DomainId: sl.Int(1745153),
					Host: sl.String("TESTHOST"),
					Type: sl.String("A"),
					Data: sl.String("1.1.2.2"),
					Ttl: sl.Int(300),
				}
				record, err := dnsManager.ResourceRecordCreate(newRecord)
				Expect(err).ToNot(HaveOccurred())
				Expect(*record.DomainId).To(Equal(1745153))
				// Get the API logs from the test handler
				fakeHandler := testhelpers.GetSessionHandler(fakeSLSession)
				ApiLogs := fakeHandler.ApiCallLogs
				// Make sure we have 1 API call
				Expect(len(ApiLogs)).To(Equal(1))
				Expect(ApiLogs[0].Service).To(Equal("SoftLayer_Dns_Domain_ResourceRecord"))
				Expect(ApiLogs[0].Method).To(Equal("createObject"))
				Expect(len(ApiLogs[0].Args)).To(Equal(1))
				// Convert the arg[0] back to a Dns_Domain_ResourceRecord
				callDomain := ApiLogs[0].Args[0].(*datatypes.Dns_Domain_ResourceRecord)
				Expect(*callDomain.Host).To(Equal("TESTHOST"))
			})
		})
	})

	Describe("Delete a resource record", func() {
		Context("Delete a resource record given its ID", func() {
			It("It returns nil", func() {
				err := dnsManager.DeleteResourceRecord(50585394)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Get a resource record", func() {
		Context("Get a resource record  given its ID", func() {
			It("It return resource record", func() {
				record, err := dnsManager.GetResourceRecord(50585394)
				Expect(err).ToNot(HaveOccurred())
				Expect(*record.Id).To(Equal(50585394))
				Expect(*record.Host).To(Equal("zookeeper1"))
				Expect(*record.Data).To(Equal("10.106.76.197"))
				Expect(*record.Type).To(Equal("a"))
				Expect(*record.Ttl).To(Equal(900))
			})
		})
	})

	Describe("List resource records under a domain", func() {
		Context("List resource records under a domain given domain ID", func() {
			It("It return a list of resource records", func() {
				records, err := dnsManager.ListResourceRecords(1745153, "", "", "", 0, "")
				Expect(err).ToNot(HaveOccurred())
				for _, record := range records {
					Expect(*record.DomainId).To(Equal(1745153))
					Expect(*record.Host).NotTo(Equal(nil))
					Expect(*record.Data).NotTo(Equal(nil))
					Expect(*record.Type).NotTo(Equal(nil))
					Expect(*record.Ttl).NotTo(Equal(nil))
				}
			})
		})
	})

	Describe("Edit a resource record", func() {
		Context("Edit a resource record", func() {
			It("It return nil", func() {
				record := datatypes.Dns_Domain_ResourceRecord{
					Id:   sl.Int(1745153),
					Data: sl.String("10.106.76.197"),
					Host: sl.String("zookeeper1"),
					Ttl:  sl.Int(901),
				}
				err := dnsManager.EditResourceRecord(record)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Output a domain content", func() {
		Context("Output a domain content", func() {
			It("It return content of the domain", func() {
				content, err := dnsManager.DumpZone(1745153)
				Expect(err).ToNot(HaveOccurred())
				expeccted := "$ORIGIN bcr01.dal06.ibmcloud.ibmcsf.net.\n$TTL 86400\n@ IN SOA ns1.softlayer.com. support.softlayer.com. (\n                       2014111108        ; Serial\n                       7200              ; Refresh\n                       600               ; Retry\n                       1728000           ; Expire\n                       43200)            ; Minimum\n\n@                      86400    IN NS    ns1.softlayer.com.\n@                      86400    IN NS    ns2.softlayer.com.\n\n@                      86400    IN MX 10 mail.bcr01.dal06.ibmcloud.ibmcsf.net.\n\ntxt                    900      IN TXT   eureka-01.bcr01.dal06.ibmcloud.ibmcsf.net eureka-02.bcr01.dal06.ibmcloud.ibmcsf.net\n@                      86400    IN A     127.0.0.1\neureka-01              900      IN A     10.106.76.194\neureka-02              900      IN A     10.106.76.195\nftp                    86400    IN A     127.0.0.1\nmail                   86400    IN A     127.0.0.1\nwebmail                86400    IN A     127.0.0.1\nwww                    86400    IN A     127.0.0.1\nzookeeper1             900      IN A     10.106.76.197\nzookeeper2             900      IN A     10.106.76.198\nzookeeper3             900      IN A     10.106.76.199\n"
				Expect(content).To(Equal(expeccted))
			})
		})
	})

})
