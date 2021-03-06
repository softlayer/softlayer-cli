package reports_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Report Suite")
}

var availableCommands = []string{
	"report-datacenter-closures",
	"report-bandwidth",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test report.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := reports.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in user.go")
			})
		}
	})

	Context("Report Namespace", func() {
		It("Report Name Space", func() {
			Expect(reports.ReportsNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(reports.ReportsNamespace().Name).To(ContainSubstring("report"))
			Expect(reports.ReportsNamespace().Description).To(ContainSubstring("Classic Infrastructure Reports"))
		})
	})

	Context("User MetaData", func() {
		It("User MetaData", func() {
			Expect(reports.ReportsMetaData().Category).To(ContainSubstring("sl"))
			Expect(reports.ReportsMetaData().Name).To(ContainSubstring("report"))
			Expect(reports.ReportsMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl report"))
			Expect(reports.ReportsMetaData().Description).To(ContainSubstring("Classic Infrastructure Reports"))
		})
	})
})
