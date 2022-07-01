package tags_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("tags set", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cmd             *tags.SetCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		cmd = tags.NewSetCommand(fakeUI, fakeTagsManager)
		cliCommand = cli.Command{
			Name:        tags.TagsSetMetaData().Name,
			Description: tags.TagsSetMetaData().Description,
			Usage:       tags.TagsSetMetaData().Usage,
			Flags:       tags.TagsSetMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("tags set", func() {

		Context("Return error", func() {
			It("Set command without required options", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flags "tags, key-name, resource-id" not set`))
			})

			It("Set command without required --tags option", func() {
				err := testhelpers.RunCommand(cliCommand, "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flag "tags" not set`))
			})

			It("Set command without required --key-name option", func() {
				err := testhelpers.RunCommand(cliCommand, "--tags='tag1,tag2'", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flag "key-name" not set`))
			})

			It("Set command without required --resource-id option", func() {
				err := testhelpers.RunCommand(cliCommand, "--tags='tag1,tag2'", "--key-name=HARDWARE")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flag "resource-id" not set`))
			})

			It("Set invalid resource-id option", func() {
				err := testhelpers.RunCommand(cliCommand, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid value "abcde" for flag -resource-id:`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeTagsManager.SetTagsReturns(false, errors.New("Failed to set tags."))
			})
			It("Failed set tags", func() {
				err := testhelpers.RunCommand(cliCommand, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set tags."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeTagsManager.SetTagsReturns(true, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--tags='tag1,tag2'", "--key-name=HARDWARE", "--resource-id=123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Set tags successfully"))
			})
		})
	})
})
