package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

func TestMetadata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SL Metadata Suite")
}

var _ = Describe("SL Metadata Tests", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	Describe("Happy Path Tests", func() {
		Context("Tests Get Version", func() {
			version := metadata.GetVersion()
			It("Make sure we got valid numbers", func() {
				Expect(version.Major).Should(BeNumerically(">", 0))
				Expect(version.Major).Should(BeNumerically("<", 99))
				Expect(version.Major).Should(BeNumerically("<", 99))
			})
		})
		Context("Tests SoftLayerNamespace", func() {
			result := metadata.SoftlayerNamespace()
			It("Make sure we got valid numbers", func() {
				Expect(result.Name).Should(Equal(metadata.NS_SL_NAME))
				Expect(result.Description).Should(ContainSubstring("Manage Classic"))
			})
		})
		Context("Tests CobraOutputFlag", func() {
			It("Getting and Setting output flag", func() {
				result := metadata.CobraOutputFlag{"test1"}
				Expect(result.String()).Should(Equal("test1"))
				err := result.Set("should Error")
				Expect(err).To(HaveOccurred())
				err = result.Set("JsoN")
				Expect(err).NotTo(HaveOccurred())
				Expect(result.String()).To(Equal("JSON"))
				Expect(result.Type()).To(Equal("string"))
			})
		})
		Context("Can create a New SoftLayer Command", func() {
			It("NewSoftlayerCommand", func() {
				result := metadata.NewSoftlayerCommand(fakeUI, fakeSession)
				Expect(result.Session).ShouldNot(BeNil())
				Expect(result.GetOutputFlag()).To(Equal(""))
			})
			It("NewSoftlayerStorageCommand", func() {
				result := metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
				Expect(result.Session).ShouldNot(BeNil())
				Expect(result.GetOutputFlag()).To(Equal(""))
			})
		})
	})

})
