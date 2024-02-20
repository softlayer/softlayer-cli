package tags_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("tags set", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cliCommand      *tags.SetCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = tags.NewSetCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.TagsManager = fakeTagsManager
	})

	Describe("tags set", func() {

		Context("Return error", func() {

			It("Set command without required --tags option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--tags' is required"))
			})

			It("Set command without required --key-name option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--tags='tag1,tag2'", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--key-name' is required"))
			})

			It("Set command without required --resource-id option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--tags='tag1,tag2'", "--key-name=HARDWARE")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--resource-id' is required"))
			})

			It("Set invalid resource-id option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid argument "abcde" for "--resource-id" flag`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTagsManager.SetTagsReturns(false, errors.New("Failed to set tags."))
			})
			It("Failed set tags", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set tags."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeTagsManager.SetTagsReturns(true, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Set tags successfully"))
			})
		})
	})
})
