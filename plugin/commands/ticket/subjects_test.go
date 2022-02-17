package ticket_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket subjects", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.SubjectsTicketCommand
		cliCommand        cli.Command
	)
	fakeTicketManager = new(testhelpers.FakeTicketManager)

	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewSubjectsTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        ticket.TicketSubjectsMetaData().Name,
			Description: ticket.TicketSubjectsMetaData().Description,
			Usage:       ticket.TicketSubjectsMetaData().Usage,
			Flags:       ticket.TicketSubjectsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket subjects", func() {
		Context("ticket subjects", func() {
			It("Normal command call", func() {
				ticketId := 12345
				ticketName := "TestSubject"
				returnData := []datatypes.Ticket_Subject{
					datatypes.Ticket_Subject{
						Id: &ticketId,
						Name: &ticketName,	
					},
				}
				fakeTicketManager.GetSubjectsReturns(&returnData, nil)
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("TestSubject"))
			})

			It("return fail 1", func() {
				err := testhelpers.RunCommand(cliCommand, "argument")
				Expect(err).To(HaveOccurred())
			})

			It("API Failure", func() {
				fakeTicketManager.GetSubjectsReturns(nil, errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
