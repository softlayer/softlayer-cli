package user_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

// These are all the commands in user.go
var availableCommands = []string{
	"user-create",
	"user-list",
	"user-delete",
	"user-detail",
	"user-permissions",
	"user-detail-edit",
	"user-permission-edit",
	"user-notifications",
	"user-edit-notifications",
	"user-device-access",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test user.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := user.GetCommandActionBindings(context, fakeUI, fakeSession)

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

	Context("User Namespace", func() {
		It("User Name Space", func() {
			Expect(user.UserNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(user.UserNamespace().Name).To(ContainSubstring("user"))
			Expect(user.UserNamespace().Description).To(ContainSubstring("Classic infrastructure Manage Users"))
		})
	})

	Context("User MetaData", func() {
		It("User MetaData", func() {
			Expect(user.UserMetaData().Category).To(ContainSubstring("sl"))
			Expect(user.UserMetaData().Name).To(ContainSubstring("user"))
			Expect(user.UserMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl user"))
			Expect(user.UserMetaData().Description).To(ContainSubstring("Classic infrastructure Manage Users"))
		})
	})
})
