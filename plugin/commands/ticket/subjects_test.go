package ticket_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket subjects", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.SubjectsTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewSubjectsTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket subjects", func() {
		Context("ticket subjects", func() {
			It("Normal command call", func() {
				ticketId := 12345
				ticketName := "TestSubject"
				returnData := []datatypes.Ticket_Subject{
					datatypes.Ticket_Subject{
						Id:   &ticketId,
						Name: &ticketName,
					},
				}
				fakeTicketManager.GetSubjectsReturns(&returnData, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("TestSubject"))
			})

			It("return fail 1", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "argument")
				Expect(err).To(HaveOccurred())
			})

			It("API Failure", func() {
				fakeTicketManager.GetSubjectsReturns(nil, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
