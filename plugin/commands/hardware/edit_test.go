package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware edit", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.EditCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewEditCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        hardware.HardwareEditMetaData().Name,
			Description: hardware.HardwareEditMetaData().Description,
			Usage:       hardware.HardwareEditMetaData().Usage,
			Flags:       hardware.HardwareEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware edit", func() {
		Context("hardware edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("hardware edit with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})
		Context("hardware edit with both -u and -F", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata", "-F", "/tmp/datafile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-u|--userdata] is not allowed with [-F|--userfile]."))
			})
		})
		Context("hardware edit with wrong public speed", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "9")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
			})
		})
		Context("hardware edit with wrong private speed", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "9")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
			})
		})

		Context("hardware edit with hostname fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the hostname/domain/note of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-H", "vs-abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the hostname/domain/note of hardware server: 1234."))
			})
		})
		Context("hardware edit with hostname succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The hostname of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-H", "vs-abc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The hostname of hardware server: 1234 was updated."))
			})
		})

		Context("hardware edit with domain fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the hostname/domain of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-D", "wilma.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the hostname/domain of hardware server: 1234."))
			})
		})
		Context("hardware edit with domain succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The domain of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-D", "wilma.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The domain of hardware server: 1234 was updated."))
			})
		})
		Context("hardware edit with userdata fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the user data of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the user data of hardware server: 1234."))
			})
		})
		Context("hardware edit with user data succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The user data of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-u", "mydata")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The user data of hardware server: 1234 was updated."))
			})
		})
		Context("hardware edit with tags fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the tags of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--tag", "mytags")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the tags of hardware server: 1234."))
			})
		})
		Context("hardware edit with tag succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The tags of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--tag", "mytags")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The tags of hardware server: 1234 was updated."))
			})
		})
		Context("hardware edit with public-speed fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the public network speed of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "1000")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the public network speed of hardware server: 1234."))
			})
		})
		Context("hardware edit with public-speed succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The public network speed of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--public-speed", "1000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The public network speed of hardware server: 1234 was updated."))
			})
		})
		Context("hardware edit with private-speed fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{false}, []string{"Failed to update the private network speed of hardware server: 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "1000")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update the private network speed of hardware server: 1234."))
			})
		})
		Context("hardware edit with private-speed succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.EditReturns([]bool{true}, []string{"The private network speed of hardware server: 1234 was updated."})
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private-speed", "1000")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The private network speed of hardware server: 1234 was updated."))
			})
		})
	})
})
