package image_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image export", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.ExportCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewExportCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image export", func() {
		Context("Image export without three arguments", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires three arguments."))
			})
		})

		Context("Image export with invalid ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Image ID"))
			})
		})

		Context("Image export with inexistent ID", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(false, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Image export with inexistent URI", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(false, errors.New("SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)"))
			})
		})

		Context("Image export with correct data", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The image 123456 was exported successfully!"))
			})
		})

	})
})
