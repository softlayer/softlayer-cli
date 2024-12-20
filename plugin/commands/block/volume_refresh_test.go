package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Block Volume Refresh", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeRefreshCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeRefreshCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Block Volume Refresh", func() {
		Context("No Arguments Error", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Bad VolumeId", func() {
			It("error resolving volume ID", func() {
				FakeStorageManager.GetVolumeIdReturns(0, errors.New("BAD Volume ID"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD Volume ID"))
			})
		})
		Context("Bad SnapshotId", func() {
			It("error resolving snapshot ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Snapshot ID'. It must be a positive integer."))
			})
		})
		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeRefreshReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				volId, dupId, force := FakeStorageManager.VolumeRefreshArgsForCall(0)
				Expect(force).To(Equal(false))
				Expect(dupId).To(Equal(5678))
				Expect(volId).To(Equal(1234))
			})
			It("Force Flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678", "--force-refresh")

				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				volId, dupId, force := FakeStorageManager.VolumeRefreshArgsForCall(0)
				Expect(force).To(Equal(true))
				Expect(dupId).To(Equal(5678))
				Expect(volId).To(Equal(1234))
			})
		})

		Context("Proper Usage, but API error", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeRefreshReturns(errors.New("Fake Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("Fake Internal Server Error"))
			})
		})
	})
})
