package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account shows Summary ", func() {
	var (
		fakeUI             *terminal.FakeUI
		cmd                *account.SummaryCommand
		cliCommand         cli.Command
		fakeSession        *session.Session
		fakeAccountManager managers.AccountManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = managers.NewAccountManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = account.NewSummaryCommand(fakeUI, fakeAccountManager)
		cliCommand = cli.Command{
			Name:        account.SummaryMetaData().Name,
			Description: account.SummaryMetaData().Description,
			Usage:       account.SummaryMetaData().Usage,
			Flags:       account.SummaryMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Account summary", func() {
		Context("Account summary, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account summary, correct use", func() {
			It("return account summary", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Account Snapshot"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                      Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company Name              IBM Cloud IaaS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Balance                   275246.130000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Upcoming Invoice          3052.870000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Image Templates           43"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dedicated Hosts           2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware                  21"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Virtual Guests            55"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Domains                   48"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Network Storage Volumes   246"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Open Tickets              6"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Network Vlans             96"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Subnets                   103"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Users                     14"))
				
			})
			It("return account summary in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Account Snapshot":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Company Name","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "IBM Cloud IaaS""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Balance","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "275246.130000""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Upcoming Invoice","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "3052.870000""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(``))
				Expect(fakeUI.Outputs()).To(ContainSubstring(``))
				Expect(fakeUI.Outputs()).To(ContainSubstring(``))
				Expect(fakeUI.Outputs()).To(ContainSubstring(``))
				
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
