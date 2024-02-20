package callapi_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CallAPI Suite")
}

var _ = Describe("CallAPI test", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *callapi.CallAPICommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
		fakeManager *testhelpers.FakeCallAPIManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeManager = new(testhelpers.FakeCallAPIManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = callapi.NewCallAPICommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.CallAPIManager = fakeManager
	})

	Describe("CallAPI command", func() {
		Context("CallAPI, Invalid Usage", func() {
			It("Set command without arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("accepts 2 arg(s), received 0"))
			})
			It("Set command with 1 argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Hardware_Server")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("accepts 2 arg(s), received 1"))
			})

		})

		Context("CallAPI, empty datas", func() {
			It("Return empty bytes", func() {
				fakeManager.CallAPIReturns([]byte{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Hardware_Server", "getVirtualGuests")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Null"))
			})
			It("Return null bytes", func() {
				fakeManager.CallAPIReturns(nil, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Hardware_Server", "getVirtualGuests")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Null"))
			})
		})

		Context("CallAPI account getUsers", func() {
			It("return users", func() {
				response := `[
	{
		"accountId": 12345,
		"address1": "4849 Alpha Rd",
		"city": "Dallas",
		"companyName": "SoftLayer Internal - Development Community",
		"country": "US",
		"createDate": "2022-11-14T10:27:12-06:00",
		"displayName": "UserTest",
		"email": "UserTest@test.com","
		"id": 654321,
		"userStatusId": 1001,
		"username": "1234_UserTest@test.com",
	},
]`
				arrayBytes := []byte(response)
				fakeManager.CallAPIReturns(arrayBytes, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "Account", "getUsers")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 654321,`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"username": "1234_UserTest@test.com",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"accountId": 12345,`))
			})
		})
	})
})
