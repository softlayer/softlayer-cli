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

var _ = Describe("Record add", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.RecordAddCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewRecordAddCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Record add", func() {
		Context("Record add with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires four arguments."))
			})
		})
		Context("Record add with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires four arguments."))
			})
		})

		Context("Record add with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get zone ID from zone name: abc.com."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record add with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create resource record under zone abc.com: type=a, record=ftp, data=127.0.0.1, ttl=7200."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "ftp", "a", "127.0.0.1", "--ttl", "3600")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create resource record under zone abc.com: type=a, record=ftp, data=127.0.0.1, ttl=3600."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Record add", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{
					Id:   sl.Int(1234),
					Type: sl.String("a"),
					Host: sl.String("ftp"),
					Data: sl.String("127.0.0.1"),
					Ttl:  sl.Int(7200),
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Created resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=7200."}))
			})
		})
		Context("Record add", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{
					Id:   sl.Int(1234),
					Type: sl.String("a"),
					Host: sl.String("ftp"),
					Data: sl.String("127.0.0.1"),
					Ttl:  sl.Int(3600),
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc.com", "ftp", "a", "127.0.0.1", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Created resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=3600."}))
			})
		})
	})
})
