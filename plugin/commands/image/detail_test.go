package image_test

import (
	"errors"
	"strings"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/image"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image detail", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeImageManager *testhelpers.FakeImageManager
		cmd              *image.DetailCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeImageManager = new(testhelpers.FakeImageManager)
		cmd = image.NewDetailCommand(fakeUI, fakeImageManager)
		cliCommand = cli.Command{
			Name:        metadata.ImageDetailMetaData().Name,
			Description: metadata.ImageDetailMetaData().Description,
			Usage:       metadata.ImageDetailMetaData().Usage,
			Flags:       metadata.ImageDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Image detail", func() {
		Context("ISCSI cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Image detail with wrong image id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Image ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Bad output format", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "12345", "--output", "FAKE")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format"))
			})
		})
		Context("JSON output format", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "12345", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("{}"))
			})
		})
		Context("Image detail with correct image id but server API call fails", func() {
			BeforeEach(func() {
				fakeImageManager.GetImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
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
					Id:               sl.Int(1234),
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
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"abcdefghijk"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"myimage"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Finished Import"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"278444"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Public"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SYSTEM"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"true"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"linux"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-29T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"200.00G"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Test_Transaction"}))
			})
			It("Test edge case output", func() {
				fakeImage.Status = nil
				fakeImage.PublicFlag = nil
				fakeImage.ImageType = nil
				fakeImageManager.GetImageReturns(fakeImage, nil)
				err := testhelpers.RunCommand(cliCommand, "1234")
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
