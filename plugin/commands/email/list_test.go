package email_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Email list Email", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *email.ListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = email.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.EmailManager = fakeEmailManager
	})

	Describe("Email list", func() {
		Context("Email list, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Email list, correct use", func() {
			It("return email list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("295324"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.test2@ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Delivery of messages through e-mail"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SENDGRID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))

			})
		})
	})
})
