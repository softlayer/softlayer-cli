package image_test

import (
	"errors"
	"strings"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image list", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.ListCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image list", func() {
		Context("Image list with both --public and --private", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--public", "--private")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--public] is not allowed with [--private].")).To(BeTrue())
			})
		})

		Context("Image list with --private only but server API call fails", func() {
			BeforeEach(func() {
				fakeImageManager.ListPrivateImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--private")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list private images.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image list with --private only", func() {
			BeforeEach(func() {
				fakeImageManager.ListPrivateImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{
					datatypes.Virtual_Guest_Block_Device_Template_Group{
						Id:   sl.Int(1234),
						Name: sl.String("image-1234"),
						ImageType: &datatypes.Virtual_Disk_Image_Type{
							Name: sl.String("SYSTEM"),
						},
						AccountId: sl.Int(278444),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--private")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"image-1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Private"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"Public"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SYSTEM"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"278444"}))
			})
		})

		Context("Image list with --public only but server API call fails", func() {
			BeforeEach(func() {
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--public")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list public images.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image list with --public only", func() {
			BeforeEach(func() {
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{
					datatypes.Virtual_Guest_Block_Device_Template_Group{
						Id:   sl.Int(1234),
						Name: sl.String("image-1234"),
						ImageType: &datatypes.Virtual_Disk_Image_Type{
							Name: sl.String("SYSTEM"),
						},
						AccountId: sl.Int(278444),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--public")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"image-1234"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"Private"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Public"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SYSTEM"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"278444"}))
			})
		})

		Context("Image list without --public or --private but list public image call fails", func() {
			BeforeEach(func() {
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list public images.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image list without --public or --private but list private image call fails", func() {
			BeforeEach(func() {
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, nil)
				fakeImageManager.ListPrivateImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))

			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list private images.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image list without --public or --private and both call succeed", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{
					datatypes.Virtual_Guest_Block_Device_Template_Group{
						Id:   sl.Int(1234),
						Name: sl.String("image-1234"),
						ImageType: &datatypes.Virtual_Disk_Image_Type{
							Name: sl.String("SYSTEM"),
						},
						AccountId: sl.Int(278444),
					},
				}, nil)
				fakeImageManager.ListPrivateImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{
					datatypes.Virtual_Guest_Block_Device_Template_Group{
						Id:   sl.Int(5678),
						Name: sl.String("image-5678"),
						ImageType: &datatypes.Virtual_Disk_Image_Type{
							Name: sl.String("SYSTEM"),
						},
						CreateDate: sl.Time(created),
						Note:       sl.String("myNote"),
						AccountId:  sl.Int(278444),
						Children: []datatypes.Virtual_Guest_Block_Device_Template_Group{
							datatypes.Virtual_Guest_Block_Device_Template_Group{
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
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "1234")).To(BeTrue())
				Expect(strings.Contains(results[1], "image-1234")).To(BeTrue())
				Expect(strings.Contains(results[1], "Public")).To(BeTrue())
				Expect(strings.Contains(results[1], "SYSTEM")).To(BeTrue())
				Expect(strings.Contains(results[1], "278444")).To(BeTrue())
				Expect(strings.Contains(results[1], "-")).To(BeTrue())

				Expect(strings.Contains(results[2], "5678")).To(BeTrue())
				Expect(strings.Contains(results[2], "image-5678")).To(BeTrue())
				Expect(strings.Contains(results[2], "Private")).To(BeTrue())
				Expect(strings.Contains(results[2], "SYSTEM")).To(BeTrue())
				Expect(strings.Contains(results[2], "Ubuntu 20.04-64 Minimal for VSI")).To(BeTrue())
				Expect(strings.Contains(results[2], "2016-12-29T00:00:00Z")).To(BeTrue())
				Expect(strings.Contains(results[2], "myNote")).To(BeTrue())
			})
		})

		Context("Image list without --public or --private and both call succeed without any image", func() {
			BeforeEach(func() {
				fakeImageManager.ListPublicImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, nil)
				fakeImageManager.ListPrivateImagesReturns([]datatypes.Virtual_Guest_Block_Device_Template_Group{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No image found."))
			})
		})
	})
})
