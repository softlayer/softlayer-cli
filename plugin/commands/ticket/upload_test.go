package ticket_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket upload", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.UploadFileTicketCommand
		cliCommand        cli.Command
	)
	BeforeEach(func() {
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewUploadFileTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        ticket.TicketUploadMetaData().Name,
			Description: ticket.TicketUploadMetaData().Description,
			Usage:       ticket.TicketUploadMetaData().Usage,
			Flags:       ticket.TicketUploadMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket upload", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("This command requires two arguments."))
			})
			It("Without Ticket id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires two arguments."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("The ticket id must be a number."))
			})
			It("With an invalid ticket id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde", "/home/user/screenshot.png")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The ticket id must be a number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(errors.New("Error: SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
			It("With ticket id that does not exist", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "/home/user/screenshot.png")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Error: SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123456'. (HTTP 404)"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeTicketManager.AttachFileToTicketReturns(nil)
			})
			It("With valid parameters", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "/home/user/screenshot.png")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
