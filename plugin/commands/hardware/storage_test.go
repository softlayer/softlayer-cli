package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware storage", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.StorageCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewStorageCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        hardware.HardwareStorageMetaData().Name,
			Description: hardware.HardwareStorageMetaData().Description,
			Usage:       hardware.HardwareStorageMetaData().Usage,
			Flags:       hardware.HardwareStorageMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware storage", func() {
		Context("hardware storage without id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})

		Context("hardware storage with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware storage ISCSI with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetStorageDetailsReturns([]datatypes.Network_Storage{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get iscsi storage detail for the hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware storage credentials with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetStorageCredentialsReturns(datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the storage credential detail for the hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware hard drives with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardDrivesReturns([]datatypes.Hardware_Component{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the hard drives detail for the hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware iscsi", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetStorageDetailsReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id:                              sl.Int(123458),
						Username:                        sl.String("SL02SEL1234567-20"),
						CapacityGb:                      sl.Int(16000),
						ServiceResourceBackendIpAddress: sl.String("10.10.10.10"),
						AllowedHardware: []datatypes.Hardware{
							datatypes.Hardware{
								Datacenter: &datatypes.Location{
									LongName: sl.String("Dallas 10"),
								},
							},
						},
						Notes: sl.String("Test notes"),
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("SL02SEL1234567-20"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("16000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.10.10.10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Test notes"))
			})
		})

		Context("hardware Credentials", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetStorageCredentialsReturns(datatypes.Network_Storage_Allowed_Host{
					Credential: &datatypes.Network_Storage_Credential{
						Username: sl.String("SL02SU1234567-H1"),
						Password: sl.String("fMnY59hjkhkj"),
					},
					Name: sl.String("iqn.2021-07.com.ibm:SL02SU1234567-h1"),
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("SL02SU1234567-H1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("fMnY59hjkhkj"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("iqn.2021-07.com.ibm:SL02SU1234567-h1"))
			})
		})

		Context("hardware hard drives", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardDrivesReturns([]datatypes.Hardware_Component{
					datatypes.Hardware_Component{
						HardwareComponentModel: &datatypes.Hardware_Component_Model{
							HardwareGenericComponentModel: &datatypes.Hardware_Component_Model_Generic{
								HardwareComponentType: &datatypes.Hardware_Component_Type{
									Type: sl.String("Hard Drive"),
								},
								Capacity: sl.Float(2000.00),
								Units:    sl.String("GB"),
							},
							Manufacturer: sl.String("Seagate"),
							Name:         sl.String("Enterprise Capacity 3.5"),
						},
						SerialNumber: sl.String("zc2fdsfsdf"),
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hard Drive"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Seagate"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Enterprise Capacity 3.5"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("zc2fdsfsdf"))
			})
		})
	})
})
