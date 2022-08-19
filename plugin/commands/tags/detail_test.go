package tags_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Tags Detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cliCommand      *tags.DetailCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = tags.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.TagsManager = fakeTagsManager
		fakeTagsManager.GetTagByTagNameReturns(FakeTags, nil)
		fakeTagsManager.GetTagReferencesReturns(FakeTagReference, nil)
		fakeTagsManager.ReferenceLookupReturns("Hardware.Name")
	})
	Describe("Tags Detail", func() {
		//sl tags list
		Context("Tags Detail, success", func() {
			It("Returns success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "testTag")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(4))
				Expect(results[2]).To(ContainSubstring("22222   HARDWARE   Hardware.Name"))
			})
		})
		//sl tags list --output=JSON
		Context("Tags Detail, success JSON", func() {
			It("Returns JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "test1", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring(`"ResourceName": "Hardware.Name"`))
			})
		})
		//sl tags list (no name)
		Context("Tags Detail, No Arguments", func() {
			It("Incorrect Usage", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
			})
		})
		//sl tags list (SLAPI ERROR)
		Context("Tags Detail, SLAPI Error", func() {
			BeforeEach(func() {
				fakeTagsManager.GetTagByTagNameReturns(nil, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("SLAPI error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "testTag")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
	})
})
