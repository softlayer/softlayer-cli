package bandwidth_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
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
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
		Context("Bandwidth Pool delete with wrong bandwidth id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Bandwidth Pool ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id but id not found", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeleteBandwidthReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound")).To(BeTrue())
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id but server API call fails", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeleteBandwidthReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Bandwidth Pool delete with correct bandwidth id", func() {
			BeforeEach(func() {
				fakeBandwidthManager.DeleteBandwidthReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"BandwidthPool associated with Id 12345678 was deleted."}))
			})
		})
	})
})
