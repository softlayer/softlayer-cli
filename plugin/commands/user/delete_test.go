package user_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Delete", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.DeleteCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewDeleteCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserDeleteMataData().Name,
			Description: user.UserDeleteMataData().Description,
			Usage:       user.UserDeleteMataData().Usage,
			Flags:       user.UserDeleteMataData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("user delete", func() {
		Context("user delete with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})

		Context("user delete with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: User ID should be a number.")).To(BeTrue())
			})
		})

		Context("user delete with id", func() {
			It("return aborted", func() {
				fakeUI.Inputs("")
				id_user := "123"
				err := testhelpers.RunCommand(cliCommand, id_user)
				Expect(err).NotTo(HaveOccurred())
				response := fmt.Sprintf("This will delete the user: %s and cannot be undone. Continue?", id_user)
				Expect(fakeUI.Outputs()).To(ContainSubstring(response))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("user delete confirmation fail", func() {
			It("return error", func() {
				fakeUI.Inputs("123456")
				id_user := "123"
				err := testhelpers.RunCommand(cliCommand, id_user)
				Expect(err).To(HaveOccurred())
				response := fmt.Sprintf("This will delete the user: %s and cannot be undone. Continue?", id_user)
				Expect(fakeUI.Outputs()).To(ContainSubstring(response))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("user delete error", func() {
			It("return error", func() {
				fakeUserManager.EditUserReturns(false, errors.New("Internal server error"))
				fakeUI.Inputs("y")
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete user."))
			})
		})

		Context("user delete with id", func() {
			It("return Ok", func() {
				fakeUI.Inputs("y")
				id_user := "123"
				err := testhelpers.RunCommand(cliCommand, id_user)
				Expect(err).NotTo(HaveOccurred())
				response := fmt.Sprintf("This will delete the user: %s and cannot be undone. Continue?", id_user)
				Expect(fakeUI.Outputs()).To(ContainSubstring(response))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

		Context("user delete with id and flag -f force", func() {
			It("return Ok", func() {
				fakeUI.Inputs("y")
				id_user := "123"
				err := testhelpers.RunCommand(cliCommand, id_user, "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})

})
