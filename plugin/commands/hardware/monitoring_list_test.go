package hardware_test

import (
	"errors"

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

var _ = Describe("hardware monitoring list", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.MonitoringListCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewMonitoringListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("Hardware monitoring list", func() {
		Context("Return error", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})

			It("Set command with an invalid output format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Internal Server Error"))
			})
			It("Command fails to get hardware", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerHardware := datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{
						Domain:                  sl.String("domain.com"),
						PrimaryIpAddress:        sl.String("9.9.9.9"),
						PrimaryBackendIpAddress: sl.String("1.1.1.1"),
						Datacenter: &datatypes.Location{
							LongName: sl.String("Dallas 10"),
						},
						NetworkMonitors: []datatypes.Network_Monitor_Version1_Query_Host{
							datatypes.Network_Monitor_Version1_Query_Host{
								Id:        sl.Int(678),
								IpAddress: sl.String("2.2.2.2"),
								Status:    sl.String("ON"),
								QueryType: &datatypes.Network_Monitor_Version1_Query_Type{
									Name: sl.String("SERVICE PING"),
								},
								ResponseAction: &datatypes.Network_Monitor_Version1_Query_ResponseType{
									ActionDescription: sl.String("Do Nothing"),
								},
							},
						},
					},
				}
				fakeHardwareManager.GetHardwareReturns(fakerHardware, nil)
			})
			It("Set command with correct hardware id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.9.9.9"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ON"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SERVICE PING"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Do Nothing"))
			})
		})
	})
})
