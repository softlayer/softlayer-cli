package ticket_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ticket create", func() {
	var (
		fakeUI            *terminal.FakeUI
		cliCommand        *ticket.CreateStandardTicketCommand
		fakeSession       *session.Session
		slCommand         *metadata.SoftlayerCommand
		fakeTicketManager *testhelpers.FakeTicketManager
		ticket_id         int
	)
	ticket_id = 12345
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = ticket.NewCreateStandardTicketCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeTicketManager = new(testhelpers.FakeTicketManager)
		cliCommand.TicketManager = fakeTicketManager
	})

	Describe("Ticket create", func() {
		Context("ticket create", func() {
			It("Success: inline body", func() {
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, `--body="Test Contents"`)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
			})

			It("Success: body from TicketManager.GetText()", func() {
				fakeTicketManager.GetTextReturns("Body goes here", nil)
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`)
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
				Expect(err).ToNot(HaveOccurred())
			})

			It("Success with all options", func() {
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ticket ID: 12345"}))
				Expect(err).ToNot(HaveOccurred())
			})

			It("Success with all options and virtual guest type", func() {
				fakeTicketManager.CreateStandardTicketReturns(&ticket_id, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", "--attachment-type=virtual", `--rootpwd="thisisapassword"`)
				Expect(fakeUI.Outputs()).To(ContainSubstring("Ticket ID: 12345"))
				Expect(err).ToNot(HaveOccurred())
			})

			It("Failure: No Subjet ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"requires the --subject-id"}))
			})

			It("Failure: No Title", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"requires the --title"}))
			})

			It("Failure: No body", func() {
				fakeTicketManager.GetTextReturns("nil", errors.New("Body Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Body Error"}))
			})

			It("Failure: No root password", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"Root password must be provided"}))
			})

			It("Failure: API failre", func() {
				fakeTicketManager.CreateStandardTicketReturns(nil, errors.New("API Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--subject-id=1522", `--title="Test Title"`, "--attachment=111111", `--rootpwd="thisisapassword"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstrings([]string{"API Error"}))

			})
		})
	})
})
