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

var _ = Describe("Account list InvoiceDetail", func() {
	var (
		fakeUI             *terminal.FakeUI
		cmd                *account.InvoiceDetailCommand
		cliCommand         cli.Command
		fakeSession        *session.Session
		fakeAccountManager managers.AccountManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = managers.NewAccountManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = account.NewInvoiceDetailCommand(fakeUI, fakeAccountManager)
		cliCommand = cli.Command{
			Name:        account.InvoiceDetailMetaData().Name,
			Description: account.InvoiceDetailMetaData().Description,
			Usage:       account.InvoiceDetailMetaData().Usage,
			Flags:       account.InvoiceDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Account invoice detail", func() {
		Context("Account invoice detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invoice ID is required."))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: The invoice ID has to be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account invoice detail, correct use", func() {
			It("return account invoice detail", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Invoice: 123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item Id        Category   Description                                                                                              Single   Monthly   Create Date            Location"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789      Server     Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-for-test-with-adition..   10.23    20.34     2022-04-04T05:10:20Z   mex01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789123   server     Dual E5-2690 v3 (12 Cores, 2.60 GHz) (test-vs.support2.com)                                              11.23    21.12     2022-04-04T05:10:21Z   ams01"))
			})
			It("return account invoice detail with additionals details", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--details")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Invoice: 123"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item Id        Category           Description                                                                                              Single   Monthly   Create Date            Location"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789      Server             Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-for-test-with-adition..   10.23    20.34     2022-04-04T05:10:20Z   mex01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring(">>>            Second Processor   Intel Xeon (12 Cores, 2.40 GHz)                                                                          10.23    20.34     ---                    ---"))
				Expect(fakeUI.Outputs()).To(ContainSubstring(">>>            Operating System   Virtual (up to 1Gbps)                                                                                    10.23    20.34     ---                    ---"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789123   server             Dual E5-2690 v3 (12 Cores, 2.60 GHz) (test-vs.support2.com)                                              11.23    21.12     2022-04-04T05:10:21Z   ams01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring(">>>            Second Processor   Intel Xeon (12 Cores, 2.40 GHz)                                                                          11.23    21.12     ---                    ---"))
				Expect(fakeUI.Outputs()).To(ContainSubstring(">>>            Operating System   Virtual (up to 1Gbps)                                                                                    11.23    21.12     ---                    ---"))
			})
			It("return account invoice detail in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Invoice: 123": "["`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Item Id": "123456789","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Category": "Server","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Description": "Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-for-test-with-adition..","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
