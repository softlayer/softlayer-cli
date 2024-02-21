package ticket_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket upload", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.UploadFileTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewUploadFileTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket upload", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("This command requires two arguments."))
			})
			It("Without Ticket id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires two arguments."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("The ticket id must be a number."))
			})
			It("With an invalid ticket id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "/home/user/screenshot.png")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The ticket id must be a number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("Error: SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
			It("With ticket id that does not exist", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "/home/user/screenshot.png")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Error: SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(nil)
			})
			It("With valid parameters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "/home/user/screenshot.png")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
