package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware reflash-firmware", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.ReflashFirmwareCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewReflashFirmwareCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware reflash-firmware", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("abcde")
			})
			It("Confirm action with invalid input", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeHardwareManager.CreateFirmwareReflashTransactionReturns(false, errors.New("Failed to reflash firmware."))
			})
			It("Failed to reflash firmware", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to reflash firmware."))
			})
		})

		Context("Retun no error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("n")
			})
			It("Abort command", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeHardwareManager.CreateFirmwareReflashTransactionReturns(true, nil)
			})
			It("Reflash firmware", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully device firmware reflashed"))
			})
		})
	})
})
