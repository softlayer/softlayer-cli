package hardware_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Hardware Server Guests", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.GuestsCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewGuestsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("Hardware Server Guests", func() {
		Context("Guests without HW ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Guests with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("Hardware guests with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareGuestsReturns([]datatypes.Virtual_Guest{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the guests instances for the hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware guests", func() {
			created, _ := time.Parse(time.RFC3339, "2021-08-30T00:00:00Z")
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareGuestsReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:          sl.Int(1234),
						Hostname:    sl.String("TestHostname"),
						MaxCpu:      sl.Int(1),
						MaxCpuUnits: sl.String("CORE"),
						MaxMemory:   sl.Int(8192),
						CreateDate:  sl.Time(created),
						Status: &datatypes.Virtual_Guest_Status{
							KeyName: sl.String("ACTIVE"),
						},
						PowerState: &datatypes.Virtual_Guest_Power_State{
							KeyName: sl.String("RUNNING"),
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("TestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CORE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8192"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-30T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RUNNING"))
			})
		})
	})
})
