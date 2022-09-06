package dns_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
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

var _ = Describe("Record edit", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.RecordEditCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewRecordEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Record edit", func() {
		Context("Record edit without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Record edit with wrong zone name", func() {
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

		Context("Record edit with listing records fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list resource records under zone: abc.com."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record edit with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(errors.New("Internal Server Error"))
			})
			It("return  error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--data", "127.0.0.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update resource record 1234 under zone abc.com."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record edit with multiple records and server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(5678),
						Type: sl.String("a"),
						Host: sl.String("mail"),
						Data: sl.String("127.0.0.8"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(errors.New("Internal Server Error"))
			})
			It("return  error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--data", "127.0.0.2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update resource record 1234 under zone abc.com."))
				Expect(err.Error()).To(ContainSubstring("Failed to update resource record 5678 under zone abc.com."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record edit with different parameters", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--data", "127.0.0.2")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=900."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--by-record", "ftp", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--by-id", "1234", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "--by-id", "5678", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
		})
	})
})
