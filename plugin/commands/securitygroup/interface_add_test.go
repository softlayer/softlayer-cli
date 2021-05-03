package securitygroup_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup interface add", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		fakeVSManager      *testhelpers.FakeVirtualServerManager
		cmd                *securitygroup.InterfaceAddCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = securitygroup.NewInterfaceAddCommand(fakeUI, fakeNetworkManager, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.SecurityGroupInterfaceAddMetaData().Name,
			Description: metadata.SecurityGroupInterfaceAddMetaData().Description,
			Usage:       metadata.SecurityGroupInterfaceAddMetaData().Usage,
			Flags:       metadata.SecurityGroupInterfaceAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Securitygroup interface add", func() {
		Context("interface add without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("interface add without componentid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Must set either -n|--network-component or both -s|--server and -i|--interface"))
			})
		})
		Context("interface add with componentid and serverid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "2345", "-s", "3456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Must set either -n|--network-component or both -s|--server and -i|--interface"))
			})
		})
		Context("interface add with componentid and interface", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "2345", "-i", "abdf")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Must set either -n|--network-component or both -s|--server and -i|--interface"))
			})
		})
		Context("interface add with serverid and wronginterface", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "2345", "-i", "abdf")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -i|--interface must be either public or private"))
			})
		})
		Context("interface add with componentID but API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AttachSecurityGroupComponentReturns(errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "4567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add network component 4567 to security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("interface add with componentID succeed", func() {
			BeforeEach(func() {
				fakeNetworkManager.AttachSecurityGroupComponentReturns(nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "4567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Network component 4567 is added to security group 1234."))
			})
		})
		Context("interface add with serverID but getInstance API call fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "public")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "private")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("interface add with serverID but get component ID fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id: sl.Int(4321),
					NetworkComponents: []datatypes.Virtual_Guest_Network_Component{
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4567),
							Port: sl.Int(1),
						},
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4569),
							Port: sl.Int(1),
						},
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4568),
							Port: sl.Int(0),
						},
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4566),
							Port: sl.Int(0),
						},
					},
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "public")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Instance 4321 has 2 public interface."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "private")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Instance 4321 has 2 private interface."))
			})
		})
		Context("interface add with serverID but attach API call fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id: sl.Int(4321),
					NetworkComponents: []datatypes.Virtual_Guest_Network_Component{
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4567),
							Port: sl.Int(1),
						},
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4568),
							Port: sl.Int(0),
						},
					},
				}, nil)
				fakeNetworkManager.AttachSecurityGroupComponentReturns(errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "public")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add network component 4567 to security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "private")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add network component 4568 to security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("interface add with server succeed", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id: sl.Int(4321),
					NetworkComponents: []datatypes.Virtual_Guest_Network_Component{
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4567),
							Port: sl.Int(1),
						},
						datatypes.Virtual_Guest_Network_Component{
							Id:   sl.Int(4568),
							Port: sl.Int(0),
						},
					},
				}, nil)
				fakeNetworkManager.AttachSecurityGroupComponentReturns(nil)
			})
			It("return succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "public")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Network component 4567 is added to security group 1234."))
			})
			It("return succeed", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-s", "4321", "-i", "private")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Network component 4568 is added to security group 1234."))
			})
		})
	})
})
