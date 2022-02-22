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

var _ = Describe("ticket detach", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.DetachDeviceTicketCommand
		cliCommand        cli.Command
	)
	BeforeEach(func() {
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewDetachDeviceTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        ticket.TicketDetachMetaData().Name,
			Description: ticket.TicketDetachMetaData().Description,
			Usage:       ticket.TicketDetachMetaData().Usage,
			Flags:       ticket.TicketDetachMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket detach", func() {
		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.RemoveDeviceFromTicketReturns(errors.New("This command requires one argument."))
			})
			It("Argument is not set", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.RemoveDeviceFromTicketReturns(errors.New("The ticket id must be a number."))
			})
			It("An invalid ticket id is set", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The ticket id must be a number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.RemoveDeviceFromTicketReturns(errors.New("either the hardware or virtual flag must be set."))
			})
			It("Hardware or virtual flag are not set", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("either the hardware or virtual flag must be set."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.RemoveDeviceFromTicketReturns(errors.New("hardware and virtual flags cannot be set at the same time."))
			})
			It("Hardware and virtual flag are set at the same time", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--hardware=987654", "--virtual=876543")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("hardware and virtual flags cannot be set at the same time."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTicketManager.RemoveDeviceFromTicketReturns(errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("Ticket id that does not exist is set", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--hardware=987654")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Remove hardware", func() {
			It("Return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--hardware=987654")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

		Context("Remove virtual server", func() {
			It("Return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--virtual=987654")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
