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
			Name:        metadata.TicketListMetaData().Name,
			Description: metadata.TicketListMetaData().Description,
			Usage:       metadata.TicketListMetaData().Usage,
			Flags:       metadata.TicketListMetaData().Flags,
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
	})
})
