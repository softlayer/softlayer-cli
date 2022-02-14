package block_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Access Password", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.AccessPasswordCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewAccessPasswordCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockAccessPasswordMetaData().Name,
			Description: block.BlockAccessPasswordMetaData().Description,
			Usage:       block.BlockAccessPasswordMetaData().Usage,
			Flags:       block.BlockAccessPasswordMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Access password", func() {
		Context("Access password without hostId", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Access password without password", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "124")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Access password with wrong hostId", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "password")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'allowed access host ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Access password with server fails", func() {
			BeforeEach(func() {
				FakeStorageManager.SetCredentialPasswordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "password")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(strings.Contains(err.Error(), "Failed to set password for host 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Access password", func() {
			BeforeEach(func() {
				FakeStorageManager.SetCredentialPasswordReturns(nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "password")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password is updated for host 1234."))
			})
		})
	})
})
