package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Edit Detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.EditCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewEditCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserEditMetaData().Name,
			Description: user.UserEditMetaData().Description,
			Usage:       user.UserEditMetaData().Usage,
			Flags:       user.UserEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("user edit detail", func() {
		Context("user edit detail with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})

		Context("user detail with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("User ID should be a number."))
			})
		})

		Context("user edit detail just with id like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--template' is required"))
			})
		})

		Context("user edit detail with empty template", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--template", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal template json: unexpected end of JSON input"))
			})
		})

		Context("user edit detail with error", func() {
			It("return error", func() {
				fakeUserManager.EditUserReturns(true, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "123", "--template", `{"firstName": "Test", "lastName": "Testerson"}`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update user 123."))
			})
		})

		Context("user edit detail with template", func() {
			It("return successful edit", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "--template", `{"firstName": "Test", "lastName": "Testerson"}`)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("User 123 updated successfully."))
			})
		})

		Context("user edit detail with template and output json", func() {
			It("return successful edit", func() {
				fakeUserManager.EditUserReturns(true, nil)
				err := testhelpers.RunCommand(cliCommand, "123", "--template", `{"firstName": "Test", "lastName": "Testerson"}`, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
			})
		})
	})
})
