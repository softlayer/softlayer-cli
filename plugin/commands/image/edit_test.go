package image_test

import (
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Image edit", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.EditCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewEditCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image edit", func() {
		Context("ISCSI cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Image edit with wrong image id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Image ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Image edit with correct image id and --name succeed", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{true}, []string{"The name of the image 1234 is updated."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--name", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The name of the image 1234 is updated."}))
			})
		})

		Context("Image edit with correct image id and --name fails", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{false}, []string{"Failed to update the image 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--name", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Failed to update the image 1234."}))
			})
		})

		Context("Image edit with correct image id and --note succeed", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{true}, []string{"The note of the image 1234 is updated."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The note of the image 1234 is updated."}))
			})
		})

		Context("Image edit with correct image id and --note fails", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{false}, []string{"Failed to update the image 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Failed to update the image 1234."}))
			})
		})

		Context("Image edit with correct image id and --tag succeed", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{true}, []string{"The tag of the image 1234 is updated."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--tag", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The tag of the image 1234 is updated."}))
			})
		})

		Context("Image edit with correct image id and --tag fails", func() {
			BeforeEach(func() {
				fakeImageManager.EditImageReturns([]bool{false}, []string{"Failed to update the image 1234."})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--tag", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Failed to update the image 1234."}))
			})
		})
	})
})
