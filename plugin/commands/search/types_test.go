package search_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/search"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("search types", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *search.SearchTypesCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
		fakeHandler *testhelpers.FakeTransportHandler
		// fakeSearchManager *testhelpers.FakeSearchManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = search.NewSearchTypesCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		// fakeSearchManager = new(testhelpers.FakeSearchManager)
		// cliCommand.SearchManager = fakeSearchManager
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
	})

	Describe("search types tests", func() {
		Context("Basic usage", func() {
			It("sl search types", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("SoftLayer_Network_Vlan"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("accountId    true       integer"))
			})
			It("sl search types JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "SoftLayer_Hardware"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"name": "accountId",`))
			})
			It("API errors", func() {
				fakeHandler.AddApiError("SoftLayer_Search", "getObjectTypes", 500, "Search Error")
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Search Error: Search Error (HTTP 500)"))
			})
		})
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
})
