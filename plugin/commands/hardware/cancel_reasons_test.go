package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware cancelreason", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.CancelReasonsCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewCancelReasonsCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        metadata.HardwareCancelReasonsMetaData().Name,
			Description: metadata.HardwareCancelReasonsMetaData().Description,
			Usage:       metadata.HardwareCancelReasonsMetaData().Usage,
			Flags:       metadata.HardwareCancelReasonsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware cancel reasons", func() {
		Context("hardware cancel reasons", func() {
			It("return nil", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
