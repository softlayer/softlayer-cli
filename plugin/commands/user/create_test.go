package user_test

import (
	"strings"

	// for ContainSubstrings()
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"

	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/user"
)

var _ = Describe("Create", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.CreateCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewCreateCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        metadata.UserCreateMetaData().Name,
			Description: metadata.UserCreateMetaData().Description,
			Usage:       metadata.UserCreateMetaData().Usage,
			Flags:       metadata.UserCreateMetaData().Flags,
			Action:      cmd.Run,
		}
		testUser := datatypes.User_Customer{
			Id: sl.Int(5555),
			Username: sl.String("ATestUser"),
			Email: sl.String("user@email.com"),
		}
		fakeUserManager.GetCurrentUserReturns(testUser, nil)
		fakeUserManager.CreateUserReturns(testUser, nil)
	})
	Describe("user create", func() {
		Context("user create with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
		Context("Basic User Create usage", func() {
			It("Create a user", func() {
				err := testhelpers.RunCommand(cliCommand, "user@email.com", "--password", "MyPassWord",  "--email", "user@email.com", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"user@email.com"}))
			})
		})
	})

	Describe("password", func() {
		Context("generatePassword", func() {
			It("return succ", func() {
				password := user.GeneratePassword(23, 4)
				Expect(len(password)).To(Equal(23))
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
