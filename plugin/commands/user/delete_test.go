package user_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Delete", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.DeleteCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewDeleteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})
	Describe("user delete", func() {
		Context("user delete with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage : This command requires one argument")).To(BeTrue())
			})
		})

		Context("user delete with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: User ID should be a number.")).To(BeTrue())
			})
		})

		Context("user delete with id", func() {
			It("return aborted", func() {
				fakeUI.Inputs("")
				id_user := "123"
				err := testhelpers.RunCobraCommand(cliCommand.Command, id_user)
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, id_user)
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete user."))
			})
		})

		Context("user delete with id", func() {
			It("return Ok", func() {
				fakeUI.Inputs("y")
				id_user := "123"
				err := testhelpers.RunCobraCommand(cliCommand.Command, id_user)
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, id_user, "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})

})
