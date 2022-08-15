package reports_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Report Suite")
}


var _ = Describe("Test report commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		commands := reports.SetupCobraCommands(slMeta)
		Expect(commands.Name()).To(Equal("report"))
	})

	Context("Report Namespace", func() {
		It("Report Name Space", func() {
			Expect(reports.ReportsNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(reports.ReportsNamespace().Name).To(ContainSubstring("report"))
			Expect(reports.ReportsNamespace().Description).To(ContainSubstring("Classic Infrastructure Reports"))
		})
	})
})
