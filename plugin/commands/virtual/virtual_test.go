package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestVirtual(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtual Server Suite")
}

var availableCommands = []string{
	"authorize-storage",
	"bandwidth",
	"billing",
	"cancel",
	"capacity-create",
	"capacity-create-options",
	"capacity-detail",
	"capacity-list",
	"capture",
	"create",
	"credentials",
	"detail",
	"dns-sync",
	"edit",
	"host-create",
	"host-list",
	"list",
	"migrate",
	"monitoring-list",
	"options",
	"pause",
	"placementgroup-create",
	"placementgroup-create-options",
	"placementgroup-detail",
	"placementgroup-list",
	"power-off",
	"power-on",
	"ready",
	"reboot",
	"reload",
	"rescue",
	"resume",
	"storage",
	"upgrade",
	"usage",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test Virtual Commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := virtual.SetupCobraCommands(slMeta)

		var arrayCommands = []string{}
		for _, command := range commands.Commands() {
			commandName := command.Name()
			arrayCommands = append(arrayCommands, commandName)
			It("available commands "+commands.Name(), func() {
				available := false
				if utils.StringInSlice(commandName, availableCommands) != -1 {
					available = true
				}
				Expect(available).To(BeTrue(), commandName+" not found in array available Commands")
			})
		}
		for _, command := range availableCommands {
			commandName := command
			It("ibmcloud sl "+commands.Name(), func() {
				available := false
				if utils.StringInSlice(commandName, arrayCommands) != -1 {
					available = true
				}
				Expect(available).To(BeTrue(), commandName+" not found in ibmcloud sl "+commands.Name())
			})
		}
	})
	
	Context("Virtual Namespace", func() {
		It("Virtual Namespace Exists", func() {
			Expect(virtual.VSNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(virtual.VSNamespace().Name).To(ContainSubstring("vs"))
			Expect(virtual.VSNamespace().Description).To(ContainSubstring("Classic infrastructure Virtual Servers"))
		})
	})
})
