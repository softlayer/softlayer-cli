package user_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var testUser datatypes.User_Customer
var _ = Describe("sl user detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		cliCommand      *user.DetailsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
		fakeHandler     *testhelpers.FakeTransportHandler
		fakeUserManager *testhelpers.FakeUserManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		fakeUserManager = new(testhelpers.FakeUserManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewDetailsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})
	AfterEach(func() {
		// Clear API call logs and any errors that might have been set after every test
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("Usage Errors", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
		})
	})

	Describe("API Errors", func() {
		It("SoftLayer_User_Customer::getObject Exception", func() {
			fakeHandler.AddApiError("SoftLayer_User_Customer", "getObject", 500, "Internal Server Error")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show user detail."))
		})
		It("SoftLayer_User_Customer::getPermissions Exception", func() {
			fakeHandler.AddApiError("SoftLayer_User_Customer", "getPermissions", 500, "Internal Server Error")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--permissions")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show user permissions."))
		})
		It("SoftLayer_User_Customer::getLoginAttempts Exception", func() {
			fakeHandler.AddApiError("SoftLayer_User_Customer", "getLoginAttempts", 500, "Internal Server Error")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--logins")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show login history."))
		})
		It("SoftLayer_Event_Log::getAllObjects Exception", func() {
			fakeHandler.AddApiError("SoftLayer_Event_Log", "getAllObjects", 500, "Internal Server Error")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--events")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show event log."))
		})
	})

	Describe("Happy Path Tests", func() {
		It("return a user", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("XXX.ASD@ibm.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login     2018-11-09T00:40:58+08:00 From: 169.60.96.34"))
			Expect(fakeUI.Outputs()).NotTo(ContainSubstring("StringKeyAuthentication"))
			Expect(fakeUI.Outputs()).NotTo(ContainSubstring("KEY_PERMISSION_2"))
		})
		It("return a user with apikey", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--keys")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("name           value"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID             345234"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("StringKeyAuthentication"))
		})
		It("return a user with permissions", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--permissions")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("APIKEY         Yes"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ACCESS_ALL_HARDWARE"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("HARDWARE_VIEW"))
		})
		It("return a user with logins", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--logins")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("asdfgn@ibm.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("2018-11-08T16:40:58Z   169.60.96.34    true"))
		})
		It("return a user with events", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--events")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("11111111 aaaa Ave - Markham ON CA L6G1C7"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("IAM Token validation successful"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("169.1.98.6"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("123_scaparro@ibm.com"))
		})
	})

	// Since this CLI makes 2 calls to the same API, we need to use the Fake Manager to handle that.
	Describe("API Errors with a Fake Manager", func() {
		var testUser datatypes.User_Customer
		BeforeEach(func() {
			testUser = datatypes.User_Customer{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_User_Customer", "getObject", nil, nil, &testUser,
			)
			Expect(txError).NotTo(HaveOccurred())
			fakeUserManager.GetUserReturnsOnCall(0, testUser, nil)
			cliCommand.UserManager = fakeUserManager
		})
		It("SoftLayer_User_Customer::getObject Second Exception Hardware", func() {
			fakeUserManager.GetUserReturnsOnCall(1, datatypes.User_Customer{}, errors.New("BAD HARDWARE"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--hardware")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show hardware."))
		})
		It("SoftLayer_User_Customer::getObject Second Exception Virtual", func() {
			fakeUserManager.GetUserReturnsOnCall(1, datatypes.User_Customer{}, errors.New("BAD VIRTUAL"))
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--virtual")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to show virual server."))
		})
	})

	Describe("Happy Path with Fake Manager", func() {
		var testUser datatypes.User_Customer
		BeforeEach(func() {
			testUser = datatypes.User_Customer{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_User_Customer", "getObject", nil, nil, &testUser,
			)
			Expect(txError).NotTo(HaveOccurred())
			fakeUserManager.GetUserReturnsOnCall(0, testUser, nil)
			cliCommand.UserManager = fakeUserManager
		})
		It("return a user with hardware", func() {
			userHardware := []datatypes.Hardware{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_Account", "getHardware", nil, nil, &userHardware,
			)
			testUser.Hardware = userHardware
			fakeUserManager.GetUserReturnsOnCall(1, testUser, nil)
			Expect(txError).NotTo(HaveOccurred())
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--hardware")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("XXX.ASD@ibm.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ibmcloud-cli-dev1.ibm.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ibmcloud-cli-dev2.ibm.com"))
		})
		It("return a user with virtual", func() {
			userGuests := []datatypes.Virtual_Guest{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_Account", "getVirtualGuests", nil, nil, &userGuests,
			)
			testUser.VirtualGuests = userGuests
			fakeUserManager.GetUserReturnsOnCall(1, testUser, nil)
			Expect(txError).NotTo(HaveOccurred())
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--virtual")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("XXX.ASD@ibm.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("3169-2-stemcell-for-dirtycow.softlayer.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("3263-10-1-stemcell-bluemix.softlayer.com"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("3263-10-2-stemcell-bluemix.softlayer.com"))
		})
		It("JSON with --hardware", func() {
			userHardware := []datatypes.Hardware{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_Account", "getHardware", nil, nil, &userHardware,
			)
			testUser.Hardware = userHardware
			fakeUserManager.GetUserReturnsOnCall(1, testUser, nil)
			Expect(txError).NotTo(HaveOccurred())
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json", "--hardware")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Virtual": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"hostname": "ibmcloud-cli-dev1"`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Events": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Permissions": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})
		It("JSON with --virtual", func() {
			userGuests := []datatypes.Virtual_Guest{}
			txError := fakeHandler.DoRequest(
				fakeSession, "SoftLayer_Account", "getVirtualGuests", nil, nil, &userGuests,
			)
			testUser.VirtualGuests = userGuests
			fakeUserManager.GetUserReturnsOnCall(1, testUser, nil)
			Expect(txError).NotTo(HaveOccurred())
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json", "--virtual")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"hostName": "3169-2-stemcell-for-dirtycow",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Hardware": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Events": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Permissions": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})
		It("return a user without apikey", func() {
			testUser.ApiAuthenticationKeys = []datatypes.User_Customer_ApiAuthentication{}
			fakeUserManager.GetUserReturnsOnCall(0, testUser, nil)
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--keys")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("Status         Active"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("APIKEY         No"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Username       test"))
		})
	})
	Describe("JSON Output tests", func() {
		It("Just user details", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Virtual": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Hardware": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Events": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Permissions": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})
		It("JSON with --permissions", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json", "--permissions")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Virtual": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Hardware": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Events": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"keyName": "ACCESS_ALL_DEDICATEDHOSTS",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})

		It("JSON with --logins", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json", "--logins")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Virtual": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Hardware": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Events": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Permissions": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"ipAddress": "169.60.96.34",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": [`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})
		It("JSON with --events", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--output=json", "--events")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": {`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"email": "XXX.ASD@ibm.com",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Virtual": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Hardware": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"eventName": "IAM Token validation successful",`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Permissions": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"Logins": null,`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"DedicatedHosts": null`))
		})
	})
})
