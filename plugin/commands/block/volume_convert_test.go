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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block Volume Convert", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.VolumeConvertCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewVolumeConvertCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.BlockVolumeConvertMetaData().Name,
			Description: metadata.BlockVolumeConvertMetaData().Description,
			Usage:       metadata.BlockVolumeConvertMetaData().Usage,
			Action:      cmd.Run,
		}
	})

	Describe("block Volume Convert", func() {
		Context("No Arguments Error", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Bad VolumeId", func() {
			It("error resolving volume ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Volume ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("Proper Usage", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeConvertReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			})
		})

		Context("Proper Usage, but API error", func() {
			BeforeEach(func() {
				FakeStorageManager.VolumeConvertReturns(errors.New("Fake Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"OK"}))
				Expect(strings.Contains(err.Error(), "Fake Internal Server Error")).To(BeTrue())
			})
		})
	})
})
