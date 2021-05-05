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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Replica partners", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *file.ReplicaPartnersCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = file.NewReplicaPartnersCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.FileReplicaPartnersMetaData().Name,
			Description: metadata.FileReplicaPartnersMetaData().Description,
			Usage:       metadata.FileReplicaPartnersMetaData().Usage,
			Flags:       metadata.FileReplicaPartnersMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Replicant partners", func() {
		Context("replicant partners without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Replicant partners with server error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationPartnersReturns(nil, errors.New("Internal server error"))
			})
			It("error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Failed to get replication partners for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
			})
		})
		Context("Replicant partners", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationPartnersReturns(
					[]datatypes.Network_Storage{
						datatypes.Network_Storage{
							Id:         sl.Int(26876939),
							Username:   sl.String("IBM02SL278444_566_REP_1"),
							AccountId:  sl.Int(278444),
							CapacityGb: sl.Int(100),
						},
					}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"26876939"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"IBM02SL278444_566_REP_1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"278444"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"278444"}))
			})
		})
		Context("Replicant partners", func() {
			BeforeEach(func() {
				FakeStorageManager.GetReplicationPartnersReturns(
					[]datatypes.Network_Storage{}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"There are no replication partners for volume 1234."}))
			})
		})
	})
})
