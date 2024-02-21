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

var _ = Describe("ticket update", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.UpdateTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewUpdateTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket update", func() {
		Context("ticket update", func() {
			It("return succ", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return succ 2", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688", "This is a test.")
				Expect(err).ToNot(HaveOccurred())
			})

			It("Editor Error", func() {
				fakeTicketManager.GetTextReturns("nil", errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767699")
				Expect(err).To(HaveOccurred())
			})

			It("Invalid Usage", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
			})

			It("Invalid Usage, not a number", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Hello.")
				Expect(err).To(HaveOccurred())
			})

			It("API Error", func() {
				fakeTicketManager.AddUpdateReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688", "Tested Update")
				Expect(err).To(HaveOccurred())
			})
		})

	})
})
