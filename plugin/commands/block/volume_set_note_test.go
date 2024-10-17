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

var _ = Describe("Block Volume Set Note", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeSetNoteCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeSetNoteCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Block Volume Set Note", func() {
		Context("No Argument Error", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(false, errors.New("This command requires one argument"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one argument"))
			})
		})

		Context("No flag error", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(false, errors.New("This command requires note flag."))
			})
			It("error resolving flag note", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`"note" not set`))
			})
		})

		Context("Bad VolumeId", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(false, errors.New("Invalid input for 'Volume ID'. It must be a positive integer."))
				FakeStorageManager.GetVolumeIdReturns(0, errors.New("BAD Volume ID"))
			})
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "--note=thisismynote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BAD Volume ID"))
			})
		})

		Context("Bad output format", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(false, errors.New("Invalid output format, only JSON is supported now."))
			})
			It("error resolving volume ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note=thisismynote", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid output format, only JSON is supported now."))
			})
		})

		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(false, errors.New("Error occurred while note was adding in block volume"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note=thisismynote")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Error occurred while note was adding in block volume"))
			})
		})

		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note=thisismynote")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The note was set successfully"))
			})
		})

		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeSetNoteReturns(true, nil)
			})
			It("return no error in json format", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--note=thisismynote", "--output=json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
			})
		})
	})
})
