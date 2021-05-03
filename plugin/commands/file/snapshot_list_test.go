package file_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/file"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Snapshot list", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.SnapshotListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewSnapshotListCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.FileSnapshotListMetaData().Name,
			Description: metadata.FileSnapshotListMetaData().Description,
			Usage:       metadata.FileSnapshotListMetaData().Usage,
			Flags:       metadata.FileSnapshotListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Snapshot list", func() {
		Context("Snapshot list without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Snapshot list with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Snapshot list with wrong --sortby", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(nil, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--sortby", "bcd")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby bcd is not supported.")).To(BeTrue())
			})
		})

		Context("Snapshot list with corrrect volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get snapshot list on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Snapshot list with correct volume id", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(
					[]datatypes.Network_Storage{
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
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0001"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sp-0001"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-26T00:12:00"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"500"}))

				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0002"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sp-0002"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-25T00:12:00"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"540"}))

				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0003"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"sp-0003"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-28T00:12:00"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"100"}))
			})
		})

		Context("Snapshot list with correct volume id and --sortby=size_bytes", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(
					[]datatypes.Network_Storage{
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
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--sortby", "size_bytes")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(rows[1], "100")).To(BeTrue())
				Expect(strings.Contains(rows[2], "500")).To(BeTrue())
				Expect(strings.Contains(rows[3], "540")).To(BeTrue())
			})
		})

		Context("Snapshot list with correct volume id and --sortby=created", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(
					[]datatypes.Network_Storage{
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
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--sortby", "created")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(rows[1], "2016-12-25T00:12:00")).To(BeTrue())
				Expect(strings.Contains(rows[2], "2016-12-26T00:12:00")).To(BeTrue())
				Expect(strings.Contains(rows[3], "2016-12-28T00:12:00")).To(BeTrue())
			})
		})

		Context("Snapshot list with correct volume id and --sortby=created", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeSnapshotListReturns(
					[]datatypes.Network_Storage{
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
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				rows := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(rows[1], "sp-0001")).To(BeTrue())
				Expect(strings.Contains(rows[2], "sp-0002")).To(BeTrue())
				Expect(strings.Contains(rows[3], "sp-0003")).To(BeTrue())
			})
		})
	})
})
