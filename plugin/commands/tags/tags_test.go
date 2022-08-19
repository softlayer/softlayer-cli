package tags_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var FakeTags = []datatypes.Tag{
	datatypes.Tag{
		Id:             sl.Int(1234),
		Name:           sl.String("TEST TAG"),
		ReferenceCount: sl.Uint(1),
	},
}

var FakeTagReference = []datatypes.Tag_Reference{
	datatypes.Tag_Reference{
		Id:              sl.Int(1111),
		ResourceTableId: sl.Int(22222),
		TagType: &datatypes.Tag_Type{
			Description: sl.String("Test Tag"),
			KeyName:     sl.String("HARDWARE"),
		},
	},
}

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tags Suite")
}

// This test suite exists to make sure commands don't get accidently removed from the SetupCobraCommands
var _ = Describe("Test tags commands", func() {
	fakeUI := terminal.NewFakeUI()
	fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
	slMeta := metadata.NewSoftlayerCommand(fakeUI, fakeSession)

	Context("New commands testable", func() {
		tagsCommands := tags.SetupCobraCommands(slMeta)
		Expect(tagsCommands.Name()).To(Equal("tags"))
	})

	Context("Report Namespace", func() {
		It("Report Name Space", func() {
			Expect(tags.TagsNamespace().ParentName).To(ContainSubstring("sl"))
			Expect(tags.TagsNamespace().Name).To(ContainSubstring("tags"))
			Expect(tags.TagsNamespace().Description).To(ContainSubstring("Classic infrastructure Tag management"))
		})
	})

})
