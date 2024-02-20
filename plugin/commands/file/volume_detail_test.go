package file_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *file.VolumeDetailCommand
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
		slCommand          *metadata.SoftlayerStorageCommand

	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Volume detail", func() {
		Context("Volume detail without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Volume detail with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})

		Context("Volume detail with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Network_Storage", "getObject", 500, "Internal Server Error")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get details of volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Volume detail with correct volume id", func() {
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                         17336531"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("User name                  IBM01SEL278444-16"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Endurance Tier             WRITEHEAVY_TIER"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Notes                      -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Encrypted                  False"))
			})
		})
	})
})
