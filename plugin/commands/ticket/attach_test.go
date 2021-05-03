package ticket_test

import (
	"strings"
	"errors"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket attach", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeTicketManager *testhelpers.FakeTicketManager
		cmd               *ticket.AttachDeviceTicketCommand
		cliCommand        cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cmd = ticket.NewAttachDeviceTicketCommand(fakeUI, fakeTicketManager)
		cliCommand = cli.Command{
			Name:        metadata.TicketAttachMetaData().Name,
			Description: metadata.TicketAttachMetaData().Description,
			Usage:       metadata.TicketAttachMetaData().Usage,
			Flags:       metadata.TicketAttachMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Ticket attach", func() {
		Context("ticket attach", func() {
			It("return succ 1", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688", "--hardware=111111")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return succ 2", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688", "--virtual=222222")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return error 1", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})

			It("return error 2", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688", "--hardware=111111", "--virtual=222222")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: hardware and virtual flags cannot be set at the same time.")).To(BeTrue())
			})

			It("Error: API Error", func() {
				fakeTicketManager.AttachDeviceToTicketReturns(errors.New("API ERROR"))
				err := testhelpers.RunCommand(cliCommand, "76767699", "--hardware=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Error: API ERROR"}))
			})

			It("return error 5", func() {
				err := testhelpers.RunCommand(cliCommand, "hello", "--hardware=111111")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage:")).To(BeTrue())
			})

			It("return error 6", func() {
				err := testhelpers.RunCommand(cliCommand, "76767688")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage:")).To(BeTrue())
			})
		})
	})
})
