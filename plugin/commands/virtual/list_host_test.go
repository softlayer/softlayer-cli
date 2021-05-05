package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS list", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.ListHostCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewListHostCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSListHostMataData().Name,
			Description: metadata.VSListHostMataData().Description,
			Usage:       metadata.VSListHostMataData().Usage,
			Flags:       metadata.VSListHostMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("list host", func() {
		Context("list with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.ListDedicatedHostReturns([]datatypes.Virtual_DedicatedHost{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list dedicated hosts on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("list with nothing found", func() {
			BeforeEach(func() {
				fakeVSManager.ListDedicatedHostReturns([]datatypes.Virtual_DedicatedHost{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No dedicated hosts are found."))
			})
		})
		Context("list with hosts found", func() {
			BeforeEach(func() {
				fakeVSManager.ListDedicatedHostReturns([]datatypes.Virtual_DedicatedHost{
					datatypes.Virtual_DedicatedHost{
						Id:             sl.Int(52001),
						Name:           sl.String("wilma-test"),
						CpuCount:       sl.Int(56),
						MemoryCapacity: sl.Int(242),
						DiskCapacity:   sl.Int(1200),
						GuestCount:     sl.Uint(1),
						BackendRouter: &datatypes.Hardware_Router_Backend{
							Hardware_Router: datatypes.Hardware_Router{
								Hardware_Switch: datatypes.Hardware_Switch{
									Hardware: datatypes.Hardware{
										Id:       sl.Int(1115295),
										Hostname: sl.String("bcr01a.wdc07"),
									},
								},
							},
						},
						AllocationStatus: &datatypes.Container_Virtual_DedicatedHost_AllocationStatus{
							CpuAllocated:    sl.Int(2),
							CpuAvailable:    sl.Int(54),
							CpuCount:        sl.Int(56),
							GuestCount:      sl.Int(1),
							DiskAllocated:   sl.Int(102),
							DiskAvailable:   sl.Int(1097),
							DiskCapacity:    sl.Int(1200),
							MemoryAllocated: sl.Int(1),
							MemoryAvailable: sl.Int(241),
							MemoryCapacity:  sl.Int(242),
						},
						Datacenter: &datatypes.Location{
							Id:       sl.Int(2017603),
							Name:     sl.String("wdc07"),
							LongName: sl.String("Washington 7"),
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("No dedicated hosts are found."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma-test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.wdc07"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("102/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/242"))
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "-n", "wilma-test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("No dedicated hosts are found."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma-test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.wdc07"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("102/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/242"))
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "-d", "wdc07")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("No dedicated hosts are found."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma-test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.wdc07"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("102/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/242"))
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--owner", "278444_wangjunl@cn.ibm.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("No dedicated hosts are found."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma-test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.wdc07"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("102/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/242"))
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--order", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("No dedicated hosts are found."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("wilma-test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.wdc07"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("102/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/242"))
			})
		})
	})
})
