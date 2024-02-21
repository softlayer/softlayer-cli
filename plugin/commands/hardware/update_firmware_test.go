package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware updatefirmware", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.UpdateFirmwareCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewUpdateFirmwareCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware update firmware", func() {
		Context("hardware update firmware without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("hardware update firmware with wrong ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware update firmware with correct ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will power off hardware server: 1234 and update device firmware. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("hardware update firmware with correct ID but server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.UpdateFirmwareReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("hardware update firmware with correct ID ", func() {
			BeforeEach(func() {
				fakeHardwareManager.UpdateFirmwareReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to update firmware for hardware server: 1234."))
			})
			It("return no error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Started to update firmware for hardware server: 1234."))
			})
		})

	})
	DescribeTable("Testing Flags",
		func(cliOptions []string, expected []bool) {
			// Adds in the Server ID and --force flag to the options
			cliOptions = append([]string{"1234", "-f"}, cliOptions...)
			err := testhelpers.RunCobraCommand(cliCommand.Command, cliOptions...)
			Expect(err).NotTo(HaveOccurred())
			hwid, ipmiFlag, raidFlag, biosFlag, hdFlag, nicFlag := fakeHardwareManager.UpdateFirmwareArgsForCall(0)
			Expect(hwid).To(Equal(1234))
			Expect(ipmiFlag).To(Equal(expected[0]))
			Expect(raidFlag).To(Equal(expected[1]))
			Expect(biosFlag).To(Equal(expected[2]))
			Expect(hdFlag).To(Equal(expected[3]))
			Expect(nicFlag).To(Equal(expected[4]))
		},
		Entry("IPMI Flag", []string{"--ipmi"}, []bool{true, false, false, false, false}),
		Entry("RAID Flag", []string{"--raid"}, []bool{false, true, false, false, false}),
		Entry("BIOS Flag", []string{"--bios"}, []bool{false, false, true, false, false}),
		Entry("HD Flag", []string{"--harddrive"}, []bool{false, false, false, true, false}),
		Entry("Network Flag", []string{"--network"}, []bool{false, false, false, false, true}),
		Entry("ALL Flags", []string{"--ipmi", "--raid", "--bios", "--harddrive", "--network"}, []bool{true, true, true, true, true}),
		Entry("No Flags (AKA All Flags", []string{}, []bool{true, true, true, true, true}),
	)
})
