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

var _ = Describe("User List", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.ListCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager

		testListUser := []datatypes.User_Customer{
			datatypes.User_Customer{
				Id:                        sl.Int(5555),
				Username:                  sl.String("ATestUser"),
				Email:                     sl.String("user2@email.com"),
				DisplayName:               sl.String("DisplayedName"),
				ExternalBindingCount:      sl.Uint(123),
				ApiAuthenticationKeyCount: sl.Uint(123456),
				UserStatus: &datatypes.User_Customer_Status{
					Name: sl.String("ACTIVE"),
				},
			},
			datatypes.User_Customer{
				Id:                        sl.Int(5556),
				Username:                  sl.String("ATestUser2"),
				Email:                     sl.String("user2@email.com"),
				DisplayName:               sl.String("DisplayedName2"),
				ExternalBindingCount:      sl.Uint(1234),
				ApiAuthenticationKeyCount: sl.Uint(1234567),
			},
		}
		fakeUserManager.ListUsersReturns(testListUser, nil)
	})

	Describe("user list ", func() {
		Context("user list with unknown column", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "noExist")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column noExist is not supported."))
			})
		})

		Context("user list fatal error", func() {
			It("return error", func() {
				fakeUserManager.ListUsersReturns([]datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list users."))
			})
		})

		Context("user list", func() {
			It("return users list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("id     username     email             displayName      2FA   classicAPIKey"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5555   ATestUser    user2@email.com   DisplayedName    yes   yes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5556   ATestUser2   user2@email.com   DisplayedName2   yes   yes"))
			})
		})

		Context("user list with column", func() {
			It("return users list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "username")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("username"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ATestUser2"))
			})
		})

		Context("user list in format json", func() {
			It("return users list", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"apiAuthenticationKeyCount": 123456,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"displayName": "DisplayedName",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "user2@email.com",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"externalBindingCount": 123,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 5555,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"userStatus": {`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "ACTIVE"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`},`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"username": "ATestUser"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"apiAuthenticationKeyCount": 1234567,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"displayName": "DisplayedName2",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "user2@email.com",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"externalBindingCount": 1234,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 5556,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"username": "ATestUser2"`))
			})
		})
	})
})
