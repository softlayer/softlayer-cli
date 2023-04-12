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

var _ = Describe("user apikey", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.ApikeyCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewApikeyCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})
	Describe("user apikey", func() {
		Context("Return error", func() {
			It("Set command without identifier", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid identifier", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
			})

			It("Set command without options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Please pass at least one of the flags."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.AddApiAuthenticationKeyReturns("", errors.New("Failed to add user's API authentication key"))
			})

			It("Failed to add user's API authentication key", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add user's API authentication key"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetApiAuthenticationKeysReturns([]datatypes.User_Customer_ApiAuthentication{}, errors.New("Failed to get user's API authentication keys"))
			})

			It("Failed to get user's API authentication keys", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user's API authentication keys"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetApiAuthenticationKeysReturns([]datatypes.User_Customer_ApiAuthentication{}, nil)
			})

			It("The user has not API authentication keys", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("The user has not API authentication keys"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerApiAuthenticationKey := []datatypes.User_Customer_ApiAuthentication{
					datatypes.User_Customer_ApiAuthentication{
						Id: sl.Int(111111),
					},
				}
				fakeUserManager.GetApiAuthenticationKeysReturns(fakerApiAuthenticationKey, nil)
				fakeUserManager.RemoveApiAuthenticationKeyReturns(false, errors.New("Failed to remove user's API authentication key"))
			})

			It("Failed to remove user's API authentication key with --remove flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove user's API authentication key"))
			})

			It("Failed to remove user's API authentication key with --refresh flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--refresh")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove user's API authentication key"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerApiAuthenticationKey := []datatypes.User_Customer_ApiAuthentication{
					datatypes.User_Customer_ApiAuthentication{
						Id: sl.Int(111111),
					},
				}
				fakeUserManager.GetApiAuthenticationKeysReturns(fakerApiAuthenticationKey, nil)
				fakeUserManager.RemoveApiAuthenticationKeyReturns(true, nil)
				fakeUserManager.AddApiAuthenticationKeyReturns("", errors.New("Failed to add user's API authentication key"))
			})

			It("Failed to add user's API authentication key with --refresh flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--refresh")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add user's API authentication key"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUserManager.AddApiAuthenticationKeyReturns("secretApiKey", nil)
			})

			It("Add API Authentication Key", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--add")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerApiAuthenticationKey := []datatypes.User_Customer_ApiAuthentication{
					datatypes.User_Customer_ApiAuthentication{
						Id: sl.Int(111111),
					},
				}
				fakeUserManager.GetApiAuthenticationKeysReturns(fakerApiAuthenticationKey, nil)
				fakeUserManager.RemoveApiAuthenticationKeyReturns(true, nil)
				fakeUserManager.AddApiAuthenticationKeyReturns("secretApiKey", nil)
			})

			It("Removed API Authentication Key", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})

			It("Refreshed API Authentication Key", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--refresh")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
