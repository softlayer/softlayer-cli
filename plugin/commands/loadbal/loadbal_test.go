package loadbal_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoadBalancer Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test load balancer commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		loadbalancerCommands := loadbal.SetupCobraCommands(slMeta)
		Expect(loadbalancerCommands.Name()).To(Equal("loadbal"))
	})
	Context("LoadBalancer Namespace", func() {
		It("LoadBalancer Name Space", func() {
			Expect(loadbal.LoadbalNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(loadbal.LoadbalNamespace().Name).To(ContainSubstring("loadbal"))
			Expect(loadbal.LoadbalNamespace().Description).To(ContainSubstring("Classic infrastructure Load Balancers"))
		})
	})
})
