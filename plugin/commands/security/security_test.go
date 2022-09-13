package security_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Security Suite")
}

// These are all the commands in security.go
var availableCommands = []string{
	"security-cert-add",
	"security-cert-download",
	"security-cert-edit",
	"security-cert-list",
	"security-cert-remove",
	"security-sshkey-add",
	"security-sshkey-edit",
	"security-sshkey-list",
	"security-sshkey-print",
	"security-sshkey-remove",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test security commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		security := security.SetupCobraCommands(slMeta)
		Expect(security.Name()).To(Equal("security"))
	})

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(security.SecurityNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(security.SecurityNamespace().Name).To(ContainSubstring("security"))
			Expect(security.SecurityNamespace().Description).To(ContainSubstring("Classic infrastructure SSH Keys and SSL Certificates"))
		})
	})
})
