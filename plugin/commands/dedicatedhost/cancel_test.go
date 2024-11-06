package dedicatedhost_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host cancel", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *dedicatedhost.CancellHostCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		FakeDedicatedhostManager *testhelpers.FakeDedicatedHostManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dedicatedhost.NewCancelHostCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedHostManager)
		cliCommand.DedicatedHostManager = FakeDedicatedhostManager
	})

	Describe("Dedicatedhost cancel usage errors", func() {
		Context("Dedicatedhost cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
	})
	Describe("Dedicatedhost cancel usage", func() {
		Context("Dedicatedhost cancel Happy Path", func() {
			It("Succss", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dedicated Host 12345 was cancelled"))
			})
		})
		Context("Dedicatedhost cancel API errors", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.DeleteHostReturns(errors.New("API ERROR"))
			})
			It("Handle API error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("API ERROR"))
			})
		})
	})
})
