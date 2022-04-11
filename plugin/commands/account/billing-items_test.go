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

var _ = Describe("Account list BillingItems", func() {
	var (
		fakeUI      *terminal.FakeUI
		cmd         *account.BillingItemsCommand
		cliCommand  cli.Command
		fakeSession *session.Session
		fakeAccountManager managers.AccountManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = managers.NewAccountManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = account.NewBillingItemsCommand(fakeUI, fakeAccountManager)
		cliCommand = cli.Command{
			Name:        account.BillingItemsMetaData().Name,
			Description: account.BillingItemsMetaData().Description,
			Usage:       account.BillingItemsMetaData().Usage,
			Flags:       account.BillingItemsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Account events", func() {
		Context("Account events, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account events, correct use", func() {
			It("return account events", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Billing Items"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id          Create Date            Cost   Category Code             Ordered By   Description                                            Notes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("81336973    2016-01-20T17:00:19Z   0.00   ssl_certificate           TestName     RapidSSL - 1 year                                      techbabble.xyz"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("933002170   2022-02-18T18:47:32Z   0.00   dedicated_virtual_hosts   testName2    virtualserver01-0c56.softlayer-internal-developmen..   -"))

			})
			It("return account events in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`Billing Items`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "81336973",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Create Date": "2016-01-20T17:00:19Z",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Cost": "0.00",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Ordered By": "TestName",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
