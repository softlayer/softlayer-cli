package account_test

import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test account.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		accountCommands := account.SetupCobraCommands(slMeta)
		Expect(accountCommands.Name()).To(Equal("account"))
	})
	Context("Account Namespace", func() {
		It("Account Name Space", func() {
			Expect(account.AccountNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(account.AccountNamespace().Name).To(ContainSubstring("account"))
			Expect(account.AccountNamespace().Description).To(ContainSubstring("Classic infrastructure Account"))
		})
	})
})
