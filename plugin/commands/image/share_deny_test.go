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

var _ = Describe("Image share", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.ShareDenyCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewShareDenyCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("Image share deny", func() {
		Context("return error", func() {
			It("image share deny without imageId", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two argument"))
			})

			It("set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Image Id'. It must be a positive integer."))
			})

			It("set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Account Id'. It must be a positive integer."))
			})

			It("image share without account-id", func() {
				fakeImageManager.ShareDenyImageReturns(false, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "654321")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to deny share image"))
			})
		})

		Context("return no error", func() {
			It("correct use", func() {
				fakeImageManager.ShareDenyImageReturns(true, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "654321")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Image"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("was deny shared with account"))
			})
		})
	})
})
