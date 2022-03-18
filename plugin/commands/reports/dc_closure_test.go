package reports_test

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Reports Datacenter-Closures", func() {
	var (
		fakeUI      *terminal.FakeUI
		cmd         *reports.DCClosuresCommand
		cliCommand  cli.Command
		fakeSession *session.Session
		handler     testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		handler = fakeSession.TransportHandler.(testhelpers.FakeTransportHandler)

		cmd = reports.NewDCClosuresCommand(fakeUI, fakeSession)
		cliCommand = cli.Command{
			Name:        reports.DCClosuresMetaData().Name,
			Description: reports.DCClosuresMetaData().Description,
			Usage:       reports.DCClosuresMetaData().Usage,
			Flags:       reports.DCClosuresMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Datacenter-Closures Testing", func() {
		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("imageTest.ibmtest.com"))
			})
			It("Outputs JSON", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("\"Name\": \"imageTest.ibmtest.com\""))
			})
		})
		Context("Error Handling", func() {
			It("SoftLayer_Search::advancedSearch() Error", func() {
				handler.AddApiError("SoftLayer_Search", "advancedSearch", 500, "BAD")
				err := testhelpers.RunCommand(cliCommand, "--output=JSON")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD: BAD (HTTP 500)"))
			})
			It("SoftLayer_Network_Pod::getAllObjects() Error", func() {
				handler.AddApiError("SoftLayer_Network_Pod", "getAllObjects", 500, "ERRRR")
				fmt.Printf("API ERRORS ARE NOW\n%v", handler.ErrorMap)
				err := testhelpers.RunCommand(cliCommand, "--output=JSON")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ERRRR: ERRRR (HTTP 500)"))
			})
			It("Outputs NOT JSON", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=boson")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			AfterEach(func() {
				handler.ClearErrors()
			})
		})
	})
})
