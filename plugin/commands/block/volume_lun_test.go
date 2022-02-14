package block_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume lun", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.VolumeLunCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewVolumeLunCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockVolumeLunMetaData().Name,
			Description: block.BlockVolumeLunMetaData().Description,
			Usage:       block.BlockVolumeLunMetaData().Usage,
			Flags:       block.BlockVolumeLunMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume lun", func() {
		Context("Volume list without volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Volume lun without lun id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("Volume lun with wrong volume id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "123")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Volume lun with wrong lun id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'LUN ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Volume lun but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{}, errors.New("Server Internal Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to set LUN ID for volume 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Server Internal Error")).To(BeTrue())
			})
		})
		Context("Volume lun but response incorrect", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{Value: sl.String("5679")}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Errors()).To(ContainSubstring("Failed to confirm the new LUN ID on volume 1234."))
			})
		})
		Context("Volume lun ", func() {
			BeforeEach(func() {
				FakeStorageManager.SetLunIdReturns(datatypes.Network_Storage_Property{Value: sl.String("5678")}, nil)
			})
			It("succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Block volume 1234 is reporting LUN ID 5678."))
			})
		})
	})
})
