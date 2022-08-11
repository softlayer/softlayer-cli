package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware reflash-firmware", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.ReflashFirmwareCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewReflashFirmwareCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        hardware.HardwareReflashFirmwareMetaData().Name,
			Description: hardware.HardwareReflashFirmwareMetaData().Description,
			Usage:       hardware.HardwareReflashFirmwareMetaData().Usage,
			Flags:       hardware.HardwareReflashFirmwareMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware reflash-firmware", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("abcde")
			})
			It("Confirm action with invalid input", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeHardwareManager.CreateFirmwareReflashTransactionReturns(false, errors.New("Failed to reflash firmware."))
			})
			It("Failed to reflash firmware", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to reflash firmware."))
			})
		})

		Context("Retun no error", func() {
			BeforeEach(func() {
				fakeUI.Inputs("n")
			})
			It("Abort command", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("Retun error", func() {
			BeforeEach(func() {
				fakeHardwareManager.CreateFirmwareReflashTransactionReturns(true, nil)
			})
			It("Reflash firmware", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully device firmware reflashed"))
			})
		})
	})
})
