package firewall_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firewall Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test firewall commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		licensesCommands := firewall.SetupCobraCommands(slMeta)
		Expect(licensesCommands.Name()).To(Equal("firewall"))
	})

	Context("Firewall Namespace", func() {
		It("Firewall Name Space", func() {
			Expect(firewall.FirewallNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(firewall.FirewallNamespace().Name).To(ContainSubstring("firewall"))
			Expect(firewall.FirewallNamespace().Description).To(ContainSubstring("Classic infrastructure Firewalls"))
		})
	})
})
