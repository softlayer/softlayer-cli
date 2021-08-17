package hardware_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Authorize Block, File Storage to a Hardware Server", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.AuthorizeStorageCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewAuthorizeStorageCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        metadata.HardwareAuthorizeStorageMataData().Name,
			Description: metadata.HardwareAuthorizeStorageMataData().Description,
			Usage:       metadata.HardwareAuthorizeStorageMataData().Usage,
			Flags:       metadata.HardwareAuthorizeStorageMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Authorize Block, File Storage to a Hardware Server", func() {
		Context("Authorize Storage without HW ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Authorize Storage with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Hardware server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Authorize storage to a Hardware Server", func() {
			BeforeEach(func() {
				fakeHardwareManager.AuthorizeStorageReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--username-storage", "SL02SL11111111-11")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			})
		})

		Context("Error Authorize Storage to a VS", func() {
			BeforeEach(func() {
				fakeHardwareManager.AuthorizeStorageReturns(false, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--username-storage", "SL02SL111")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to authorize storage to the hardware server instance: {{.Storage}}.\n{{.Error}}")).To(BeTrue())
			})
		})
	})
})
