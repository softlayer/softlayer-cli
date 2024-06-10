package user_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("User Permissions", func() {
	var (
		fakeUI          *terminal.FakeUI
		cliCommand      *user.PermissionsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
		fakeHandler *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewPermissionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")

	})

	Describe("user permissions ", func() {
		Context("user permissions", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("user permissions with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("User ID should be a number."))
			})
		})

		Context("user permissions error user", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_User_Customer", "getObject",
										500, "Internal Server Error")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user."))
			})
		})

		Context("user permissions error", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_User_Permission_Department", "getAllObjects",
										500, "Internal Server Error")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get permissions."))
			})
		})

		Context("user permissions", func() {
			It("return user permissions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID   Role Name   Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2    role name   description of the role"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ADMINISTRATIVE   KeyName                           Assigned   Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACCOUNT_BRAND_ADD                 false      Permission to create sub brands"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACCOUNT_BILLING_SYSTEM            true       Permission to access account billing system type determination"))
			})
		})

		Context("user permissions - master account", func() {
			It("return user permissions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12345")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This account is the Master User and has all permissions enabled"))
			})
		})
	})
})
