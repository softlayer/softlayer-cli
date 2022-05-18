package email_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Edit email", func() {
	var (
		fakeUI           *terminal.FakeUI
		cmd              *email.EditCommand
		cliCommand       cli.Command
		fakeSession      *session.Session
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = email.NewEditCommand(fakeUI, fakeEmailManager)
		cliCommand = cli.Command{
			Name:        email.EditMetaData().Name,
			Description: email.EditMetaData().Description,
			Usage:       email.EditMetaData().Usage,
			Flags:       email.EditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Email edit", func() {
		Context("Email edit, Invalid Usage", func() {
			It("Send command without emailID", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
			It("Send command with bad emailID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'email ID'. It must be a positive integer."))
			})
			It("Send command without any flag", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Please pass at least one of the flags."))
			})
		})

		Context("Email edit, correct use", func() {
			It("return ok email account updated", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--username", "newusername", "--password", "xxxxxxxxxxxx")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(""))
			})
			It("return ok email account updated", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--email", "newemail@test.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(""))
			})
		})
	})
})
