package search_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/search"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SEarch Suite")
}

var availableCommands = []string{
	"types",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test search commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession([]string{"advancedSearch-allTypes"})
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	cliCommand := search.SetupCobraCommands(slMeta)
	cliCommand.PersistentFlags().Var(slMeta.OutputFlag, "output", "--output=JSON for json output.")
	// fakeSearchManager := new(testhelpers.FakeSearchManager)

	Context("New commands testable", func() {
		commands := search.SetupCobraCommands(slMeta)

		var arrayCommands = []string{}
		for _, command := range commands.Commands() {
			commandName := command.Name()
			arrayCommands = append(arrayCommands, commandName)
			It("available commands "+commands.Name(), func() {
				available := false
				if utils.StringInSlice(commandName, availableCommands) != -1 {
					available = true
				}
				Expect(available).To(BeTrue(), commandName+" not found in array available Commands")
			})
		}
		for _, command := range availableCommands {
			commandName := command
			It("ibmcloud sl "+commands.Name(), func() {
				available := false
				if utils.StringInSlice(commandName, arrayCommands) != -1 {
					available = true
				}
				Expect(available).To(BeTrue(), commandName + " not found in ibmcloud sl " + commands.Name())
			})
		}
	})

	Context("Seach Namespace", func() {
		It("Search Name Space", func() {
			Expect(search.SearchNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(search.SearchNamespace().Name).To(ContainSubstring("search"))
			Expect(search.SearchNamespace().Description).To(ContainSubstring("Perform a query against the SoftLayer search database."))
		})
	})

	var fakeHandler *testhelpers.FakeTransportHandler
	BeforeEach(func() {
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
	})
	// Search command is a bit special is an actual command, not a command group like most others.
	Context("Search Command tests", func() {
		It("Basic Search Command", func() {
			err := testhelpers.RunCobraCommand(cliCommand)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Basic Search API Error", func() {
			fakeHandler.AddApiError("SoftLayer_Search", "advancedSearch", 500, "Search Error")
			err := testhelpers.RunCobraCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Search Error: Search Error (HTTP 500)"))
		})
		It("Basic Search Command with query", func() {
			err := testhelpers.RunCobraCommand(cliCommand , "-q", `"_objectTpye:SoftLayer_Virtual_Guest test.com`)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("SoftLayer_Network_Vlan"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("VLAN |match|       ID: 675037"))
		})
		It("Basic Search Command with JSON output", func() {
			err := testhelpers.RunCobraCommand(cliCommand, "--output=JSON", "-q", `"test.com"`)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"resourceType": "SoftLayer_Ticket"`))
			Expect(fakeUI.Outputs()).To(ContainSubstring(`"id": 85346218,`))
		})
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
})
