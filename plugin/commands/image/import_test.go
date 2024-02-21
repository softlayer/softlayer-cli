package image_test

import (
	"errors"
	"time"

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

var _ = Describe("Image import", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.ImportCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewImportCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image import", func() {
		Context("Image import without three arguments", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "myimage")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires three arguments."))
			})
		})

		Context("Image import with inexistent URI", func() {
			BeforeEach(func() {
				fakeImageManager.ImportImageReturns(datatypes.Virtual_Guest_Block_Device_Template_Group{}, errors.New("SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "myimage", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)"))
			})
		})

		Context("Image import with correct data", func() {
			fakeImage := datatypes.Virtual_Guest_Block_Device_Template_Group{}
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-29T00:00:00Z")
				fakeImage = datatypes.Virtual_Guest_Block_Device_Template_Group{
					Id:               sl.Int(123456),
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
				fakeImageManager.ImportImageReturns(fakeImage, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "myimage", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myimage"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-29T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abcdefghijk"))
			})
		})

	})
})
