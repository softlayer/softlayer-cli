package metadata_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"testing"
)

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metadata Suite")
}

var _ = Describe("Metadata list Metadata", func() {
	var (
		fakeUI              *terminal.FakeUI
		cmd                 *metadata.MetadataCommand
		cliCommand          cli.Command
		fakeMetadataManager *testhelpers.FakeMetadataManager
	)
	BeforeEach(func() {
		fakeMetadataManager = new(testhelpers.FakeMetadataManager)
		fakeUI = terminal.NewFakeUI()
		cmd = metadata.NewMetadataCommand(fakeUI, fakeMetadataManager)
		cliCommand = cli.Command{
			Name:        metadata.MetadataMetadata().Name,
			Description: metadata.MetadataMetadata().Description,
			Usage:       metadata.MetadataMetadata().Usage,
			Flags:       metadata.MetadataMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Metadata command", func() {
		Context("Metadata options, Invalid Usage", func() {
			It("Set command without option", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Set unavailable option", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This option is not available."))
			})
		})

		Context("Metadata options, correct use", func() {
			It("return metadata", func() {
				fakeMetadataManager.CallAPIServiceReturns([]byte{},nil)
				err := testhelpers.RunCommand(cliCommand, "backend_ip")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("not yet equals"))
			})
		})
	})
})
