package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestVirtual(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Server Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test Virtual Commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := virtual.SetupCobraCommands(slMeta)
		It("SetupCobraCommand works", func() {
			Expect(commands.Name()).To(Equal("vs"))
		})

	})
	Context("Virtual Namespace", func() {
		It("Virtual Namespace Exists", func() {
			Expect(virtual.VSNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(virtual.VSNamespace().Name).To(ContainSubstring("vs"))
			Expect(virtual.VSNamespace().Description).To(ContainSubstring("Classic infrastructure Virtual Servers"))
		})
	})
})
