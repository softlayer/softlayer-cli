package account_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}

var availableCommands = []string{
	"account-bandwidth-pools",
	"account-events",
	"account-event-detail",
	"account-invoices",
	"account-invoice-detail",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test account.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := account.GetCommandActionBindings(context, fakeUI, fakeSession)

	Context("Test Actions", func() {
		for _, cmdName := range availableCommands {
			//necessary to ensure the correct value is passed to the closure
			cmdName := cmdName
			It("ibmcloud sl "+cmdName, func() {
				command, exists := commands[cmdName]
				Expect(exists).To(BeTrue(), cmdName+" not found")
				// Checks to make sure we actually have a function here.
				// Test the actual function works in the specific commands test file.
				Expect(reflect.ValueOf(command).Kind().String()).To(Equal("func"))
				context := testhelpers.GetCliContext(cmdName)
				err := command(context)
				// some commands work without arguments
				if err == nil {
					Expect(err).NotTo(HaveOccurred())
				} else {
					Expect(err).To(HaveOccurred())
				}
			})
		}
	})

	Context("New commands testable", func() {
		for cmdName, _ := range commands {
			//necessary to ensure the correct value is passed to the closure
			cmdName := cmdName
			It("availableCommands["+cmdName+"]", func() {
				found := false
				for _, value := range availableCommands {
					if value == cmdName {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in account.go")
			})
		}
	})

	Context("Account Namespace", func() {
		It("Account Name Space", func() {
			Expect(account.AccountNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(account.AccountNamespace().Name).To(ContainSubstring("account"))
			Expect(account.AccountNamespace().Description).To(ContainSubstring("Classic infrastructure Account"))
		})
	})

	Context("Account MetaData", func() {
		It("Account MetaData", func() {
			Expect(account.AccountMetaData().Category).To(ContainSubstring("sl"))
			Expect(account.AccountMetaData().Name).To(ContainSubstring("account"))
			Expect(account.AccountMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl account"))
			Expect(account.AccountMetaData().Description).To(ContainSubstring("Classic infrastructure Account"))
		})
	})
})
