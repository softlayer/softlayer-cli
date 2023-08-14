package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("User Permissions", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.PermissionsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewPermissionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
		testUser := datatypes.User_Customer{
			Roles: []datatypes.User_Permission_Role{
				datatypes.User_Permission_Role{
					Id:          sl.Int(123),
					Name:        sl.String("role name"),
					Description: sl.String("description"),
				},
			},
			Permissions: []datatypes.User_Customer_CustomerPermission_Permission{
				datatypes.User_Customer_CustomerPermission_Permission{
					KeyName: sl.String("KEY_PERMISSION_1"),
					Name:    sl.String("Permission 1"),
				},
			},
		}
		testAllPermissions := []datatypes.User_Customer_CustomerPermission_Permission{
			datatypes.User_Customer_CustomerPermission_Permission{
				KeyName: sl.String("KEY_PERMISSION_1"),
				Name:    sl.String("Permission 1"),
			},
			datatypes.User_Customer_CustomerPermission_Permission{
				KeyName: sl.String("KEY_PERMISSION_2"),
				Name:    sl.String("Permission 2"),
			},
		}
		fakeUserManager.GetUserReturns(testUser, nil)
		fakeUserManager.GetAllPermissionReturns(testAllPermissions, nil)
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
			It("return error", func() {
				fakeUserManager.GetUserReturns(datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user."))
			})
		})

		Context("user permissions error", func() {
			It("return error", func() {
				fakeUserManager.GetAllPermissionReturns([]datatypes.User_Customer_CustomerPermission_Permission{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get permissions."))
			})
		})

		Context("user permissions", func() {
			It("return user permissions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID    Role Name   Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123   role name   description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Description    KeyName            Assigned"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permission 1   KEY_PERMISSION_1   true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permission 2   KEY_PERMISSION_2   false"))
			})
		})

		Context("hide user permissions", func() {
			It("return not equal user permissions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(Equal("ACCOUNT_SUMMARY_VIEW"))
				Expect(fakeUI.Outputs()).To(Not(Equal("REQUEST_COMPLIANCE_REPORT")))
				Expect(fakeUI.Outputs()).To(Not(Equal("COMPANY_EDIT")))
				Expect(fakeUI.Outputs()).To(Not(Equal("ONE_TIME_PAYMENTS")))
				Expect(fakeUI.Outputs()).To(Not(Equal("UPDATE_PAYMENT_DETAILS")))
				Expect(fakeUI.Outputs()).To(Not(Equal("EU_LIMITED_PROCESSING_MANAGE")))
				Expect(fakeUI.Outputs()).To(Not(Equal("TICKET_ADD")))
				Expect(fakeUI.Outputs()).To(Not(Equal("TICKET_EDIT")))
				Expect(fakeUI.Outputs()).To(Not(Equal("TICKET_SEARCH")))
				Expect(fakeUI.Outputs()).To(Not(Equal("TICKET_VIEW")))
				Expect(fakeUI.Outputs()).To(Not(Equal("TICKET_VIEW_ALL")))
			})
		})

		Context("user permissions - master account", func() {
			It("return user permissions", func() {
				fakeUserManager.GetUserReturns(datatypes.User_Customer{
					IsMasterUserFlag: sl.Bool(true),
				}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This account is the Master User and has all permissions enabled"))
			})
		})
	})
})
