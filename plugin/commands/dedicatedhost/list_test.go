package dedicatedhost_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host detail", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *dedicatedhost.ListCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		FakeDedicatedhostManager *testhelpers.FakeDedicatedHostManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dedicatedhost.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedHostManager)
		cliCommand.DedicatedHostManager = FakeDedicatedhostManager
	})

	Describe("Dedicatedhost list", func() {
		Context("list with server fails", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListDedicatedHostReturns([]datatypes.Virtual_DedicatedHost{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list dedicated hosts on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("list with nothing found", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListDedicatedHostReturns([]datatypes.Virtual_DedicatedHost{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No dedicated hosts are found."))
			})
		})

		Context("Set invalid --sortby option", func() {
			It("Set invalid --sortby option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=User")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid --sortBy option."))
			})
		})

		Context("list with hosts found", func() {
			BeforeEach(func() {
				fakerDedicatedHosts := []datatypes.Virtual_DedicatedHost{
					datatypes.Virtual_DedicatedHost{
						Id:         sl.Int(111111),
						Name:       sl.String("dedicatedhost01"),
						GuestCount: sl.Uint(0),
						BackendRouter: &datatypes.Hardware_Router_Backend{
							Hardware_Router: datatypes.Hardware_Router{
								Hardware_Switch: datatypes.Hardware_Switch{
									Hardware: datatypes.Hardware{
										Hostname: sl.String("bcr01a.dal13"),
									},
								},
							},
						},
						AllocationStatus: &datatypes.Container_Virtual_DedicatedHost_AllocationStatus{
							CpuAllocated:    sl.Int(0),
							CpuCount:        sl.Int(56),
							DiskAllocated:   sl.Int(0),
							DiskCapacity:    sl.Int(1200),
							MemoryAllocated: sl.Int(0),
							MemoryCapacity:  sl.Int(242),
						},
						Datacenter: &datatypes.Location{
							Name: sl.String("dal13"),
						},
					},
					datatypes.Virtual_DedicatedHost{
						Id:         sl.Int(222222),
						Name:       sl.String("dedicatedhost02"),
						GuestCount: sl.Uint(0),
						BackendRouter: &datatypes.Hardware_Router_Backend{
							Hardware_Router: datatypes.Hardware_Router{
								Hardware_Switch: datatypes.Hardware_Switch{
									Hardware: datatypes.Hardware{
										Hostname: sl.String("bcr01a.dal13"),
									},
								},
							},
						},
						AllocationStatus: &datatypes.Container_Virtual_DedicatedHost_AllocationStatus{
							CpuAllocated:    sl.Int(0),
							CpuCount:        sl.Int(56),
							DiskAllocated:   sl.Int(0),
							DiskCapacity:    sl.Int(1200),
							MemoryAllocated: sl.Int(0),
							MemoryCapacity:  sl.Int(242),
						},
						Datacenter: &datatypes.Location{
							Name: sl.String("dal13"),
						},
					},
				}
				FakeDedicatedhostManager.ListDedicatedHostReturns(fakerDedicatedHosts, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Name")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Datacenter")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Router")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Cpu")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Memory")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Disk")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=Guests")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dedicatedhost02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bcr01a.dal13"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/56"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/242 "))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0/1200"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0"))
			})
		})
	})
})
