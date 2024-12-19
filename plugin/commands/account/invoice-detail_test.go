package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list InvoiceDetail", func() {
	var (
		fakeUI      	*terminal.FakeUI
		cliCommand  	*account.InvoiceDetailCommand
		fakeSession 	*session.Session
		slCommand  		*metadata.SoftlayerCommand
		fakeHandler     *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewInvoiceDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
    AfterEach(func() {
        fakeHandler.ClearApiCallLogs()
        fakeHandler.ClearErrors()
    })

	Describe("Account invoice detail", func() {
		Context("Account invoice detail, Invalid Usage", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Set command with id like letters", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Invoice ID'. It must be a positive integer."))
			})
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account invoice detail, correct use", func() {
			It("return account invoice detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Item Id        Category   Description                                                                           Single   Monthly   Create Date   Location"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789      Server     Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-f...   22.59    35.26     2022-04-04    mex01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789123   server     Dual E5-2690 v3 (12 Cores, 2.60 GHz) (test-vs.support2.com)                           23.81    36.04     2022-04-04    ams01"))
			})
			It("return account invoice detail with additionals details", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--details")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(Equal(
`Item Id        Category           Description                                                                           Single   Monthly   Create Date   Location
123456789      Server             Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-f...   22.59    35.26     2022-04-04    mex01
>>>            Server             Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-f...   10.23    20.34     ---           ---
>>>            Second Processor   Intel Xeon (12 Cores, 2.40 GHz)                                                       5.24     6.12      ---           ---
>>>            Operating System   Virtual (up to 1Gbps)                                                                 7.12     8.79      ---           ---
123456789123   server             Dual E5-2690 v3 (12 Cores, 2.60 GHz) (test-vs.support2.com)                           23.81    36.04     2022-04-04    ams01
>>>            server             Dual E5-2690 v3 (12 Cores, 2.60 GHz) (test-vs.support2.com)                           11.23    21.12     ---           ---
>>>            Second Processor   Intel Xeon (12 Cores, 2.40 GHz)                                                       5.35     6.23      ---           ---
>>>            Operating System   Virtual (up to 1Gbps)                                                                 7.23     8.68      ---           ---
`))
			})
			It("return account invoice detail in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Item Id": "123456789",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Category": "Server",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Description": "Dual Intel Xeon Silver 4210 (20 Cores, 2.20 GHz) (test-gpu.softlayer-community-f...",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Single": "22.59",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Monthly": "35.26",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Location": "mex01"`))
			})
		})
		Context("issues856", func() {
			It("Handle large int invoices", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "999")
				Expect(err).NotTo(HaveOccurred())
				output := fakeUI.Outputs()
				Expect(output).To(ContainSubstring("testlb-307608-dal13.lb.bluemix.net"))
			})
			It("Missing properties dont break", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "888")
				Expect(err).NotTo(HaveOccurred())
				output := fakeUI.Outputs()
				Expect(output).To(ContainSubstring("2020-05-04    None"))
				Expect(output).To(ContainSubstring("0.00     0.00      2020-05-04    tok02"))
			})
		})
	})
})
