package vlan_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VLAN Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test vlan commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		vlanCommands := vlan.SetupCobraCommands(slMeta)
		Expect(vlanCommands.Name()).To(Equal("vlan"))
	})
	Context("Vlan Namespace", func() {
		It("Vlan Name Space", func() {
			Expect(vlan.VlanNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(vlan.VlanNamespace().Name).To(ContainSubstring("vlan"))
			Expect(vlan.VlanNamespace().Description).To(ContainSubstring("Classic infrastructure Network VLANs"))
		})
	})
})
