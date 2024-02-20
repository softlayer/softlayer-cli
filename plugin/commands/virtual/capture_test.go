package virtual_test

import (
	"errors"
	"time"

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

var _ = Describe("VS capture", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CaptureCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCaptureCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("VS capture", func() {
		Context("VS capture without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS capture with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("VS capture without --name", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n|--name' is required"))
			})
		})

		Context("VS capture fails to get VS info", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage", "--all")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual server instance"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VS capture with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.CaptureImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage", "--all")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to capture image for virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VS capture ", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-30T00:00:00Z")
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					BlockDevices: []datatypes.Virtual_Guest_Block_Device{
						datatypes.Virtual_Guest_Block_Device{
							DiskImage: &datatypes.Virtual_Disk_Image{
								MetadataFlag: sl.Bool(true),
							},
						},
						datatypes.Virtual_Guest_Block_Device{
							DiskImage: &datatypes.Virtual_Disk_Image{
								Type: &datatypes.Virtual_Disk_Image_Type{
									KeyName: sl.String("SWAP"),
								},
							},
						},
						datatypes.Virtual_Guest_Block_Device{
							MountType: sl.String("CD"),
						},
						datatypes.Virtual_Guest_Block_Device{
							DiskImage: &datatypes.Virtual_Disk_Image{
								MetadataFlag: sl.Bool(false),
								Type: &datatypes.Virtual_Disk_Image_Type{
									KeyName: sl.String("SYSTEM"),
								},
							},
							Device: sl.String("0"),
						},
						datatypes.Virtual_Guest_Block_Device{
							DiskImage: &datatypes.Virtual_Disk_Image{
								MetadataFlag: sl.Bool(false),
								Type: &datatypes.Virtual_Disk_Image_Type{
									KeyName: sl.String("SYSTEM"),
								},
							},
							Device: sl.String("4"),
						},
					},
				}, nil)
				fakeVSManager.CaptureImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{
					Id:         sl.Int(12345678),
					CreateDate: sl.Time(created),
					Note:       sl.String("-"),
				}, nil)
			})
			It("--device option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage", "--device", "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-30T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
			})
			It("--all option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage", "--all")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-30T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
			})
			It("only system disk", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-30T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
			})
		})
	})
})
