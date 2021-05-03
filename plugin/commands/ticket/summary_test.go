package ticket_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
)

var _ = Describe("ticket list", func() {
	var (
		fakeUI            *terminal.FakeUI
		FakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.SummaryTicketCommand
		cliCommand        cli.Command
	)
	
	BeforeEach(func() {
		FakeTicketManager = new(testhelpers.FakeTicketManager)
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewSummaryTicketCommand(fakeUI, FakeTicketManager)
		cliCommand = cli.Command{
			Name:        metadata.TicketSummaryMetaData().Name,
			Description: metadata.TicketSummaryMetaData().Description,
			Usage:       metadata.TicketSummaryMetaData().Usage,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket summary", func() {
		Context("ticket summary success", func() {
			It("return succ 1", func() {

				fakeReturn := managers.TicketSummary{
					Accounting: 0,
					Billing: 1,
					Sales: 2,
					Support: 3,
					Other: 4,
					Closed: 5,
					Open: 6,
				}
				FakeTicketManager.SummaryReturns(&fakeReturn, nil)
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
			})

		})
		Context("ticket summary failure", func() {
			It("return fail 1", func() {

				fakeReturn := managers.TicketSummary{}
				FakeTicketManager.SummaryReturns(&fakeReturn, errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
			})

		})
	})
})
