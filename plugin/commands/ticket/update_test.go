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
)

var _ = Describe("ticket update", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.UpdateTicketCommand
		cliCommand        cli.Command
	)
	fakeTicketManager = new(testhelpers.FakeTicketManager)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		cmd = ticket.NewUpdateTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        metadata.TicketUpdataMetaData().Name,
			Description: metadata.TicketUpdataMetaData().Description,
			Usage:       metadata.TicketUpdataMetaData().Usage,
			Flags:       metadata.TicketUpdataMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket update", func() {
		Context("ticket update", func() {
			It("return succ", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return succ 2", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688", "This is a test.")
				Expect(err).ToNot(HaveOccurred())
			})

			It("Editor Error", func() {
				fakeTicketManager.GetTextReturns("nil", errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "76767699")
				Expect(err).To(HaveOccurred())
			})

			It("Invalid Usage", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
			})

			It("Invalid Usage, not a number", func() {
				err := testhelpers.RunCommand(cliCommand, "Hello.")
				Expect(err).To(HaveOccurred())
			})

			It("API Error", func() {
				fakeTicketManager.AddUpdateReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCommand(cliCommand, "76767688", "Tested Update")
				Expect(err).To(HaveOccurred())
			})
		})

	})
})
