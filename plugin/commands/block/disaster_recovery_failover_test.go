package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Disaster Recovery Failover", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.DisasterRecoveryFailoverCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewDisasterRecoveryFailoverCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockDisasterRecoveryFailoverMetaData().Name,
			Description: block.BlockDisasterRecoveryFailoverMetaData().Description,
			Usage:       block.BlockDisasterRecoveryFailoverMetaData().Usage,
			Flags:       block.BlockDisasterRecoveryFailoverMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Replicant failover", func() {
		Context("replicant failover without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("replicant failover without replicant id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Replicant fail over with wrong volume id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Replicant fail over with wrong replica id", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Replica ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Replicant fail over with correct volume id and replica id", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Failover of volume 1234 to replica 5678 is now in progress."}))
			})
		})

		Context("Replicant fail over with correct volume id but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.DisasterRecoveryFailoverReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{
					"OK",
				}))
				Expect(strings.Contains(err.Error(), "Failover operation could not be initiated for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})
