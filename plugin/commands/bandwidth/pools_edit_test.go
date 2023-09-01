package bandwidth_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
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
		Context("Bandwidth Pool cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
		Context("Bandwidth Pool edit with correct Bandwidth Pool associated with id and --name ibm-internal-test", func() {
			BeforeEach(func() {
				fakeBandwidthManager.EditBandwidthReturns(true, errors.New("The Bandwidth Pool associated with Id 12345678 was edited successfully."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345678", "--name", "ibm-internal-test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The Bandwidth Pool associated with Id 12345678 was edited successfully."))
			})
		})
	})
})
