package licenses_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Licenses list Cancel Item", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *licenses.CancelItemCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeLicensesManager *testhelpers.FakeLicensesManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = licenses.NewCancelItemCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLicensesManager = new(testhelpers.FakeLicensesManager)
		cliCommand.LicensesManager = fakeLicensesManager
	})

	Describe("Licenses cancel item", func() {
		Context("Licenses cancel item, Invalid Usage", func() {
			It("Set command without any datacenter and keyName", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`This command requires one argument`))
			})
		})

		Context("Licenses cancel item, correct use", func() {
			It("return licenses cancel item", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "XXX_XXX_XXX", "--immediate")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("License: XXX_XXX_XXX was cancelled."))
			})
		})

		Context("Licenses cancel errors", func() {
			It("return license error", func() {
				fakeLicensesManager.CancelItemReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "XXX_XXX_XXX")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find license with key: XXX_XXX_XXX."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return license error", func() {
				fakeLicensesManager.CancelItemReturns(errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "XXX_XXX_XXX")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel license: XXX_XXX_XXX."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
	})
})
