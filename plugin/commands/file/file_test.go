package file_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Suite")
}

var availableCommands = []string{
	"access-authorize",
	"access-list",
	"access-revoke",
	"disaster-recovery-failover",
	"duplicate-convert-status",
	"replica-failback",
	"replica-failover",
	"replica-locations",
	"replica-order",
	"replica-partners",
	"snapshot-cancel",
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
	"volume-set-note",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test file.GetCommandActionBindings()", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		commands := file.SetupCobraCommands(slMeta)

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

	Context("File Namespace", func() {
		It("File Name Space", func() {
			Expect(file.FileNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(file.FileNamespace().Name).To(ContainSubstring("file"))
			Expect(file.FileNamespace().Description).To(ContainSubstring("Classic infrastructure File"))
		})
	})
})
