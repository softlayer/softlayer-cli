package licenses_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Licenses Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test licenses commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		licensesCommands := licenses.SetupCobraCommands(slMeta)
		Expect(licensesCommands.Name()).To(Equal("licenses"))
	})

	Context("Licenses Namespace", func() {
		It("Licenses Name Space", func() {
			Expect(licenses.LicensesNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(licenses.LicensesNamespace().Name).To(ContainSubstring("licenses"))
			Expect(licenses.LicensesNamespace().Description).To(ContainSubstring("Classic infrastructure Licenses"))
		})
	})
})
