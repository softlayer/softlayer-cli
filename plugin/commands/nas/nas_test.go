package nas_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/nas"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Network Attached Storage Suite")
}

// These are all the commands in nas.go
var availableCommands = []string{
	"nas-list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test nas.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := nas.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in nas.go")
			})
		}
	})

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(nas.NasNetworkStorageNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(nas.NasNetworkStorageNamespace().Name).To(ContainSubstring("nas"))
			Expect(nas.NasNetworkStorageNamespace().Description).To(ContainSubstring("Classic infrastructure Network Attached Storage"))
		})
	})

	Context("Network Attached Storage MetaData", func() {
		It("Network Attached Storage MetaData", func() {
			Expect(nas.NasNetworkStorageMetaData().Category).To(ContainSubstring("sl"))
			Expect(nas.NasNetworkStorageMetaData().Name).To(ContainSubstring("nas"))
			Expect(nas.NasNetworkStorageMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl nas"))
			Expect(nas.NasNetworkStorageMetaData().Description).To(ContainSubstring("Classic infrastructure Network Attached Storage"))
		})
	})
})
