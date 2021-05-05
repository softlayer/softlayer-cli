package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware reboot", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.RebootCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewRebootCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        metadata.HardwarePowerRebootMetaData().Name,
			Description: metadata.HardwarePowerRebootMetaData().Description,
			Usage:       metadata.HardwarePowerRebootMetaData().Usage,
			Flags:       metadata.HardwarePowerRebootMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware reboot", func() {
		Context("hardware reboot without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("hardware reboot with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware reboot with wrong parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--soft", "--hard")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Can only specify either --hard or --soft."))
			})
		})

		Context("hardware reboot with correct ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will reboot hardware server: 1234. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("hardware reboot with correct ID but server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.RebootReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to reboot hardware server: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("hardware reboot with correct ID ", func() {
			BeforeEach(func() {
				fakeHardwareManager.RebootReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--soft")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234", "--soft")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--hard")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234", "--hard")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server: 1234 was rebooted."))
			})
		})
	})
})
