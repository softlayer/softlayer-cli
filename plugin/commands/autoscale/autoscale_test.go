package autoscale_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autoscale Suite")
}

// These are all the commands in autoscale.go
var availableCommands = []string{
	"autoscale-edit",
	"autoscale-tag",
	"autoscale-logs",
	"autoscale-detail",
	"autoscale-list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test autoscale.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := autoscale.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in autoscale.go")
			})
		}
	})

	Context("AutoScale Namespace", func() {
		It("AutoScale Name Space", func() {
			Expect(autoscale.AutoScaleNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(autoscale.AutoScaleNamespace().Name).To(ContainSubstring("autoscale"))
			Expect(autoscale.AutoScaleNamespace().Description).To(ContainSubstring("Classic infrastructure Autoscale Group"))
		})
	})

	Context("User MetaData", func() {
		It("User MetaData", func() {
			Expect(autoscale.AutoScaleMetaData().Category).To(ContainSubstring("sl"))
			Expect(autoscale.AutoScaleMetaData().Name).To(ContainSubstring("autoscale"))
			Expect(autoscale.AutoScaleMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl autoscale"))
			Expect(autoscale.AutoScaleMetaData().Description).To(ContainSubstring("Classic infrastructure Autoscale Group"))
		})
	})
})
