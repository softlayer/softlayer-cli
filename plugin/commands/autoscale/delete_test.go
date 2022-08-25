package autoscale_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale delete", func() {
	var (
		fakeUI               *terminal.FakeUI
		cliCommand           *autoscale.DeleteCommand
		fakeSession          *session.Session
		slCommand            *metadata.SoftlayerCommand
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = autoscale.NewDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AutoScaleManager = fakeAutoScaleManager
	})

	Describe("autoscale delete", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Autoscale Group ID'. It must be a positive integer."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.DeleteReturns(false, errors.New("Failed to delete Auto Scale Group."))
			})
			It("Failed delete scale group", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete Auto Scale Group."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("abcde")
			})
			It("Cancel with invalid input", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.DeleteReturns(true, nil)
				fakeUI.Inputs("y")
			})

			It("Delete scale group", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Auto Scale Group was deleted successfully"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.DeleteReturns(true, nil)
			})

			It("Delete scale group without confirmation", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Auto Scale Group was deleted successfully"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("n")
			})

			It("Cancel", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
	})
})
