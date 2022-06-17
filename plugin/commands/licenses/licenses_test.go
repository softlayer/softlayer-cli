package licenses_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Licenses Suite")
}

var availableCommands = []string{
	"licenses-create-options",
	"licenses-create",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test licenses.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := licenses.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in licenses.go")
			})
		}
	})

	Context("Licenses Namespace", func() {
		It("Licenses Name Space", func() {
			Expect(licenses.LicensesNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(licenses.LicensesNamespace().Name).To(ContainSubstring("licenses"))
			Expect(licenses.LicensesNamespace().Description).To(ContainSubstring("Classic infrastructure Licenses"))
		})
	})

	Context("Licenses MetaData", func() {
		It("Licenses MetaData", func() {
			Expect(licenses.LicensesMetaData().Category).To(ContainSubstring("sl"))
			Expect(licenses.LicensesMetaData().Name).To(ContainSubstring("licenses"))
			Expect(licenses.LicensesMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl licenses"))
			Expect(licenses.LicensesMetaData().Description).To(ContainSubstring("Classic infrastructure Licenses"))
		})
	})
})
