package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var fakeReturn = []datatypes.Network_Storage{
	datatypes.Network_Storage{
		Id:                        sl.Int(0001),
		Username:                  sl.String("sp-0001"),
		SnapshotCreationTimestamp: sl.String("2016-12-26T00:12:00"),
		SnapshotSizeBytes:         sl.String("500"),
	},
	datatypes.Network_Storage{
		Id:                        sl.Int(0002),
		Username:                  sl.String("sp-0002"),
		SnapshotCreationTimestamp: sl.String("2016-12-25T00:12:00"),
		SnapshotSizeBytes:         sl.String("540"),
	},
	datatypes.Network_Storage{
		Id:                        sl.Int(0003),
		Username:                  sl.String("sp-0003"),
		SnapshotCreationTimestamp: sl.String("2016-12-28T00:12:00"),
		SnapshotSizeBytes:         sl.String("100"),
	},
}

var _ = Describe("Snapshot list", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.SnapshotListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewSnapshotListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
		FakeStorageManager.GetVolumeIdReturns(1234, nil)
	})

	Describe("Snapshot list tests", func() {
		Context("Usage Errors", func() {
			It("No volumeid", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Bad --sortby", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "bcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby bcd is not supported."))
			})
		})

		Context("Snapshot list with corrrect volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(nil, errors.New("Internal Server Error"))
			})
			It("SL API ERROR", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get snapshot list on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Snapshot list with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(fakeReturn, nil)
			})
			It("Happy Path", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				// I don't like ContainSubstrings, but its useful for checking for multiple strings in a single line
				// ContainSubstrings will not check multiple lines though, which is why we have 3 calls here.
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1", "sp-0001", "2016-12-26T00:12:00", "500"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2", "sp-0002", "2016-12-25T00:12:00", "540"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"3", "sp-0003", "2016-12-28T00:12:00", "100"}))

			})
			It("Sorted by size_bytes", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "size_bytes")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(rows[1]).To(ContainSubstring("100"))
				Expect(rows[2]).To(ContainSubstring("500"))
				Expect(rows[3]).To(ContainSubstring("540"))
			})
			It("Sorted by created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "created")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(rows[1]).To(ContainSubstring("2016-12-25T00:12:00"))
				Expect(rows[2]).To(ContainSubstring("2016-12-26T00:12:00"))
				Expect(rows[3]).To(ContainSubstring("2016-12-28T00:12:00"))
			})
			It("Sorted by name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(rows[1]).To(ContainSubstring("sp-0001"))
				Expect(rows[2]).To(ContainSubstring("sp-0002"))
				Expect(rows[3]).To(ContainSubstring("sp-0003"))
			})
		})
	})
})
