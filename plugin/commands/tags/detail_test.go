package tags_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Tags Detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cmd             *tags.DetailCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		cmd = tags.NewDetailCommand(fakeUI, fakeTagsManager)
		cliCommand = cli.Command{
			Name:        metadata.TagsDetailsMetaData().Name,
			Description: metadata.TagsDetailsMetaData().Description,
			Usage:       metadata.TagsDetailsMetaData().Usage,
			Flags:       metadata.TagsDetailsMetaData().Flags,
			Action:      cmd.Run,
		}
		fakeTagsManager.GetTagByTagNameReturns(FakeTags, nil)
		fakeTagsManager.GetTagReferencesReturns(FakeTagReference, nil)
		fakeTagsManager.ReferenceLookupReturns("Hardware.Name")
	})
	Describe("Tags Detail", func() {
		//sl tags list
		Context("Tags Detail, success", func() {
			It("Returns success", func() {
				err := testhelpers.RunCommand(cliCommand, "testTag")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(4))
				Expect(results[2]).To(ContainSubstring("22222   HARDWARE   Hardware.Name"))
			})
		})
		//sl tags list --output=JSON
		Context("Tags Detail, success JSON", func() {
			It("Returns JSON", func() {
				err := testhelpers.RunCommand(cliCommand, "test1", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring(`"ResourceName": "Hardware.Name"`))
			})
		})
		//sl tags list (no name)
		Context("Tags Detail, No Arguments", func() {
			It("Incorrect Usage", func() {
				err := testhelpers.RunCommand(cliCommand)
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
				err := testhelpers.RunCommand(cliCommand, "testTag")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
	})
})
