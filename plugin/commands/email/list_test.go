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

var _ = Describe("Email list Email", func() {
	var (
		fakeUI           *terminal.FakeUI
		cmd              *email.ListCommand
		cliCommand       cli.Command
		fakeSession      *session.Session
		fakeEmailManager managers.EmailManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeEmailManager = managers.NewEmailManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = email.NewListCommand(fakeUI, fakeEmailManager)
		cliCommand = cli.Command{
			Name:        email.ListMetaData().Name,
			Description: email.ListMetaData().Description,
			Usage:       email.ListMetaData().Usage,
			Flags:       email.ListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Email list", func() {
		Context("Email list, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Email list, correct use", func() {
			It("return email list", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email information   Id       Username             Description                           Vendor"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("295324   test.test2@ibm.com   Delivery of messages through e-mail   SENDGRID"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email overview      Credit allowed   Credits remain   Credits overage   Credits used   Package        Reputation   Requests"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("25000            25000            0                 0              Free Package   100          56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Statistics          Delivered   Requests   Bounces   Opens   Clicks   Spam reports"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0           0          0         0       0        0"))

			})
			It("return email email in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Email information",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "295324","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Username": "test.test2@ibm.com","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Description": "Delivery of messages through e-mail","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Email overview",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Credit allowed": "25000","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Credits remain": "25000","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Statistics",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Delivered": "0","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Requests": "0","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
