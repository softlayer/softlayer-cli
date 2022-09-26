package loadbal_test

import (
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoadBalancer Suite")
}

var availableCommands = []string{
	"cancel",
	"detail",
	"health-edit",
	"l7member-add",
	"l7member-delete",
	"l7policies",
	"l7policy-add",
	"l7policy-delete",
	"l7policy-edit",
	"l7pool-add",
	"l7pool-delete",
	"l7pool-detail",
	"l7pool-edit",
	"l7rule-add",
	"l7rule-delete",
	"l7rules",
	"list",
	"member-add",
	"member-delete",
	"ns-detail",
	"ns-list",
	"order",
	"order-options",
	"protocol-add",
	"protocol-delete",
	"protocol-edit",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test load balancer commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := loadbal.SetupCobraCommands(slMeta)

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
	Context("LoadBalancer Namespace", func() {
		It("LoadBalancer Name Space", func() {
			Expect(loadbal.LoadbalNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(loadbal.LoadbalNamespace().Name).To(ContainSubstring("loadbal"))
			Expect(loadbal.LoadbalNamespace().Description).To(ContainSubstring("Classic infrastructure Load Balancers"))
		})
	})
})
