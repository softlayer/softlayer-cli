package dns_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	mdns "github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var defaultCreateZoneReturn = datatypes.Dns_Domain{
	Id:   sl.Int(12345),
	Name: sl.String("default.com"),
}

var defaultResourceRecordReturn = datatypes.Dns_Domain_ResourceRecord{
	Id:       sl.Int(99999),
	DomainId: sl.Int(12345),
	Host:     sl.String("www.default.com"),
	Type:     sl.String("A"),
	Data:     sl.String("192.168.1.1"),
	Ttl:      sl.Int(100),
}

// See `plugin/testfixtures/dns_import.bind` for where these come from
var complexZoneArgs = []datatypes.Dns_Domain_ResourceRecord{
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("@"),
		Type: sl.String("MX"), Ttl: sl.Int(1814400),
		Data: sl.String("mail.example.com."),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("@"),
		Type: sl.String("MX"), Ttl: sl.Int(1814400),
		Data: sl.String("mail.example.net."),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("ns1"),
		Type: sl.String("A"), Ttl: sl.Int(1814400),
		Data: sl.String("192.168.254.2"),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("mail"),
		Type: sl.String("A"), Ttl: sl.Int(1814400),
		Data: sl.String("192.168.254.4"),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("joe"),
		Type: sl.String("A"), Ttl: sl.Int(1814400),
		Data: sl.String("192.168.254.6"),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("www"),
		Type: sl.String("A"), Ttl: sl.Int(1814400),
		Data: sl.String("192.168.254.7"),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("ftp"),
		Type: sl.String("CNAME"), Ttl: sl.Int(1814400),
		Data: sl.String("ftp.example.net."),
	},
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("@"),
		Type: sl.String("TXT"), Ttl: sl.Int(3600),
		Data: sl.String("v=spf1 includespf.dynect.net ~all"),
	},
}

var srvZoneArgs = []datatypes.Dns_Domain_ResourceRecord{
	datatypes.Dns_Domain_ResourceRecord{
		DomainId: sl.Int(12345), Host: sl.String("_serviceTest._tls.host.local"),
		Type: sl.String("SRV"), Ttl: sl.Int(900),
		Data: sl.String("target.field.btest1.com."),
	},
}

var _ = Describe("DNS Import", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.ImportCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewImportCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("DNS import", func() {
		Context("DNS import without file", func() {
			It("without any argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("with an inexist file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "not-exist.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read file: not-exist.txt."))
			})
		})
		Context("Test GetZone Function", func() {
			Context("Test Happy Path", func() {
				BeforeEach(func() {
					fakeDNSManager.CreateZoneReturns(defaultCreateZoneReturn, nil)
				})
				It("Success", func() {
					zone, err := dns.CreateOrGetZone("test.com", fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(*zone.Id).To(Equal(12345))
				})
			})
			Context("Test Error Handling", func() {
				It("Zone Exists", func() {
					fakeDNSManager.CreateZoneReturns(datatypes.Dns_Domain{}, errors.New("Zone Exists"))
					fakeDNSManager.GetZoneIdFromNameReturns(99999, nil)
					zone, err := dns.CreateOrGetZone("test.com", fakeDNSManager)
					apiCall := fakeDNSManager.GetZoneIdFromNameArgsForCall(0)
					Expect(apiCall).To(Equal("test.com"))
					Expect(err).NotTo(HaveOccurred())
					Expect(*zone.Id).To(Equal(99999))
				})
				It("Everything went wrong", func() {
					fakeDNSManager.CreateZoneReturns(datatypes.Dns_Domain{}, errors.New("Zone Exists"))
					fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("No Zone"))
					_, err := dns.CreateOrGetZone("test.com", fakeDNSManager)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Zone Exists"))
				})
			})
		})
		Context("Test CreateRecord", func() {
			Context("Test NS record", func() {
				It("No Host", func() {
					rr, _ := mdns.NewRR("@ 3600 IN NS testNs.com.")
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(false))
					apiCalls := fakeDNSManager.ResourceRecordCreateCallCount()
					Expect(apiCalls).To(Equal(0))
				})
				It("Create NS record", func() {
					rr, _ := mdns.NewRR("@ 3600 IN NS testNs.com")
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("testHost")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("testNs.com."))
					Expect(*apiCall.Host).To(Equal("testHost"))
				})
			})
			Context("Test CNAME", func() {
				It("Create CNAME", func() {
					rr, _ := mdns.NewRR("someHost 3600 IN CNAME testcname.")
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("testcname."))
				})
			})
			Context("Test MX", func() {
				It("Create MX", func() {
					rr, _ := mdns.NewRR("mail 3600 IN MX 99 mail.example.net.")
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("mail")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("mail.example.net."))
					Expect(*apiCall.MxPriority).To(Equal(99))
				})
			})
			Context("Test TXT", func() {
				It("Create TXT", func() {
					rr, _ := mdns.NewRR(`@        3600 TXT   "v=spf1 includespf.dynect.net ~all"`)
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("v=spf1 includespf.dynect.net ~all"))
				})
			})
			Context("Test SRV", func() {
				It("Create SRV", func() {
					rr, _ := mdns.NewRR(`_serviceTest._tls.host.local 900      IN SRV   15 20 5005 target.field.`)
					record := datatypes.Dns_Domain_ResourceRecord{
						Host: sl.String("@"), DomainId: sl.Int(12344), Type: sl.String("SRV"), Ttl: sl.Int(900),
					}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					createCalls := fakeDNSManager.ResourceRecordCreateCallCount()
					Expect(createCalls).To(Equal(0))
					apiCall := fakeDNSManager.SrvResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.DomainId).To(Equal(12344))
					Expect(*apiCall.Type).To(Equal("SRV"))
					Expect(*apiCall.Ttl).To(Equal(900))
					Expect(*apiCall.Host).To(Equal("@"))
					Expect(*apiCall.Priority).To(Equal(15))
					Expect(*apiCall.Port).To(Equal(5005))
					Expect(*apiCall.Weight).To(Equal(20))
					Expect(*apiCall.Data).To(Equal("target.field."))
					Expect(*apiCall.Service).To(Equal("_serviceTest"))
					Expect(*apiCall.Protocol).To(Equal("_tls"))
				})
				It("Some error during create", func() {
					rr, _ := mdns.NewRR(`_serviceTest._tls.host.local 900      IN SRV   15 20 5005 target.field.`)
					record := datatypes.Dns_Domain_ResourceRecord{
						Host: sl.String("@"), DomainId: sl.Int(12344), Type: sl.String("SRV"), Ttl: sl.Int(900),
					}
					fakeDNSManager.SrvResourceRecordCreateReturns(datatypes.Dns_Domain_ResourceRecord_SrvType{}, errors.New("fake error"))
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).To(HaveOccurred())
					Expect(response).To(Equal(false))
					Expect(err.Error()).To(ContainSubstring("fake error"))
				})
				It("Bad SRV record", func() {
					rr, _ := mdns.NewRR(`local 900      IN SRV   15 20 5005 target.field.`)
					record := datatypes.Dns_Domain_ResourceRecord{
						Host: sl.String("@"), DomainId: sl.Int(12344), Type: sl.String("SRV"), Ttl: sl.Int(900),
					}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).To(HaveOccurred())
					Expect(response).To(Equal(false))
					Expect(err.Error()).To(ContainSubstring("Invalid SRV record:"))
				})
			})
			Context("Test SOA", func() {
				It("Create SOA", func() {
					rr, _ := mdns.NewRR(`@ IN SOA ns1.softlayer.com. support.softlayer.com. (
                       2020042204        ; Serial
                       7200              ; Refresh
                       600               ; Retry
                       1728000           ; Expire
                       43200)            ; Minimum`)
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(false))
					createCalls := fakeDNSManager.ResourceRecordCreateCallCount()
					Expect(createCalls).To(Equal(0))
				})
			})
			Context("Test Default", func() {
				It("Create AAAA", func() {
					rr, _ := mdns.NewRR(`www        3600 AAAA  1100:2202:33c4:4410:55d3:6635:77a4:8ecc`)
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("1100:2202:33c4:4410:55d3:6635:77a4:8ecc"))
				})
				It("Create A", func() {
					rr, _ := mdns.NewRR(`www        3600 A  192.68.1.1`)
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(Equal(true))
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(0)
					Expect(*apiCall.Data).To(Equal("192.68.1.1"))
				})
			})
			Context("Test Error Handling", func() {
				It("Create A", func() {
					rr, _ := mdns.NewRR(`www        3600 A  192.68.1.1`)
					record := datatypes.Dns_Domain_ResourceRecord{Host: sl.String("@")}
					fakeDNSManager.ResourceRecordCreateReturns(datatypes.Dns_Domain_ResourceRecord{}, errors.New("fake error"))

					response, err := dns.CreateRecord(&record, rr, fakeDNSManager)
					Expect(err).To(HaveOccurred())
					Expect(response).To(Equal(false))
					Expect(err.Error()).To(ContainSubstring("fake error"))
				})
			})
		})
		Context("Error Handling", func() {
			It("Fails to get or create zone", func() {
				fakeDNSManager.CreateZoneReturns(defaultCreateZoneReturn, errors.New("CreateZoneFail"))
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("GetZoneIdFromNameFail"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_2.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create zone: dal06.bluemix.ibmcsf.net."))
				Expect(err.Error()).To(ContainSubstring("CreateZoneFail"))
			})
			It("CreateRecord fails", func() {
				fakeDNSManager.CreateZoneReturns(defaultCreateZoneReturn, nil)
				fakeDNSManager.ResourceRecordCreateReturns(datatypes.Dns_Domain_ResourceRecord{}, errors.New("ResourceRecordCreateReturnsFail"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_2.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ResourceRecordCreateReturnsFail"))
			})
		})
		Context("DNS send a file import", func() {
			BeforeEach(func() {
				fakeDNSManager.CreateZoneReturns(defaultCreateZoneReturn, nil)
				fakeDNSManager.CreateResourceRecordReturns(defaultResourceRecordReturn, nil)
				fakeDNSManager.ResourceRecordCreateReturns(defaultResourceRecordReturn, nil)
			})
			It("send a empty file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/empty_file.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to parse zone from BIND file."))
			})

			It("no send a TTL in the file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/no_ttl.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("dns: not a TTL: \"$ORIGIN\" at line: 1:7"))
			})

			It("send a good file with --dry-run argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_2.bind", "--dry-run")
				Expect(err).NotTo(HaveOccurred())
				createCalls := fakeDNSManager.ResourceRecordCreateCallCount()
				Expect(createCalls).To(Equal(0))
			})

			It("send a good file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_3.bind")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Domain: dal06.bluemix.ibmcsf.net Id: 12345"))
			})

			It("Complex Bind File", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import.bind")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Domain: example.com Id: 12345"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created Record: @ 1814400 MX mail.example.com."))
				for i := 0; i < 8; i++ {
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(i)
					Expect(apiCall.DomainId).To(Equal(complexZoneArgs[i].DomainId))
					Expect(apiCall.Type).To(Equal(complexZoneArgs[i].Type))
					Expect(apiCall.Host).To(Equal(complexZoneArgs[i].Host))
					Expect(apiCall.Ttl).To(Equal(complexZoneArgs[i].Ttl))
					Expect(apiCall.Data).To(Equal(complexZoneArgs[i].Data))
				}
			})
			It("SRV record", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_srv.bind")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Domain: btest1.com Id: 12345"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created Record: _serviceTest._tls.host.local 900 SRV target.field.btest1.com."))

				apiCall := fakeDNSManager.SrvResourceRecordCreateArgsForCall(0)
				Expect(apiCall.DomainId).To(Equal(srvZoneArgs[0].DomainId))
				Expect(apiCall.Type).To(Equal(srvZoneArgs[0].Type))
				Expect(apiCall.Host).To(Equal(srvZoneArgs[0].Host))
				Expect(apiCall.Ttl).To(Equal(srvZoneArgs[0].Ttl))
				Expect(apiCall.Data).To(Equal(srvZoneArgs[0].Data))
			})
		})
	})
})
