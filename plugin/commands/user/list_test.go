package user_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("sl user list", func() {
	var (
		fakeUI *terminal.FakeUI
		// fakeUserManager *testhelpers.FakeUserManager
		cliCommand  *user.ListCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
		fakeHandler *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		// fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		// Clear API call logs and any errors that might have been set after every test
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("Usage Errors ", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "noExist")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column noExist is not supported."))
		})
	})

	Describe("API Errors", func() {
		It("SoftLayer_Account::getUsers API Error", func() {
			fakeHandler.AddApiError("SoftLayer_Account", "getUsers", 500, "Internal Server Error")
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to list users."))
		})
	})

	Describe("Happy Path", func() {
		It("return users list", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).NotTo(HaveOccurred())
			// Remove whitespace to make testing output a bit less rigid
			trimmed := strings.ReplaceAll(fakeUI.Outputs(), " ", "")
			Expect(trimmed).To(ContainSubstring("idusernameemaildisplayName2FAclassicAPIKeyvpn"))
			Expect(trimmed).To(ContainSubstring("1468361IBM27821sdfasdfasd@one.comazzz---"))
			Expect(trimmed).To(ContainSubstring("1468362IBM27832sdfasdfasd@two.comaccc-yesYes"))
			Expect(trimmed).To(ContainSubstring("1468363IBM27843sdfasdfasd@three.comabc--No"))
		})
		It("return users list simple columns", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--column=id", "--column=displayName")
			Expect(err).NotTo(HaveOccurred())
			// Remove whitespace to make testing output a bit less rigid
			trimmed := strings.ReplaceAll(fakeUI.Outputs(), " ", "")
			Expect(trimmed).To(ContainSubstring("iddisplayName"))
			Expect(trimmed).To(ContainSubstring("1468361azzz"))
			Expect(trimmed).To(ContainSubstring("1468362accc"))
			Expect(trimmed).To(ContainSubstring("1468363abc"))
		})
		It("return users list with just username", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "username")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("username"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("IBM27821"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("IBM27832"))
		})
		It("return users list with json output", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "sdfasdfasd@three.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"displayName": "azzz",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"apiAuthenticationKeyCount": 1,`))
		})
	})
})
