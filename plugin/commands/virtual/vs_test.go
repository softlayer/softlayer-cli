package virtual_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestVirtual(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Server Suite")
}

// These are all the commands in virtual.go
var availableCommands = []string{
	"vs-authorize-storage",
	"vs-cancel",
	"vs-capture",
	"vs-create",
	"vs-host-create",
	"vs-options",
	"vs-credentials",
	"vs-detail",
	"vs-dns-sync",
	"vs-edit",
	"vs-list",
	"vs-host-list",
	"vs-migrate",
	"vs-pause",
	"vs-power-off",
	"vs-power-on",
	"vs-ready",
	"vs-billing",
	"vs-reboot",
	"vs-reload",
	"vs-rescue",
	"vs-resume",
	"vs-upgrade",
	"vs-capacity-create-options",
	"vs-capacity-detail",
	"vs-bandwidth",
	"vs-storage",
	"vs-placementgroup-list",
	"vs-placementgroup-create-options",
	"vs-placementgroup-create",
	"vs-capacity-list",
	"vs-capacity-create",
	"vs-usage",
	"vs-placementgroup-details",
	"vs-monitoring-list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test virtual.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := virtual.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				context := testhelpers.GetCliContextHelp(cmdName)
				errSet := context.GlobalSet("help", "true")
				Expect(errSet).NotTo(HaveOccurred())
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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in virtual.go")
			})
		}
	})

	Context("Virtual Namespace", func() {
		It("virtual Name Space", func() {
			Expect(virtual.VSNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(virtual.VSNamespace().Name).To(ContainSubstring("vs"))
			Expect(virtual.VSNamespace().Description).To(ContainSubstring("Classic infrastructure Virtual Servers"))
		})
	})

	Context("Virtual MetaData", func() {
		It("virtual MetaData", func() {
			Expect(virtual.VSMetaData().Category).To(ContainSubstring("sl"))
			Expect(virtual.VSMetaData().Name).To(ContainSubstring("vs"))
			Expect(virtual.VSMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl vs"))
			Expect(virtual.VSMetaData().Description).To(ContainSubstring("Classic infrastructure Virtual Servers"))
		})
	})
})
