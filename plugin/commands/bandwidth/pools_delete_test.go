package bandwidth_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/bandwidth"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Bandwidth Pool delete", func() {
	var (
		fakeUI               *terminal.FakeUI
		cliCommand           *bandwidth.DeleteCommand
		fakeSession          *session.Session
		slCommand            *metadata.SoftlayerCommand
		fakeBandwidthManager *testhelpers.FakeBandwidthManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeBandwidthManager = new(testhelpers.FakeBandwidthManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = bandwidth.NewDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.BandwidthManager = fakeBandwidthManager
	})

	Describe("Bandwidth Pool delete", func() {
		Context("Bandwidth Pool delete without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(""))
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))

			})
		})
		Context("Bandwidth Pool delete with wrong bandwidth id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'IDENTIFIER'. It must be a positive integer."))
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id but id not found", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeletePoolReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound"))
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id but server API call fails", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeletePoolReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeletePoolReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bandwidth pool 12345678 was deleted."))
			})
		})
	})
})
