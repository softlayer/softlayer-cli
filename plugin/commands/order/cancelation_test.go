package order_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("order cancelation", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *order.CancelationCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = order.NewCancelationCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("order cancelation", func() {
		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return no error", func() {
			It("List order cancelations", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Case Number"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("153572280"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Number"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Approved"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Requested by"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("UserTest UserLastName"))
			})

			It("List order cancelations in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Case Number": "153572280",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Number Of Items Cancelled": "1",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Status": "Approved",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Requested by": "UserTest UserLastName"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})

	})
})
