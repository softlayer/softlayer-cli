package tags_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Tags Delete", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cmd             *tags.DeleteCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		cmd = tags.NewDeleteCommand(fakeUI, fakeTagsManager)
		cliCommand = cli.Command{
			Name:        metadata.TagsDeleteMetaData().Name,
			Description: metadata.TagsDeleteMetaData().Description,
			Usage:       metadata.TagsDeleteMetaData().Usage,
			Flags:       metadata.TagsDeleteMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("Tags Delete", func() {
		//sl tags delete
		Context("Tags Delete, success", func() {
			BeforeEach(func() {
				fakeTagsManager.DeleteTagReturns(true, nil)
			})
			It("Returns success", func() {
				err := testhelpers.RunCommand(cliCommand, "testTag")
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
				err := testhelpers.RunCommand(cliCommand, "testTag", "--output", "JSON")
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
				err := testhelpers.RunCommand(cliCommand, "testTag")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		Context("Tags Delete, No Arguments", func() {
			It("Incorrect Usage", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
			})
		})
	})
})
