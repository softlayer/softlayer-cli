package block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Block Suite")
}


var _ = Describe("Test block Commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		Commands := block.SetupCobraCommands(slMeta)
		Expect(Commands.Name()).To(Equal("block"))
	})
	Context("Account Namespace", func() {
		It("Account Name Space", func() {
			Expect(block.BlockNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(block.BlockNamespace().Name).To(ContainSubstring("block"))
			Expect(block.BlockNamespace().Description).To(ContainSubstring("Classic infrastructure Block Storage"))
		})
	})
})
