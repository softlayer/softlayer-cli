package tags_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Tags Delete", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cliCommand      *tags.DeleteCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = tags.NewDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.TagsManager = fakeTagsManager
	})
	Describe("Tags Delete", func() {
		//sl tags delete
		Context("Tags Delete, success", func() {
			BeforeEach(func() {
				fakeTagsManager.DeleteTagReturns(true, nil)
			})
			It("Returns success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "testTag")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring("true"))
			})
		})
		Context("Tags Delete, JSON Output", func() {
			BeforeEach(func() {
				fakeTagsManager.DeleteTagReturns(true, nil)
			})
			It("Returns success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "testTag", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring("true"))
			})
		})
		Context("Tags Delete, Error", func() {
			BeforeEach(func() {
				fakeTagsManager.DeleteTagReturns(false, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("API Error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "testTag")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		Context("Tags Delete, No Arguments", func() {
			It("Incorrect Usage", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
			})
		})
	})
})
