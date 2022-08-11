package cdn_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cdn Suite")
}

var availableCommands = []string{
	"cdn-list",
	"cdn-detail",
	"cdn-edit",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test cdn.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := cdn.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in cdn.go")
			})
		}
	})

	Context("Cdn Namespace", func() {
		It("Cdn Name Space", func() {
			Expect(cdn.CdnNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(cdn.CdnNamespace().Name).To(ContainSubstring("cdn"))
			Expect(cdn.CdnNamespace().Description).To(ContainSubstring("Classic infrastructure CDN"))
		})
	})

	Context("Cdn MetaData", func() {
		It("Cdn MetaData", func() {
			Expect(cdn.CdnMetaData().Category).To(ContainSubstring("sl"))
			Expect(cdn.CdnMetaData().Name).To(ContainSubstring("cdn"))
			Expect(cdn.CdnMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl cdn"))
			Expect(cdn.CdnMetaData().Description).To(ContainSubstring("Classic infrastructure CDN"))
		})
	})
})
