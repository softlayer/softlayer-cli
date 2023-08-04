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

var _ = Describe("Email list Detail", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *email.DetailCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = email.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.EmailManager = fakeEmailManager
	})

	Describe("Email detail", func() {
		Context("Email detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Email ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Email  detail, correct use", func() {
			It("return email  detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("295324"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.test2@ie.ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.test3@ie.ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2020-07-06T16:29:11Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("network_message_delivery"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Free Package"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Delivery of messages through e-mail"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("EMAIL"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SENDGRID"))

			})
		})
	})
})
