package block_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Access Password", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.AccessPasswordCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = block.NewAccessPasswordCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager

	})

	Describe("Access password", func() {
		Context("Access password without hostId", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Access password without password", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "124")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Access password with wrong hostId", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "password")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'allowed access host ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Access password with server fails", func() {
			BeforeEach(func() {
				FakeStorageManager.SetCredentialPasswordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "password")
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "password")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Password is updated for host 1234."))
			})
		})
	})
})
