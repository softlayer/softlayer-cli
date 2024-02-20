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

var _ = Describe("Account list BillingItems", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.BillingItemsCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewBillingItemsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account events", func() {
		Context("Account events, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account events, correct use", func() {
			It("return account events", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--category=ssl_certificate", "--create=2016-01-20", "--ordered=TestName")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("81336973"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-01-20T17:00:19Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.00"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ssl_certificate"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("TestName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RapidSSL - 1 year"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("techbabble.xyz"))

			})
			It("return account events in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
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
