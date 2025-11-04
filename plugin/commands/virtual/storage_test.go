package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("virtual storage", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.StorageCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
		fakeHandler   *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewStorageCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		
	})

	AfterEach(func() {
		// Clear API call logs and any errors that might have been set after every test
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("virtual storage", func() {
		BeforeEach(func() {
			cliCommand.VirtualServerManager = fakeVSManager
		})
		Context("User Input Checks", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("API Errors", func() {
			It("Failed to get ISCSI Storage", func() {
				fakeVSManager.GetStorageDetailsReturns([]datatypes.Network_Storage{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get iscsi storage detail for the virtual server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("virtual storage credentials with server fails", func() {
				fakeVSManager.GetStorageCredentialsReturns(datatypes.Network_Storage_Allowed_Host{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the storage credential detail for the virtual server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("virtual portable storage with server fails", func() {
				fakeVSManager.GetPortableStorageReturns([]datatypes.Virtual_Disk_Image{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the portable storage detail for the virtual server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("virtual local disks with server fails", func() {
				fakeVSManager.GetLocalDisksReturns([]datatypes.Virtual_Guest_Block_Device{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the local disks detail for the virtual server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})

		})


		Context("hardware iscsi", func() {
			BeforeEach(func() {
				fakeVSManager.GetStorageDetailsReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id:                              sl.Int(123458),
						Username:                        sl.String("SL02SEL1234567-20"),
						CapacityGb:                      sl.Int(16000),
						ServiceResourceBackendIpAddress: sl.String("10.10.10.10"),
						AllowedVirtualGuests: []datatypes.Virtual_Guest{
							datatypes.Virtual_Guest{
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
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
				fakeVSManager.GetStorageCredentialsReturns(datatypes.Network_Storage_Allowed_Host{
					Credential: &datatypes.Network_Storage_Credential{
						Username: sl.String("SL02SU1234567-H1"),
						Password: sl.String("fMnY59hjkhkj"),
					},
					Name: sl.String("iqn.2021-07.com.ibm:SL02SU1234567-h1"),
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("SL02SU1234567-H1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("fMnY59hjkhkj"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("iqn.2021-07.com.ibm:SL02SU1234567-h1"))
			})
		})

		Context("hardware portable storage", func() {
			BeforeEach(func() {
				fakeVSManager.GetPortableStorageReturns([]datatypes.Virtual_Disk_Image{
					datatypes.Virtual_Disk_Image{
						Description: sl.String("Test Description"),
						Capacity:    sl.Int(16000),
						BillingItem: &datatypes.Billing_Item_Virtual_Disk_Image{
							Billing_Item: datatypes.Billing_Item{
								Location: &datatypes.Location{
									LongName: sl.String("Dallas 10"),
								},
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Test Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("16000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 10"))
			})
		})

		Context("hardware local disks", func() {
			BeforeEach(func() {
				fakeVSManager.GetLocalDisksReturns([]datatypes.Virtual_Guest_Block_Device{
					datatypes.Virtual_Guest_Block_Device{
						MountType: sl.String("Disk"),
						Device:    sl.String("1"),
						DiskImage: &datatypes.Virtual_Disk_Image{
							Capacity:    sl.Int(100),
							Units:       sl.String("GB"),
							Description: sl.String("Tes description SWAP"),
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Disk"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("100"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("GB"))
			})
		})
	})
	Describe("virtual storage with fixtures", func() {
		Context("Issues943", func() {
			It("Successful", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "934")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Disk"))
			})
		})
	})

})
