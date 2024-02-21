package tags_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Tags list", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeTagsManager *testhelpers.FakeTagsManager
		cliCommand      *tags.ListCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeTagsManager = new(testhelpers.FakeTagsManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = tags.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.TagsManager = fakeTagsManager
	})
	Describe("Tags list", func() {
		//sl tags list
		Context("Tags list, no arguments", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(3))
				Expect(results[1]).To(ContainSubstring("TEST TAG"))
			})
		})
		//sl tags list --output JSON
		Context("Tags list, Details", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
				fakeTagsManager.GetTagReferencesReturns(FakeTagReference, nil)
				fakeTagsManager.ReferenceLookupReturns("Hardware.Name")
			})
			It("Returns JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring("TEST TAG"))
			})
		})
		Context("Tags list, no arguments, ListTags error", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(nil, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		Context("Tags list, no arguments, ListEmptyTags error", func() {
			BeforeEach(func() {
				fakeTagsManager.ListEmptyTagsReturns(nil, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		//sl tags list -d
		Context("Tags list, Details", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
				fakeTagsManager.GetTagReferencesReturns(FakeTagReference, nil)
				fakeTagsManager.ReferenceLookupReturns("Hardware.Name")
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(4))
				Expect(results[2]).To(ContainSubstring("22222   HARDWARE   Hardware.Name"))
			})
		})
		//sl tags list -d --output JSON
		Context("Tags list, Details", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
				fakeTagsManager.GetTagReferencesReturns(FakeTagReference, nil)
				fakeTagsManager.ReferenceLookupReturns("Hardware.Name")
			})
			It("Returns JSON", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "--output", "JSON")
				Expect(err).NotTo(HaveOccurred())
				results := fakeUI.Outputs()
				Expect(results).To(ContainSubstring(`"ResourceName": "Hardware.Name"`))
			})
		})
		Context("Tags list, Details, ListTags error", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(nil, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		Context("Tags list, Details, GetTagReferences error", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
				fakeTagsManager.GetTagReferencesReturns([]datatypes.Tag_Reference{}, errors.New("SoftLayer_Exception_ApiError"))
			})
			It("Handle Exception", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(4))
				Expect(results[1]).To(ContainSubstring("Resource"))
				Expect(results[2]).To(ContainSubstring("SoftLayer_Exception_ApiError"))
			})
		})
		Context("Tags list, Details, Missing KeyName", func() {
			BeforeEach(func() {
				fakeTagsManager.ListTagsReturns(FakeTags, nil)
				some_number := 1111
				fake_tag := []datatypes.Tag_Reference{
					datatypes.Tag_Reference{
						Id:              &some_number,
						ResourceTableId: &some_number,
						TagType:         nil,
					},
				}
				fakeTagsManager.GetTagReferencesReturns(fake_tag, nil)
			})
			It("Handle missing KeyName", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(4))
				Expect(results[2]).To(ContainSubstring("1111   None"))
			})
		})
	})
})
