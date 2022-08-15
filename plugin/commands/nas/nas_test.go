package nas_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/nas"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Network Attached Storage Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test nas.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		nas := nas.SetupCobraCommands(slMeta)
		Expect(nas.Name()).To(Equal("nas"))
	})

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(nas.NasNetworkStorageNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(nas.NasNetworkStorageNamespace().Name).To(ContainSubstring("nas"))
			Expect(nas.NasNetworkStorageNamespace().Description).To(ContainSubstring("Classic infrastructure Network Attached Storage"))
		})
	})
})
