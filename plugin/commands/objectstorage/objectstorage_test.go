package objectstorage_test

import (
	"reflect"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ObjectStorage Suite")
}

var availableCommands = []string{
	"object-storage-accounts",
	"object-storage-endpoints",
	"object-storage-credential-list",
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test objectstorage.GetCommandActionBindings()", func() {
	var (
		context plugin.PluginContext
	)
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	context = plugin.InitPluginContext("softlayer")
	commands := objectstorage.GetCommandActionBindings(context, fakeUI, fakeSession)

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
				Expect(found).To(BeTrue(), cmdName+" needs to be added to availableCommands[] in objectstorage.go")
			})
		}
	})

	Context("ObjectStorage Namespace", func() {
		It("ObjectStorage Name Space", func() {
			Expect(objectstorage.ObjectStorageNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(objectstorage.ObjectStorageNamespace().Name).To(ContainSubstring("object-storage"))
			Expect(objectstorage.ObjectStorageNamespace().Description).To(ContainSubstring("Classic infrastructure Object Storage commands"))
		})
	})

	Context("ObjectStorage MetaData", func() {
		It("ObjectStorage MetaData", func() {
			Expect(objectstorage.ObjectStorageMetaData().Category).To(ContainSubstring("sl"))
			Expect(objectstorage.ObjectStorageMetaData().Name).To(ContainSubstring("object-storage"))
			Expect(objectstorage.ObjectStorageMetaData().Usage).To(ContainSubstring("${COMMAND_NAME} sl object-storage"))
			Expect(objectstorage.ObjectStorageMetaData().Description).To(ContainSubstring("Classic infrastructure Object Storage commands"))
		})
	})
})
