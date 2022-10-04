package globalip_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GlobalIP Suite")
}

var availableCommands = []string{
	"assign",
	"cancel",
	"create",
	"list",
	"unassign",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test globalip commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		commands := globalip.SetupCobraCommands(slMeta)

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

	Context("Network Attached Storage Namespace", func() {
		It("Network Attached Storage Name Space", func() {
			Expect(globalip.GlobalIpNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(globalip.GlobalIpNamespace().Name).To(ContainSubstring("globalip"))
			Expect(globalip.GlobalIpNamespace().Description).To(ContainSubstring("Classic infrastructure Global IP addresses"))
		})
	})
})
