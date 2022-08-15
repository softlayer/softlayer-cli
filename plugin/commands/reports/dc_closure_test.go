package reports_test

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)



var _ = Describe("Reports Datacenter-Closures", func() {
    var (
        fakeUI          *terminal.FakeUI
        cliCommand      *reports.DCClosuresCommand
        fakeSession     *session.Session
        slCommand       *metadata.SoftlayerCommand
        fakeHandler *testhelpers.FakeTransportHandler
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        fakeHandler = testhelpers.GetSessionHandler(fakeSession)
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = reports.NewDCClosuresCommand(slCommand)
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
    })
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("Datacenter-Closures Testing", func() {
		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("imageTest.ibmtest.com"))
			})
			It("Outputs JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("\"Name\": \"imageTest.ibmtest.com\""))
			})
		})
		Context("Error Handling", func() {
			It("SoftLayer_Search::advancedSearch() Error", func() {
				fakeHandler.AddApiError("SoftLayer_Search", "advancedSearch", 500, "BAD")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
			It("SoftLayer_Network_Pod::getAllObjects() Error", func() {
				fakeHandler.AddApiError("SoftLayer_Network_Pod", "getAllObjects", 500, "ERRRR")
				fmt.Printf("API ERRORS ARE NOW\n%v", fakeHandler.ErrorMap)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ERRRR: ERRRR (HTTP 500)"))
			})
			It("Outputs NOT JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=boson")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			AfterEach(func() {
				fakeHandler.ClearErrors()
			})
		})
	})
})
