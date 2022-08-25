package autoscale_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autoscale Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test autoscale commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		autoscaleCommands := autoscale.SetupCobraCommands(slMeta)
		Expect(autoscaleCommands.Name()).To(Equal("autoscale"))
	})

	Context("AutoScale Namespace", func() {
		It("AutoScale Name Space", func() {
			Expect(autoscale.AutoScaleNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(autoscale.AutoScaleNamespace().Name).To(ContainSubstring("autoscale"))
			Expect(autoscale.AutoScaleNamespace().Description).To(ContainSubstring("Classic infrastructure Autoscale Group"))
		})
	})
})
