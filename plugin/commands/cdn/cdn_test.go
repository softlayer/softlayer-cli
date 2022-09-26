package cdn_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cdn Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test cdn commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		cdn := cdn.SetupCobraCommands(slMeta)
		Expect(cdn.Name()).To(Equal("cdn"))
	})

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(cdn.CdnNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(cdn.CdnNamespace().Name).To(ContainSubstring("cdn"))
			Expect(cdn.CdnNamespace().Description).To(ContainSubstring("Classic infrastructure CDN commands"))
		})
	})
})
