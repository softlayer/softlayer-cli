package dns_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Record list", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.RecordListCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewRecordListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Record list", func() {
		Context("Record list without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Record list with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get zone ID from zone name: abc.com."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record list with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list resource records under zone: abc.com"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))

			})
		})

		Context("Record list", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585314),
						Type: sl.String("txt"),
						Host: sl.String("txt"),
						Data: sl.String("bcr01.dal06.bluemix.ibmcsf.net"),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585306),
						Type: sl.String("soa"),
						Host: sl.String("@"),
						Data: sl.String("ns1.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585307),
						Type: sl.String("ns"),
						Host: sl.String("@"),
						Data: sl.String("ns1.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585308),
						Type: sl.String("ns"),
						Host: sl.String("@"),
						Data: sl.String("ns2.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585313),
						Type: sl.String("mx"),
						Host: sl.String("@"),
						Data: sl.String("mail.dal06.bluemix.ibmcsf.net."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585305),
						Type: sl.String("a"),
						Host: sl.String("@"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585312),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585309),
						Type: sl.String("a"),
						Host: sl.String("mail"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585310),
						Type: sl.String("a"),
						Host: sl.String("webmail"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585311),
						Type: sl.String("a"),
						Host: sl.String("www"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("50585314"))
				Expect(results[1]).To(ContainSubstring("txt       TXT    900     bcr01.dal06.bluemix.ibmcsf.net"))
				Expect(results[2]).To(ContainSubstring("50585306"))
				Expect(results[2]).To(ContainSubstring("@         SOA    900     ns1.softlayer.com."))
				Expect(results[3]).To(ContainSubstring("50585307"))
				Expect(results[3]).To(ContainSubstring("@         NS     900     ns1.softlayer.com."))
				Expect(results[4]).To(ContainSubstring("50585308"))
				Expect(results[4]).To(ContainSubstring("@         NS     900     ns2.softlayer.com."))
				Expect(results[5]).To(ContainSubstring("50585313"))
				Expect(results[5]).To(ContainSubstring("@         MX     900     mail.dal06.bluemix.ibmcsf.net."))
				Expect(results[6]).To(ContainSubstring("50585305"))
				Expect(results[6]).To(ContainSubstring("@         A      900     127.0.0.1"))
				Expect(results[7]).To(ContainSubstring("50585312"))
				Expect(results[7]).To(ContainSubstring("ftp       A      86400   127.0.0.1"))
				Expect(results[8]).To(ContainSubstring("50585309"))
				Expect(results[8]).To(ContainSubstring("mail      A      86400   127.0.0.1"))
				Expect(results[9]).To(ContainSubstring("50585310"))
				Expect(results[9]).To(ContainSubstring("webmail   A      86400   127.0.0.1"))
				Expect(results[10]).To(ContainSubstring("50585311"))
				Expect(results[10]).To(ContainSubstring("www       A      86400   127.0.0.1"))
			})
		})
	})
})
