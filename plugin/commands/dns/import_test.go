package dns_test

import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)
var defaultCreateZoneReturn = datatypes.Dns_Domain{
	Id: sl.Int(12345),
	Name: sl.String("default.com"),
}

var defaultResourceRecordReturn = datatypes.Dns_Domain_ResourceRecord{
	Id: sl.Int(99999),
	DomainId: sl.Int(12345),
	Host: sl.String("www.default.com"),
	Type: sl.String("A"),
	Data: sl.String("192.168.1.1"),
	Ttl: sl.Int(100),
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

		Context("DNS send a file import", func() {
			BeforeEach(func() {
				fakeDNSManager.CreateZoneReturns(defaultCreateZoneReturn, nil)
				fakeDNSManager.CreateResourceRecordReturns(defaultResourceRecordReturn, nil)
				fakeDNSManager.ResourceRecordCreateReturns(defaultResourceRecordReturn, nil)
			})
			It("send a empty file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/empty_file.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse file."))
				Expect(err.Error()).To(ContainSubstring("Unable to parse zone from BIND file."))
			})

			
			It("no send a TTL in the file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/no_ttl.bind")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse file."))
				Expect(err.Error()).To(ContainSubstring("dns: not a TTL: \"$ORIGIN\" at line: 1:7"))
			})

			It("send a good file with --dry-run argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_2.bind", "--dry-run")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})

			It("send a good file", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import_3.bind")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Zone default.com was created."))
			})

			It("Complex Bind File", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "../../testfixtures/Dns_Import_Tests/dns_import.bind")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Zone default.com was created."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created resource record under zone example.com: ID=99999,"))
				for i := 0; i < 8; i++ {
					apiCall := fakeDNSManager.ResourceRecordCreateArgsForCall(i)
					Expect(apiCall.DomainId).To(Equal(complexZoneArgs[i].DomainId))
					Expect(apiCall.Type).To(Equal(complexZoneArgs[i].Type))
					Expect(apiCall.Host).To(Equal(complexZoneArgs[i].Host))
					Expect(apiCall.Ttl).To(Equal(complexZoneArgs[i].Ttl))
					Expect(apiCall.Data).To(Equal(complexZoneArgs[i].Data))
				}
			})
		})
	})
})
