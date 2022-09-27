package securitygroup_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Security Group Suite")
}

var availableCommands = []string{
	"create",
	"delete",
	"detail",
	"edit",
	"interface-add",
	"interface-list",
	"interface-remove",
	"list",
	"rule-add",
	"rule-edit",
	"rule-list",
	"rule-remove",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test securitygroup.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := securitygroup.SetupCobraCommands(slMeta)

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

	Context("Securitygroup Namespace", func() {
		It("Securitygroup Name Space", func() {
			Expect(securitygroup.SecurityGroupNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(securitygroup.SecurityGroupNamespace().Name).To(ContainSubstring("securitygroup"))
			Expect(securitygroup.SecurityGroupNamespace().Description).To(ContainSubstring("Classic infrastructure network security groups"))
		})
	})
})
