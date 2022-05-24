package licenses_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Licenses list Cancel Item", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeLicensesManager *testhelpers.FakeLicensesManager
		cmd                 *licenses.CancelItemCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLicensesManager = new(testhelpers.FakeLicensesManager)
		cmd = licenses.NewCancelItemCommand(fakeUI, fakeLicensesManager)
		cliCommand = cli.Command{
			Name:        licenses.CancelItemMetaData().Name,
			Description: licenses.CancelItemMetaData().Description,
			Usage:       licenses.CancelItemMetaData().Usage,
			Flags:       licenses.CancelItemMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Licenses cancel item", func() {
		Context("Licenses cancel item, Invalid Usage", func() {
			It("Set command without any datacenter and keyName", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`This command requires one argument.`))
			})
		})

		Context("Licenses cancel item, correct use", func() {
			It("return licenses cancel item", func() {
				err := testhelpers.RunCommand(cliCommand, "XXX_XXX_XXX", "--immediate")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("License: XXX_XXX_XXX was cancelled."))
			})
		})

		Context("Licenses cancel errors", func() {
			It("return license error", func() {
				fakeLicensesManager.CancelItemReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
				err := testhelpers.RunCommand(cliCommand, "XXX_XXX_XXX")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find license with key: XXX_XXX_XXX."))
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return license error", func() {
				fakeLicensesManager.CancelItemReturns(errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "XXX_XXX_XXX")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel license: XXX_XXX_XXX."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
	})
})
