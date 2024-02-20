package bandwidth_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/bandwidth"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Bandwidth Pool edit", func() {
	var (
		fakeUI               *terminal.FakeUI
		cliCommand           *bandwidth.EditCommand
		fakeSession          *session.Session
		slCommand            *metadata.SoftlayerCommand
		fakeBandwidthManager *testhelpers.FakeBandwidthManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeBandwidthManager = new(testhelpers.FakeBandwidthManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = bandwidth.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.BandwidthManager = fakeBandwidthManager
	})

	Describe("Bandwidth Pool edit", func() {
		Context("Bandwidth Pool invalid usage", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))

			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "name" not set`))
			})
		})
		Context("Happy Path", func() {
			BeforeEach(func() {
				fakeBandwidthManager.EditPoolReturns(true, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678", "--name", "ibm-internal-test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bandwidth pool 12345678 was edited successfully."))
			})
		})
		Context("API Errors", func() {
			BeforeEach(func() {
				fakeBandwidthManager.EditPoolReturns(false, errors.New("API ERROR"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678", "--name", "ibm-internal-test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("API ERROR"))
			})
		})
	})
})
