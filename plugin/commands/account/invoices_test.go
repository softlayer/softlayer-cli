package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list Invoices", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.InvoicesCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewInvoicesCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account invoices", func() {
		Context("Account invoices, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account invoices, correct use", func() {
			It("return account invoices", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--limit", "10", "--closed", "--all")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id         Created                Type   Status   Starting Balance   Ending Balance   Invoice Amount   Items"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("76602936   2021-11-24T21:07:42Z   NEW    OPEN     264111.300000      264111.300000    0.000000         14"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("77186102   2021-12-10T13:44:59Z   NEW    CLOSED   266803.650000      266803.650000    0.000000         3"))
			})
			It("return account invoices in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "76602936",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Created": "2021-11-24T21:07:42Z",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "OPEN",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
