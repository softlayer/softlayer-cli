package globalip_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GlobalIP Suite")
}

var _ = Describe("Test globalip commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		globalipCommands := globalip.SetupCobraCommands(slMeta)
		Expect(globalipCommands.Name()).To(Equal("globalip"))
	})

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(globalip.GlobalIpNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(globalip.GlobalIpNamespace().Name).To(ContainSubstring("globalip"))
			Expect(globalip.GlobalIpNamespace().Description).To(ContainSubstring("Classic infrastructure Global IP addresses"))
		})
	})
})
