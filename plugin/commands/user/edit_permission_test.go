package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Edit Permission", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.EditPermissionCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewEditPermissionCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})

	Describe("user edit permission", func() {
		Context("user edit permission with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})

		Context("user edit permission with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
			})
		})

		Context("user edit permission just with id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: one of --permission and --from-user should be used to specify permissions"))
			})
		})

		Context("user edit permission fatal error", func() {
			It("return error", func() {
				fakeUserManager.AddPermissionReturns(false, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--permission", "PERMISSION")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update permissions\nInternal server error"))
			})
		})

		Context("user edit permission with correct id, permission and not send true o false", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--permission", "PERMISSION", "--enable", "notTrue o False")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("options for enable are true, false"))
			})
		})

		Context("user edit permission with correct id and permission by default enable", func() {
			It("updated permission", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--permission", "PERMISSION")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permissions updated successfully: PERMISSION"))
			})
		})

		Context("user edit permission to dissable with correct id and permission", func() {
			It("updated permission", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--permission", "PERMISSION", "--enable", "false")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permissions updated successfully: PERMISSION"))
			})
		})

		Context("user edit permission to same the another user", func() {
			It("updated permission", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "--from-user", "456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permissions updated successfully:"))
			})
		})
	})
})
