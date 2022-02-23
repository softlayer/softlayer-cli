package hardware_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware reload", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.ReloadCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewReloadCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        hardware.HardwareReloadMetaData().Name,
			Description: hardware.HardwareReloadMetaData().Description,
			Usage:       hardware.HardwareReloadMetaData().Usage,
			Flags:       hardware.HardwareReloadMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware reload", func() {
		Context("hardware reload without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("hardware reload with wrong ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Hardware server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("hardware reload with correct ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will reload operating system for hardware server: 1234. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("hardware reload with correct ID but server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.ReloadReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to reload operating system for hardware server: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("hardware reload with correct ID ", func() {
			BeforeEach(func() {
				fakeHardwareManager.ReloadReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to reload operating system for hardware server: 1234"))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to reload operating system for hardware server: 1234"))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234", "-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to reload operating system for hardware server: 1234"))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCommand(cliCommand, "1234", "-w")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to reload operating system for hardware server: 1234"))
			})
		})
	})
})
