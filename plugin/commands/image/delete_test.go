package image_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/image"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image delete", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeImageManager *testhelpers.FakeImageManager
		cmd              *image.DeleteCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeImageManager = new(testhelpers.FakeImageManager)
		cmd = image.NewDeleteCommand(fakeUI, fakeImageManager)
		cliCommand = cli.Command{
			Name:        metadata.ImageDelMetaData().Name,
			Description: metadata.ImageDelMetaData().Description,
			Usage:       metadata.ImageDelMetaData().Usage,
			Flags:       metadata.ImageDelMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Image delete", func() {
		Context("Image delete without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Image delete with wrong image id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Image ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Image delete with correct image id but image not found", func() {
			BeforeEach(func() {
				fakeImageManager.DeleteImageReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Unable to find image with ID 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound")).To(BeTrue())
			})
		})

		Context("Image delete with correct image id but server API call fails", func() {
			BeforeEach(func() {
				fakeImageManager.DeleteImageReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to delete image: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Image delete with correct image id", func() {
			BeforeEach(func() {
				fakeImageManager.DeleteImageReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Image 1234 was deleted."}))
			})
		})
	})
})
