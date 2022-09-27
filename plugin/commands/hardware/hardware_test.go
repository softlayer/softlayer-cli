package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hardware server Suite")
}

var availableCommands = []string{
	"add-notification",
	"authorize-storage",
	"bandwidth",
	"billing",
	"cancel",
	"cancel-reasons",
	"create",
	"create-options",
	"credentials",
	"detail",
	"edit",
	"guests",
	"list",
	"monitoring-list",
	"power-cycle",
	"power-off",
	"power-on",
	"reboot",
	"notifications",
	"reflash-firmware",
	"reload",
	"rescue",
	"sensor",
	"storage",
	"toggle-ipmi",
	"update-firmware",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test hardware.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := hardware.SetupCobraCommands(slMeta)

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

	Context("Hardware Namespace", func() {
		It("Hardware Name Space", func() {
			Expect(hardware.HardwareNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(hardware.HardwareNamespace().Name).To(ContainSubstring("hardware"))
			Expect(hardware.HardwareNamespace().Description).To(ContainSubstring("Classic infrastructure hardware servers"))
		})
	})
})
