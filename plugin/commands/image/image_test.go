package image_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test image commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		imageCommands := image.SetupCobraCommands(slMeta)
		Expect(imageCommands.Name()).To(Equal("image"))
	})

	Context("Image Namespace", func() {
		It("Image Name Space", func() {
			Expect(image.ImageNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(image.ImageNamespace().Name).To(ContainSubstring("image"))
			Expect(image.ImageNamespace().Description).To(ContainSubstring("Classic infrastructure Compute images"))
		})
	})
})
