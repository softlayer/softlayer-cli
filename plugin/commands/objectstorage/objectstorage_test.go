package objectstorage_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ObjectStorage Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the actionBindings
var _ = Describe("Test objectstorage commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
	Context("New commands testable", func() {
		objectstorageCommands := objectstorage.SetupCobraCommands(slMeta)
		Expect(objectstorageCommands.Name()).To(Equal("object-storage"))
	})
	Context("ObjectStorage Namespace", func() {
		It("ObjectStorage Name Space", func() {
			Expect(objectstorage.ObjectStorageNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(objectstorage.ObjectStorageNamespace().Name).To(ContainSubstring("object-storage"))
			Expect(objectstorage.ObjectStorageNamespace().Description).To(ContainSubstring("Classic infrastructure Object Storage commands"))
		})
	})
})
