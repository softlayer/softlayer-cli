package ticket_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket detail", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.DetailTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewDetailTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket detail", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.GetTicketReturns(datatypes.Ticket{}, errors.New("This command requires one argument."))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.GetTicketReturns(datatypes.Ticket{}, errors.New("The ticket id must be a positive non-zero number."))
			})
			It("Invalid ticket id is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "0")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The ticket id must be a positive non-zero number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.GetTicketReturns(datatypes.Ticket{}, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("Ticket id that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.GetTicketReturns(datatypes.Ticket{}, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("Ticket id that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.GetAllUpdatesReturns([]datatypes.Ticket_Update{}, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("Ticket id that does not exist is set", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Return no error", func() {
			fakeTicket := datatypes.Ticket{}
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				lastEdited, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:59Z")
				fakeTicket = datatypes.Ticket{
					Id:        sl.Int(123456),
					Title:     sl.String("My title"),
					Priority:  sl.Int(0),
					AccountId: sl.Int(278444),
					Status: &datatypes.Ticket_Status{
						Name: sl.String("Open"),
					},
					CreateDate:   sl.Time(created),
					LastEditDate: sl.Time(lastEdited),
					AssignedUser: &datatypes.User_Customer{
						FirstName: sl.String("Juan"),
						LastName:  sl.String("Perez"),
					},
				}
				fakeTicketManager.GetTicketReturns(fakeTicket, nil)
			})

			It("Get ticket with one update", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("My title"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Open"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:59Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Juan Perez"))
			})
		})

		Context("Return no error", func() {
			ticketUpdates := []datatypes.Ticket_Update{}
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				lastEdited, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:59Z")
				ticketUpdates = []datatypes.Ticket_Update{
					datatypes.Ticket_Update{
						CreateDate: sl.Time(created),
						Id:         sl.Int(111111),
						Entry:      sl.String("Entry 1"),
					},
					datatypes.Ticket_Update{
						CreateDate: sl.Time(lastEdited),
						Id:         sl.Int(222222),
						Entry:      sl.String("Entry 2"),
					},
				}
				fakeTicketManager.GetAllUpdatesReturns(ticketUpdates, nil)
			})

			It("Get ticket with two updates", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--count=2")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Entry 1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:59Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Entry 2"))
			})
		})
	})
})
