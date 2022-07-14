package tags_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("tags cleanup", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cmd             *tags.CleanupCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		cmd = tags.NewCleanupCommand(fakeUI, fakeTagsManager)
		cliCommand = cli.Command{
			Name:        tags.TagsCleanupMetaData().Name,
			Description: tags.TagsCleanupMetaData().Description,
			Usage:       tags.TagsCleanupMetaData().Usage,
			Flags:       tags.TagsCleanupMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("tags cleanup", func() {

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTagsManager.GetUnattachedTagsReturns([]datatypes.Tag{}, errors.New("Failed to get Unattached Tags."))
			})
			It("Failed get Unattached Tags", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Unattached Tags."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerTags := []datatypes.Tag{
					datatypes.Tag{
						Name: sl.String("mytag")},
				}
				fakeTagsManager.GetUnattachedTagsReturns(fakerTags, nil)
				fakeTagsManager.DeleteTagReturns(false, errors.New("Failed to delete Tag"))
			})
			It("Failed get Unattached Tags", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Failed to delete Tag"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerTags := []datatypes.Tag{
					datatypes.Tag{
						Name: sl.String("mytag"),
					},
				}
				fakeTagsManager.GetUnattachedTagsReturns(fakerTags, nil)
			})
			It("Set command with --dry-run option", func() {
				err := testhelpers.RunCommand(cliCommand, "--dry-run")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("(Dry Run) Removing Tag"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerTags := []datatypes.Tag{
					datatypes.Tag{
						Name: sl.String("mytag"),
					},
				}
				fakeTagsManager.GetUnattachedTagsReturns(fakerTags, nil)
				fakeTagsManager.DeleteTagReturns(true, nil)
			})
			It("Remove tag", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Removing Tag"))
			})
		})

	})
})
