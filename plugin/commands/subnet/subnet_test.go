package subnet_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subnet Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test subnet commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		subnetCommands := subnet.SetupCobraCommands(slMeta)
		Expect(subnetCommands.Name()).To(Equal("subnet"))
	})
	Context("Subnet Namespace", func() {
		It("Subnet Name Space", func() {
			Expect(subnet.SubnetNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(subnet.SubnetNamespace().Name).To(ContainSubstring("subnet"))
			Expect(subnet.SubnetNamespace().Description).To(ContainSubstring("Classic infrastructure Network subnets"))
		})
	})
})
