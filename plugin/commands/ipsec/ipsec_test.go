package ipsec_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IPSEC Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test ipsec commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		ipsecCommands := ipsec.SetupCobraCommands(slMeta)
		Expect(ipsecCommands.Name()).To(Equal("ipsec"))
	})

	Context("ipsec Namespace", func() {
		It("ipsec Name Space", func() {
			Expect(ipsec.IpsecNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(ipsec.IpsecNamespace().Name).To(ContainSubstring("ipsec"))
			Expect(ipsec.IpsecNamespace().Description).To(ContainSubstring("Classic infrastructure IPSEC VPN"))
		})
	})
})
