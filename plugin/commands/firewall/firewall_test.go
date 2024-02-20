package firewall_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Firewall Suite")
}

var availableCommands = []string{
	"add",
	"cancel",
	"detail",
	"edit",
	"list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test firewall commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := firewall.SetupCobraCommands(slMeta)

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

	Context("Firewall Namespace", func() {
		It("Firewall Name Space", func() {
			Expect(firewall.FirewallNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(firewall.FirewallNamespace().Name).To(ContainSubstring("firewall"))
			Expect(firewall.FirewallNamespace().Description).To(ContainSubstring("Classic infrastructure Firewalls"))
		})
	})
})
