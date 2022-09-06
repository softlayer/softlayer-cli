package dns_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DNS Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test DNS commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		dnsCommands := dns.SetupCobraCommands(slMeta)
		Expect(dnsCommands.Name()).To(Equal("dns"))
	})
	Context("DNS Namespace", func() {
		It("DNS Name Space", func() {
			Expect(dns.DnsNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(dns.DnsNamespace().Name).To(ContainSubstring("dns"))
			Expect(dns.DnsNamespace().Description).To(ContainSubstring("Classic infrastructure Domain Name System"))
		})
	})
})
