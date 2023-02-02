package block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Block Suite")
}

var availableCommands = []string{
	"access-authorize",
	"access-list",
	"access-password",
	"access-revoke",
	"disaster-recovery-failover",
	"duplicate-convert-status",
	"object-list",
	"object-storage-detail",
	"replica-failback",
	"replica-failover",
	"replica-locations",
	"replica-order",
	"replica-partners",
	"snapshot-cancel",
	"snapshot-create",
	"snapshot-create",
	"snapshot-delete",
	"snapshot-disable",
	"snapshot-enable",
	"snapshot-get-notification-status",
	"snapshot-list",
	"snapshot-order",
	"snapshot-restore",
	"snapshot-schedule-list",
	"snapshot-set-notification",
	"subnets-assign",
	"subnets-list",
	"subnets-remove",
	"volume-cancel",
	"volume-convert",
	"volume-count",
	"volume-detail",
	"volume-duplicate",
	"volume-limits",
	"volume-list",
	"volume-modify",
	"volume-options",
	"volume-order",
	"volume-refresh",
	"volume-set-lun-id",
	"volume-set-note",
}

var _ = Describe("Test block Commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := block.SetupCobraCommands(slMeta)

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

	Context("Account Namespace", func() {
		It("Account Name Space", func() {
			Expect(block.BlockNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(block.BlockNamespace().Name).To(ContainSubstring("block"))
			Expect(block.BlockNamespace().Description).To(ContainSubstring("Classic infrastructure Block Storage"))
		})
	})
})
