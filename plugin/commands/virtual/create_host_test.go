package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS create host", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *virtual.CreateHostCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeVSManager      *testhelpers.FakeVirtualServerManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCreateHostCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Create host", func() {
		Context("Create host without hostname", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-H|--hostname' is required"))
			})
		})
		Context("Create host without domain", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-D|--domain' is required"))
			})
		})
		Context("Create host without datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-d|--datacenter' is required"))
			})
		})
		Context("Create host without vlan", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-v|--vlan-private' is required"))
			})
		})
		Context("Create host with wrong billing", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "dbd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-b|--billing] has to be either hourly or monthly."))
			})
		})
		Context("Create host with get vlan fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get vlan 123."))
			})
		})
		Context("Create host without -f and not continue", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{}, nil)
			})
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
		Context("Create host with create host fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.wdc07"),
							},
						},
					},
				}, nil)
				fakeVSManager.CreateDedicatedHostReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create dedicated host."))
			})
		})
		Context("Create host with no vlan router id host fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       nil,
								Hostname: sl.String("bcr01a.wdc07"),
							},
						},
					},
				}, nil)
				fakeVSManager.CreateDedicatedHostReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get vlan primary router ID."))
			})
		})

		Context("create host with succeed", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetVlanReturns(datatypes.Network_Vlan{
					Id: sl.Int(123),
					PrimaryRouter: &datatypes.Hardware_Router{
						Hardware_Switch: datatypes.Hardware_Switch{
							Hardware: datatypes.Hardware{
								Id:       sl.Int(1115295),
								Hostname: sl.String("bcr01a.wdc07")},
						},
					},
				}, nil)
				fakeVSManager.CreateDedicatedHostReturns(datatypes.Container_Product_Order_Receipt{OrderId: sl.Int(345678)}, nil)
			})
			It("return order", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
			It("return order", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "monthly", "-v", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
			It("return order", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "monthly", "-v", "123", "-s", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
			It("return order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "test", "-D", "softlayer.com", "-d", "dal09", "-b", "hourly", "-v", "123", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order 345678 was placed."))
			})
		})
	})
})
