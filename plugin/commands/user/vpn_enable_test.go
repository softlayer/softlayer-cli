package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("user vpn-enable", func() {
	var (
		fakeUI          *terminal.FakeUI
		cliCommand      *user.VpnEnableCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
		fakeUserManager *testhelpers.FakeUserManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeUserManager = new(testhelpers.FakeUserManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewVpnEnableCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager

		testUser = datatypes.User_Customer{
			SslVpnAllowedFlag: sl.Bool(true),
		}
	})

	Describe("user vpn-enable", func() {

		Context("Return error", func() {
			It("Set command without Argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid user Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'User ID'. It must be a positive integer."))
			})
		})

		Context("Account cancel-item, softlayer errors", func() {
			It("Set command with unknow item ID", func() {
				fakeUserManager.EditUserReturns(true, errors.New("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ObjectNotFound: Unable to find object with id of '123'. (HTTP 404)"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUserManager.EditUserReturns(true, nil)
			})
			It("enable", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
