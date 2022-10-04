package hardware_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware powercycle", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.PowerCycleCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewPowerCycleCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware power-cycle", func() {
		Context("hardware powercycle without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
		Context("hardware powercycle with wrong ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Hardware server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("hardware powercycle with correct ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will power cycle hardware server: 1234. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("hardware powercycle with correct ID but server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.PowerCycleReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to power cycle hardware server: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("hardware powercycle with correct ID ", func() {
			BeforeEach(func() {
				fakeHardwareManager.PowerCycleReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was power cycle."))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was power cycle."))
			})
		})
	})
})
