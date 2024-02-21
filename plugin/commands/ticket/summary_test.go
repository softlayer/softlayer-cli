package ticket_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket list", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.SummaryTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewSummaryTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket summary", func() {
		Context("ticket summary success", func() {
			It("return succ 1", func() {

				fakeReturn := managers.TicketSummary{
					Accounting: 0,
					Billing:    1,
					Sales:      2,
					Support:    3,
					Other:      4,
					Closed:     5,
					Open:       6,
				}
				fakeTicketManager.SummaryReturns(&fakeReturn, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("ticket summary failure", func() {
			It("return fail 1", func() {

				fakeReturn := managers.TicketSummary{}
				fakeTicketManager.SummaryReturns(&fakeReturn, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
