package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
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
	Describe("User Create Command", func() {
		Context("Invalid Paramter Checks", func() {
			It("Needs one argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Input Checks", func() {
			It("Not Y/N", func() {
				fakeUI.Inputs("123456")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
			It("No confirmation", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
		Context("Error Handling", func() {
			It("API Error", func() {
				fakeUserManager.CreateUserReturns(datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add user."))
			})
			It("Bad Template", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f", "--template", ``)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to unmarshal template json: unexpected end of JSON input"))
			})
		})

		Context("Happy Path Tests", func() {
			It("Create a user --force", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
			It("Create a user with confirmation", func() {
				fakeUI.Inputs("Y")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "MyPassWord")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("You are about to create the following user: createdUser@email.com. Do you wish to continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
			It("Create a user from another user", func() {
				fakeUI.Inputs("Y")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--from-user", "456", "--password", "MyPassWord", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password   MyPassWord"))
			})
			It("Create a user from a template", func() {
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
			It("Create a user with a generated password", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "createdUser@email.com", "--email", "createdUser@email.com", "--password", "generate", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name       value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username   createdUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email      createdUser@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password"))
			})
		})
	})

	// dataValues are a set of 4 strings we set Default and UserValues to
	// expected is a set of 2 strings that we check were set properly
	DescribeTable("StructAssignment Tests",
		func(dataValues []string, expected []string) {
			Expect(len(dataValues)).To(Equal(4))
			Expect(len(expected)).To(Equal(2))
			Default := datatypes.User_Customer{Address1: &dataValues[0], Address2: &dataValues[1]}
			UserValues := datatypes.User_Customer{Address1: &dataValues[2], Address2: &dataValues[3]}
			// Can't set nil in the dataValues value because its a string, so we just do this
			if dataValues[0] == "nil" {
				Default.Address1 = nil
			}
			if dataValues[1] == "nil" {
				Default.Address2 = nil
			}
			if dataValues[2] == "nil" {
				UserValues.Address1 = nil
			}
			if dataValues[3] == "nil" {
				UserValues.Address2 = nil
			}
			user.StructAssignment(&Default, &UserValues)
			Expect(*Default.Address1).To(Equal(expected[0]))
			Expect(*Default.Address2).To(Equal(expected[1]))
		},
		Entry("Test1", []string{"Def1", "Def2", "UserInput1", "nil"}, []string{"UserInput1", "Def2"}),
		Entry("Test2", []string{"Def1", "Def2", "nil", "UserInput2"}, []string{"Def1", "UserInput2"}),
		Entry("Test3", []string{"Def1", "nil", "UserInput1", "UserInput2"}, []string{"UserInput1", "UserInput2"}),
		Entry("Test4", []string{"nil", "Def2", "UserInput1", "UserInput2"}, []string{"UserInput1", "UserInput2"}),
	)
})
