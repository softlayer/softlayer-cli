package file_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("file duplicate-convert-status", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.DuplicateConvertStatusCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewDuplicateConvertStatusCommand(fakeUI, fakeStorageManager)
		cliCommand = cli.Command{
			Name:        file.FileDuplicateConvertStatusMetaData().Name,
			Description: file.FileDuplicateConvertStatusMetaData().Description,
			Usage:       file.FileDuplicateConvertStatusMetaData().Usage,
			Flags:       file.FileDuplicateConvertStatusMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("file duplicate-convert-status", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeStorageManager.GetDuplicateConversionStatusReturns(datatypes.Container_Network_Storage_DuplicateConversionStatusInformation{}, errors.New("Failed to get duplicate conversion status"))
			})
			It("Failed get duplicate conversion status", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get duplicate conversion status"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerDuplicateConversionStatus := datatypes.Container_Network_Storage_DuplicateConversionStatusInformation{
					ActiveConversionStartTime:       sl.String("2022-06-13 14:59:17"),
					DeDuplicateConversionPercentage: sl.Int(68),
					VolumeUsername:                  sl.String("SL02SEVC123456_74"),
				}
				fakeStorageManager.GetDuplicateConversionStatusReturns(fakerDuplicateConversionStatus, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2022-06-13 14:59:17"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("68"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SL02SEVC123456_74"))
			})
		})
	})
})
