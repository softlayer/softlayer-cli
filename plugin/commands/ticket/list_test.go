package ticket_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket list", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.ListTicketCommand
		cliCommand        cli.Command
	)
	fakeTicketManager = new(testhelpers.FakeTicketManager)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewListTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        ticket.TicketListMetaData().Name,
			Description: ticket.TicketListMetaData().Description,
			Usage:       ticket.TicketListMetaData().Usage,
			Flags:       ticket.TicketListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket list", func() {
		Context("ticket list", func() {
			It("return succ 1", func() {
				err := testhelpers.RunCommand(cliCommand, "--open")
				Expect(err).ToNot(HaveOccurred())
			})
			It("return succ 2", func() {
				err := testhelpers.RunCommand(cliCommand, "--closed")
				Expect(err).ToNot(HaveOccurred())
			})
			It("return succ 3", func() {
				err := testhelpers.RunCommand(cliCommand, "--open", "--closed")
				Expect(err).ToNot(HaveOccurred())
			})
			It("return fail 1", func() {
				err := testhelpers.RunCommand(cliCommand, "argument")
				Expect(err).To(HaveOccurred())
			})

			It("API Failure", func() {
				fakeTicketManager.ListOpenTicketsReturns(nil, errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "--open")
				Expect(err).To(HaveOccurred())
			})

		})

		Context("Return no error", func() {
			tickets := []datatypes.Ticket{}
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				lastEdited, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:59Z")
				tickets = []datatypes.Ticket{
					datatypes.Ticket{
						Id: sl.Int(111111),
						AssignedUser: &datatypes.User_Customer{
							FirstName: sl.String("Juan"),
							LastName:  sl.String("Perez"),
						},
						CreateDate:   sl.Time(created),
						LastEditDate: sl.Time(lastEdited),
						Title:        sl.String("My ticket"),
						Status: &datatypes.Ticket_Status{
							Name: sl.String("Open"),
						},
						TotalUpdateCount: sl.Int(2),
						Priority:         sl.Int(0),
					},
				}
				fakeTicketManager.ListOpenTicketsReturns(tickets, nil)
			})

			It("List ticket", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Juan Perez"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("My ticket"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:59Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Open"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
		})
	})
})
