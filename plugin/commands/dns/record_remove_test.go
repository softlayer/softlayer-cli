package dns_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Record remove", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.RecordRemoveCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewRecordRemoveCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Record remove", func() {
		Context("Record remove without record ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Record remove with record ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Record ID'. It must be a positive integer."))
			})
		})

		Context("Record remove with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.DeleteResourceRecordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete resource record: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))

			})
		})

		Context("Record remove with record not found", func() {
			BeforeEach(func() {
				fakeDNSManager.DeleteResourceRecordReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find resource record with ID: 1234."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound"))

			})
		})

		Context("Record remove", func() {
			BeforeEach(func() {
				fakeDNSManager.DeleteResourceRecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Resource record 1234 was removed."}))
			})
		})
	})
})
