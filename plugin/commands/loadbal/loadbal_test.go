package loadbal_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoadBalancer Suite")
}

// These are all the commands in loadbal.go
var availableCommands = []string{
	"loadbal-cancel",
	"loadbal-detail",
	"loadbal-health-edit",
	"loadbal-l7member-add",
	"loadbal-l7member-delete",
	"loadbal-l7policies",
	"loadbal-l7policy-add",
	"loadbal-l7policy-delete",
	"loadbal-l7policy-edit",
	"loadbal-l7pool-add",
	"loadbal-l7pool-delete",
	"loadbal-l7pool-detail",
	"loadbal-l7pool-edit",
	"loadbal-l7rule-add",
	"loadbal-l7rule-delete",
	"loadbal-l7rules",
	"loadbal-list",
	"loadbal-member-add",
	"loadbal-member-delete",
	"loadbal-order",
	"loadbal-order-options",
	"loadbal-protocol-add",
	"loadbal-protocol-delete",
	"loadbal-protocol-edit",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test loadbal.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := loadbal.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in loadbal.go")
			})
		}
	})

	Context("Loadbal Namespace", func() {
		It("Loadbal Name Space", func() {
			Expect(loadbal.LoadbalNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(loadbal.LoadbalNamespace().Name).To(ContainSubstring("loadbal"))
			Expect(loadbal.LoadbalNamespace().Description).To(ContainSubstring("Classic infrastructure Load Balancers"))
		})
	})

	Context("Loadbal MetaData", func() {
		It("Loadbal MetaData", func() {
			Expect(loadbal.LoadbalMetaData().Category).To(ContainSubstring("sl"))
			Expect(loadbal.LoadbalMetaData().Name).To(ContainSubstring("loadbal"))
			Expect(loadbal.LoadbalMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl loadbal"))
			Expect(loadbal.LoadbalMetaData().Description).To(ContainSubstring("Classic infrastructure Load Balancers"))
		})
	})
})
