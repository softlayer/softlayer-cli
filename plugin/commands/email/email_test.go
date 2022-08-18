package email_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Email Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test eventlog commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		eventlogCommands := email.SetupCobraCommands(slMeta)
		Expect(eventlogCommands.Name()).To(Equal("email"))
	})

	Context("Email Namespace", func() {
		It("Email Name Space", func() {
			Expect(email.EmailNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(email.EmailNamespace().Name).To(ContainSubstring("email"))
			Expect(email.EmailNamespace().Description).To(ContainSubstring("Classic infrastructure Email"))
		})
	})
})
