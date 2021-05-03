package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/block"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block Snapshot Schedule List", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.SnapshotScheduleListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewSnapshotScheduleListCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.BlockSnapshotScheduleListMetaData().Name,
			Description: metadata.BlockSnapshotScheduleListMetaData().Description,
			Usage:       metadata.BlockSnapshotScheduleListMetaData().Usage,
			Action:      cmd.Run,
		}
	})

	Describe("block Snapshot Schedule List", func() {
		Context("No Arguments Error", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Proper Usage", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"replication"}))
			})
		})
		Context("Proper Usage, but API error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotSchedulesReturns(datatypes.Network_Storage{}, errors.New("Fake Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"replication"}))
				Expect(strings.Contains(err.Error(), "Fake Internal Server Error")).To(BeTrue())
			})
		})
	})
})
