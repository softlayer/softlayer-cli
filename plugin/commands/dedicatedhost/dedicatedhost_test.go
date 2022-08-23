package dedicatedhost_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dedicatedhost Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test dedicatedhost.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		dedicatedhostCommands := dedicatedhost.SetupCobraCommands(slMeta)
		Expect(dedicatedhostCommands.Name()).To(Equal("dedicatedhost"))
	})
	Context("Dedicatedhost Namespace", func() {
		It("Dedicatedhost Name Space", func() {
			Expect(dedicatedhost.DedicatedhostNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(dedicatedhost.DedicatedhostNamespace().Name).To(ContainSubstring("dedicatedhost"))
			Expect(dedicatedhost.DedicatedhostNamespace().Description).To(ContainSubstring("Classic infrastructure Dedicatedhost"))
		})
	})
})
