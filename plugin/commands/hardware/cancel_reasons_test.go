package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
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
			Name:        hardware.HardwareCancelReasonsMetaData().Name,
			Description: hardware.HardwareCancelReasonsMetaData().Description,
			Usage:       hardware.HardwareCancelReasonsMetaData().Usage,
			Flags:       hardware.HardwareCancelReasonsMetaData().Flags,
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
