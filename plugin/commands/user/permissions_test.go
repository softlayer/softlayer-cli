package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Edit Detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.PermissionsCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewPermissionsCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserPermissionsMetaData().Name,
			Description: user.UserPermissionsMetaData().Description,
			Usage:       user.UserPermissionsMetaData().Usage,
			Flags:       user.UserPermissionsMetaData().Flags,
			Action:      cmd.Run,
		}
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
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument."))
			})
		})

		Context("user permissions with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("User ID should be a number."))
			})
		})

		Context("user permissions error user", func() {
			It("return error", func() {
				fakeUserManager.GetUserReturns(datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user."))
			})
		})

		Context("user permissions error", func() {
			It("return error", func() {
				fakeUserManager.GetAllPermissionReturns([]datatypes.User_Customer_CustomerPermission_Permission{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get permissions."))
			})
		})

		Context("user permissions", func() {
			It("return user permissions", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID    Role Name   Description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123   role name   description"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Description    KeyName            Assigned"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permission 1   KEY_PERMISSION_1   true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Permission 2   KEY_PERMISSION_2   false"))
			})
		})

		Context("user permissions - master account", func() {
			It("return user permissions", func() {
				fakeUserManager.GetUserReturns(datatypes.User_Customer{
					IsMasterUserFlag: sl.Bool(true),
				}, nil)
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This account is the Master User and has all permissions enabled"))
			})
		})
	})
})
