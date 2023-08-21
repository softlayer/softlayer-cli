package image_test

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image detail", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.DetailCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image detail", func() {
		Context("ISCSI cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
		Context("Image detail with wrong image id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Image ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Bad output format", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--output", "FAKE")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})
		Context("JSON output format", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("{}"))
			})
		})
		Context("Image detail with correct image id but server API call fails", func() {
			BeforeEach(func() {
				fakeImageManager.GetImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get image: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image image with correct image id", func() {
			fakeImage := datatypes.Virtual_Guest_Block_Device_Template_Group{}
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				fakeImage = datatypes.Virtual_Guest_Block_Device_Template_Group{
					Id: sl.Int(1234),
					AccountReferences: []datatypes.Virtual_Guest_Block_Device_Template_Group_Accounts{
						{
							AccountId:  sl.Int(654),
							CreateDate: sl.Time(created),
						},
					},
					GlobalIdentifier: sl.String("abcdefghijk"),
					Name:             sl.String("myimage"),
					Status: &datatypes.Virtual_Guest_Block_Device_Template_Group_Status{
						Name: sl.String("Finished Import"),
					},
					AccountId:  sl.Int(278444),
					PublicFlag: sl.Int(1),
					ImageType: &datatypes.Virtual_Disk_Image_Type{
						Name: sl.String("SYSTEM"),
					},
					FlexImageFlag: sl.Bool(true),
					Note:          sl.String("linux"),
					CreateDate:    sl.Time(created),
					Children: []datatypes.Virtual_Guest_Block_Device_Template_Group{
						datatypes.Virtual_Guest_Block_Device_Template_Group{
							BlockDevicesDiskSpaceTotal: sl.Float(107374182400),
							Datacenter: &datatypes.Location{
								Name: sl.String("tok02"),
							},
							Transaction: &datatypes.Provisioning_Version1_Transaction{
								TransactionStatus: &datatypes.Provisioning_Version1_Transaction_Status{
									Name: sl.String("Test_Transaction"),
								},
							},
							BlockDevices: []datatypes.Virtual_Guest_Block_Device_Template{
								datatypes.Virtual_Guest_Block_Device_Template{
									DiskImage: &datatypes.Virtual_Disk_Image{
										SoftwareReferences: []datatypes.Virtual_Disk_Image_Software{
											datatypes.Virtual_Disk_Image_Software{
												SoftwareDescription: &datatypes.Software_Description{
													LongDescription: sl.String("Ubuntu 20.04-64 Minimal for VSI"),
												},
											},
										},
									},
								},
							},
						},
						datatypes.Virtual_Guest_Block_Device_Template_Group{
							BlockDevicesDiskSpaceTotal: sl.Float(107374182400),
							Datacenter: &datatypes.Location{
								Name: sl.String("dal10"),
							},
						},
					},
				}
			})
			It("return no error", func() {
				fakeImageManager.GetImageReturns(fakeImage, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abcdefghijk"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myimage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Finished Import"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("278444"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Public"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SYSTEM"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("linux"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Ubuntu 20.04-64 Minimal for VSI"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("share image"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("654"))
			})
			It("Test edge case output", func() {
				fakeImage.Status = nil
				fakeImage.PublicFlag = nil
				fakeImage.ImageType = nil
				fakeImageManager.GetImageReturns(fakeImage, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				// Removes whitespace from the string for easier testing.
				output := strings.ReplaceAll(fakeUI.Outputs(), " ", "")

				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(ContainSubstring("1234"))
				Expect(output).To(ContainSubstring("status-"))
				Expect(output).To(ContainSubstring("visibilityPrivate"))
				Expect(output).To(ContainSubstring("type-"))

			})
		})
	})
})
