package hardware_test

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware create options", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.CreateOptionsCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware create options", func() {
		Context("hardware create options", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetCreateOptionsReturns(map[string]map[string]string{
					managers.KEY_LOCATIONS:  map[string]string{"dal10": "Dallas 10"},
					managers.KEY_SIZES:      map[string]string{"D2620_128GB_2X1T_SATA_RAID_1xM60_GPU": "Dual Xeon 2620v4, 128GB Ram, 2x800GB SSD disks, RAID1"},
					managers.KEY_OS:         map[string]string{"CENTOS_6_32": "CentOS 6.5-32"},
					managers.KEY_NAME_OS:         map[string]string{"CENTOS_6_32": "OS_CENTOS_6_X_32_BIT"},
					managers.KEY_PORT_SPEED: map[string]string{"100_MBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS": "10000"},
					managers.KEY_PORT_SPEED_DESCRIPTION: map[string]string{"100_MBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS": "10 Gbps Redundant Public & Private Network Uplinks"},
					managers.KEY_EXTRAS:     map[string]string{"8_PUBLIC_IP_ADDRESSES": "8 Public IP Addresses"},
				})
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("D2620_128GB_2X1T_SATA_RAID_1xM60_GPU"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("CENTOS_6_32"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OS_CENTOS_6_X_32_BIT"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("100_MBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8_PUBLIC_IP_ADDRESSES"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Routers"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("fcr01a.ams01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Amsterdam 1"))
			})
		})
	})
})
