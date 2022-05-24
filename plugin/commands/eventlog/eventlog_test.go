package eventlog_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Log Suite")
}

// These are all the commands in eventlog.go
var availableCommands = []string{
	"event-log-get",
	"event-log-types",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test eventlog.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := eventlog.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in eventlog.go")
			})
		}
	})

	Context("Event Log Namespace", func() {
		It("Event Log Name Space", func() {
			Expect(eventlog.EventLogNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(eventlog.EventLogNamespace().Name).To(ContainSubstring("event-log"))
			Expect(eventlog.EventLogNamespace().Description).To(ContainSubstring("Classic infrastructure Event Log Group"))
		})
	})

	Context("User MetaData", func() {
		It("User MetaData", func() {
			Expect(eventlog.EventLogMetaData().Category).To(ContainSubstring("sl"))
			Expect(eventlog.EventLogMetaData().Name).To(ContainSubstring("event-log"))
			Expect(eventlog.EventLogMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl event-log"))
			Expect(eventlog.EventLogMetaData().Description).To(ContainSubstring("Classic infrastructure Event Log Group"))
		})
	})
})
