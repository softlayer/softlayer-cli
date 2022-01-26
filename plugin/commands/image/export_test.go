package image_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image export", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeImageManager *testhelpers.FakeImageManager
		cmd              *image.ExportCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeImageManager = new(testhelpers.FakeImageManager)
		cmd = image.NewExportCommand(fakeUI, fakeImageManager)
		cliCommand = cli.Command{
			Name:        metadata.ImageDelMetaData().Name,
			Description: metadata.ImageDelMetaData().Description,
			Usage:       metadata.ImageDelMetaData().Usage,
			Flags:       metadata.ImageDelMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Image export", func() {
		Context("Image export without three arguments", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires three arguments.")).To(BeTrue())
			})
		})

		Context("Image export with invalid ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Image ID")).To(BeTrue())
			})
		})

		Context("Image export with inexistent ID", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(false, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)")).To(BeTrue())
			})
		})

		Context("Image export with inexistent URI", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(false, errors.New("SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)"))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_Public: Template configuration uri specified an invalid network storage service resource protocol. (HTTP 500)")).To(BeTrue())
			})
		})

		Context("Image export with correct data", func() {
			BeforeEach(func() {
				fakeImageManager.ExportImageReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "swift://SLOS123456-10@dal05/OS/testImage4f.iso", "PI-ABCDE-abcde1234567890abcdefgrty1234567890")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The image 123456 was exported successfully!"}))
			})
		})

	})
})
