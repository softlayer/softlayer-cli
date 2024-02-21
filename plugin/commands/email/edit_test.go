package email_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Edit email", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *email.EditCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = email.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.EmailManager = fakeEmailManager
	})

	Describe("Email edit", func() {
		Context("Email edit, Invalid Usage", func() {
			It("Send command without emailID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
			It("Send command with bad emailID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'email ID'. It must be a positive integer."))
			})
			It("Send command without any flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Please pass at least one of the flags."))
			})
		})

		Context("Email edit, correct use", func() {
			It("return ok email account updated", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--username", "newusername", "--password", "xxxxxxxxxxxx")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email account 123 was updated."))
			})
			It("return ok email account updated", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--email", "newemail@test.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email address 123 was updated."))
			})
		})
	})
})
