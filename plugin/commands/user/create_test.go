package user_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
)

var _ = Describe("Create", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.CreateCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager

		testUser := datatypes.User_Customer{
			Id:       sl.Int(6666),
			Username: sl.String("createdUser"),
			Email:    sl.String("createdUser@email.com"),
		}
		fakeUserManager.CreateUserReturns(testUser, nil)
	})
	Describe("user create", func() {
		Context("user create with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage : This command requires one argument")).To(BeTrue())
			})
		})

		Context("create user with fail confirmation", func() {
			It("return error", func() {
				fakeUI.Inputs("123456")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("create user with No confirmation", func() {
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})

		Context("User Create error", func() {
			It("return error", func() {
				fakeUserManager.CreateUserReturns(datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add user."))
			})
		})

		Context("Basic User Create usage", func() {
			It("Create a user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
		})

		Context("User Create", func() {
			It("Create a user", func() {
				fakeUI.Inputs("Y")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
		})

		Context("User Create from user", func() {
			It("Create a user", func() {
				fakeUI.Inputs("Y")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--from-user", "456", "--password", "MyPassWord", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
		})

		Context("User Create from wrong template", func() {
			It("Create a user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f", "--template", ``)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal template json: unexpected end of JSON input"))
			})
		})

		Context("User Create from template", func() {
			It("Create a user", func() {
				testUser := datatypes.User_Customer{
					Id:        sl.Int(6666),
					Username:  sl.String("createdUser"),
					Email:     sl.String("createdUser@email.com"),
					FirstName: sl.String("Test"),
					LastName:  sl.String("Testerson"),
				}
				fakeUserManager.CreateUserReturns(testUser, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f", "--template", `{"firstName": "Test", "lastName": "Testerson"}`)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))

			})
		})

		Context("User Create with generated password", func() {
			It("Create a user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "generate", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password"))
			})
		})

	})

	Describe("structAssignment", func() {

		A1 := "11"
		A2 := "12"
		B1 := "21"
		B2 := "22"
		var S1, S2 datatypes.User_Customer
		Context("structAssignment", func() {
			BeforeEach(func() {
				S1 = datatypes.User_Customer{Address1: &A1, Address2: &A2}
				S2 = datatypes.User_Customer{Address1: &B1, Address2: nil}
			})

			It("return succ", func() {
				user.StructAssignment(&S1, &S2)
				Expect(*S1.Address1).To(Equal("21"))
				Expect(*S1.Address2).To(Equal("12"))
			})
		})
		Context("structAssignment", func() {
			BeforeEach(func() {
				S1 = datatypes.User_Customer{Address1: &A1, Address2: &A2}
				S2 = datatypes.User_Customer{Address1: &B1, Address2: &B2}
			})

			It("return succ", func() {
				user.StructAssignment(&S1, &S2)
				Expect(*S1.Address1).To(Equal("21"))
				Expect(*S1.Address2).To(Equal("22"))
			})
		})
		Context("structAssignment", func() {
			BeforeEach(func() {
				S1 = datatypes.User_Customer{Address1: nil, Address2: &A2}
				S2 = datatypes.User_Customer{Address1: &B1, Address2: &B2}
			})

			It("return succ", func() {
				user.StructAssignment(&S1, &S2)
				Expect(*S1.Address1).To(Equal("21"))
				Expect(*S1.Address2).To(Equal("22"))
			})
		})

		Context("structAssignment", func() {
			BeforeEach(func() {
				S1 = datatypes.User_Customer{Address1: &A1, Address2: &A2}
				S2 = datatypes.User_Customer{Address1: nil, Address2: &B2}
			})

			It("return succ", func() {
				user.StructAssignment(&S1, &S2)
				Expect(*S1.Address1).To(Equal("11"))
				Expect(*S1.Address2).To(Equal("22"))
			})
		})
	})
})
