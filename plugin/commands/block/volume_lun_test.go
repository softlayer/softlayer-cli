package block_test

import (
	"errors"

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

var _ = Describe("Volume lun", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeLunCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeLunCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume lun", func() {
		Context("Volume list without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Volume lun without lun id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("Volume lun with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})
		})
		Context("Volume lun with wrong lun id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'LUN ID'. It must be a positive integer."))
			})
		})
		Context("Volume lun but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{}, errors.New("Server Internal Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to set LUN ID for volume 1234."))
				Expect(err.Error()).To(ContainSubstring("Server Internal Error"))
			})
		})
		Context("Volume lun but response incorrect", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{Value: sl.String("5679")}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Errors()).To(ContainSubstring("Failed to confirm the new LUN ID on volume 1234."))
			})
		})
		Context("Volume lun ", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{Value: sl.String("5678")}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Block volume 1234 is reporting LUN ID 5678."))
			})
		})
	})
})
