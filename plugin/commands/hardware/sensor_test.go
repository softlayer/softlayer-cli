package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hadware sensor", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.SensorCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewSensorCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware sensor", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetSensorDataReturns([]datatypes.Container_RemoteManagement_SensorReading{}, errors.New("Failed to get hardware sensor data."))
			})
			It("Failed get hardware sensor data", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware sensor data."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerSensorData := []datatypes.Container_RemoteManagement_SensorReading{
					datatypes.Container_RemoteManagement_SensorReading{
						SensorUnits:      sl.String("degrees C"),
						SensorId:         sl.String("CPU1 Temperature"),
						Status:           sl.String("ok"),
						SensorReading:    sl.String("20.0"),
						LowerCritical:    sl.String("5.0"),
						LowerNonCritical: sl.String("10.0"),
						UpperCritical:    sl.String("85.0"),
						UpperNonCritical: sl.String("90.0"),
					},
					datatypes.Container_RemoteManagement_SensorReading{
						SensorUnits:      sl.String("Volts"),
						SensorId:         sl.String("12V"),
						Status:           sl.String("ok"),
						SensorReading:    sl.String("12.3"),
						LowerCritical:    sl.String("10.53"),
						LowerNonCritical: sl.String("10.78"),
						UpperCritical:    sl.String("12.91"),
						UpperNonCritical: sl.String("13.28"),
					},
					datatypes.Container_RemoteManagement_SensorReading{
						SensorUnits:      sl.String("Watts"),
						SensorId:         sl.String("CPU Power"),
						Status:           sl.String("ok"),
						SensorReading:    sl.String("26.0"),
						LowerCritical:    sl.String("20.0"),
						LowerNonCritical: sl.String("21.0"),
						UpperCritical:    sl.String("29.0"),
						UpperNonCritical: sl.String("30.0"),
					},
					datatypes.Container_RemoteManagement_SensorReading{
						SensorUnits:      sl.String("RPM"),
						SensorId:         sl.String("FAN1"),
						Status:           sl.String("ok"),
						SensorReading:    sl.String("8700.0"),
						LowerCritical:    sl.String("500.0"),
						LowerNonCritical: sl.String("700.0"),
						UpperCritical:    sl.String("25300.0"),
						UpperNonCritical: sl.String("25400.0"),
					},
					datatypes.Container_RemoteManagement_SensorReading{
						SensorUnits:   sl.String("discrete"),
						SensorId:      sl.String("CPU1 Temperature"),
						Status:        sl.String("ok"),
						SensorReading: sl.String("PS1 Status"),
					},
				}
				fakeHardwareManager.GetSensorDataReturns(fakerSensorData, nil)
			})
			It("display hardware sensor data", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--discrete")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("CPU1 Temperature"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ok"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("85.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("90.0"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("12V"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ok"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12.3"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.53"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.78"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12.91"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("13.28"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("CPU Power"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ok"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("26.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("21.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("29.0"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("30.0"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("FAN1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ok"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8700"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("500"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("700"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("25300"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("25400"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("CPU1 Temperature"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ok"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PS1 Status"))
			})
		})
	})
})
