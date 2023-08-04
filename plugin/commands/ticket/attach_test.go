package ticket_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket attach", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.AttachDeviceTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewAttachDeviceTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket attach", func() {
		Context("ticket attach", func() {
			It("return succ 1", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688", "--hardware=111111")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return succ 2", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688", "--virtual=222222")
				Expect(err).ToNot(HaveOccurred())
			})

			It("return error 1", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})

			It("return error 2", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688", "--hardware=111111", "--virtual=222222")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: hardware and virtual flags cannot be set at the same time."))
			})

			It("Error: API Error", func() {
				fakeTicketManager.AttachDeviceToTicketReturns(errors.New("API ERROR"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767699", "--hardware=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("API ERROR"))
			})

			It("return error 5", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "hello", "--hardware=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage:"))
			})

			It("return error 6", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "76767688")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage:"))
			})
		})
	})
})
