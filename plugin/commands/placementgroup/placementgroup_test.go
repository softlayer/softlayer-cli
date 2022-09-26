package placementgroup_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PlacementGroup Suite")
}

var availableCommands = []string{
	"create",
	"create-options",
	"delete",
	"detail",
	"list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test placementgroup.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := placementgroup.SetupCobraCommands(slMeta)

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

	Context("PlacementGroup Namespace", func() {
		It("PlacementGroup Name Space", func() {
			Expect(placementgroup.PlacementGroupNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(placementgroup.PlacementGroupNamespace().Name).To(ContainSubstring("placement-group"))
			Expect(placementgroup.PlacementGroupNamespace().Description).To(ContainSubstring("Classic infrastructure Placement Group"))
		})
	})
})
