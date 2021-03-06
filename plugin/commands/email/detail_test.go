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

var _ = Describe("Email list Detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		cmd                *email.DetailCommand
		cliCommand         cli.Command
		fakeSession        *session.Session
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = email.NewDetailCommand(fakeUI, fakeEmailManager)
		cliCommand = cli.Command{
			Name:        email.DetailMetaData().Name,
			Description: email.DetailMetaData().Description,
			Usage:       email.DetailMetaData().Usage,
			Flags:       email.DetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Email detail", func() {
		Context("Email detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Email ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Email  detail, correct use", func() {
			It("return email  detail", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name               Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id                 295324"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username           test.test2@ie.ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email address      test.test3@ie.ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Create date        2020-07-06T16:29:11Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Category code      network_message_delivery"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Description        Free Package"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Type description   Delivery of messages through e-mail"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Type               EMAIL"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Vendor             SENDGRID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Statistics         Delivered   Requests   Bounces   Opens   Clicks   Spam reports"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0           0          0         0       0        0"))
				
			})
			It("return email  detail in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Id",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "295324"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Username",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "test.test2@ie.ibm.com"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Type description",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "Delivery of messages through e-mail"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Statistics",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "["`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
