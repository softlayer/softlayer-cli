package securitygroup_test

import (
	"errors"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("end to end test", func() {
	var (
		fakeSession        *session.Session
		fakeNetworkManager managers.NetworkManager
		fakeUI             *terminal.FakeUI
		cliCommand         *securitygroup.InterfaceListCommand
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeNetworkManager = managers.NewNetworkManager(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewInterfaceListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	It("return no error", func() {
		err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
		Expect(err).NotTo(HaveOccurred())
		results := strings.Split(fakeUI.Outputs(), "\n")
		Expect(strings.Contains(results[1], "smsgwmongos-app-02-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[1], "private")).To(BeTrue())
		Expect(strings.Contains(results[1], "10.209.128.248")).To(BeTrue())
		Expect(strings.Contains(results[2], "smsgwmongos-app-03-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[2], "private")).To(BeTrue())
		Expect(strings.Contains(results[2], "10.209.128.85")).To(BeTrue())
		Expect(strings.Contains(results[3], "smsgwmongos-app-01-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[3], "private")).To(BeTrue())
		Expect(strings.Contains(results[3], "10.209.128.228")).To(BeTrue())
	})

})
var _ = Describe("Securitygroup interface list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cliCommand         *securitygroup.InterfaceListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewInterfaceListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Securitygroup interface list", func() {
		Context("interface list without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("interface list with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("interface list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})
		Context("interface list but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("list zero result", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id: sl.Int(1234),
				}, nil)
			})
			It("return not found", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No interfaces are binded to security group 1234."))
			})
		})
		Context("list non-zero result", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id: sl.Int(1234),
					NetworkComponentBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(87654321),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
									Hostname:                sl.String("bluemix-cll-analytics"),
									PrimaryIpAddress:        sl.String("169.48.97.229"),
									PrimaryBackendIpAddress: sl.String("10.186.44.219"),
								},
								Port: sl.Int(0),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(12345678),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671788),
									Hostname:                sl.String("bluemix-cli-analytics"),
									PrimaryIpAddress:        sl.String("169.48.97.229"),
									PrimaryBackendIpAddress: sl.String("10.186.44.219"),
								},
								Port: sl.Int(1),
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[2], "87654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[2], "87654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "virtualServerId")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[2], "36671788")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "bluemix-cli-analytics")).To(BeTrue())
				Expect(strings.Contains(results[2], "bluemix-cll-analytics")).To(BeTrue())
			})
		})

		Context("list non-zero result without PrimaryIpAddress or PrimaryBackendIpAddress", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id: sl.Int(1234),
					NetworkComponentBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(87654321),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
									Hostname:                sl.String("bluemix-cll-analytics"),
									PrimaryBackendIpAddress: sl.String("10.186.44.219"),
								},
								Port: sl.Int(0),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(12345678),
								Guest: &datatypes.Virtual_Guest{
									Id:               sl.Int(36671790),
									Hostname:         sl.String("bluemix-cli-analytics"),
									PrimaryIpAddress: sl.String("169.48.97.229"),
								},
								Port: sl.Int(1),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(13345678),
								Guest: &datatypes.Virtual_Guest{
									Id:               sl.Int(36671789),
									Hostname:         sl.String("bluemix-cli-analytics"),
									PrimaryIpAddress: sl.String("169.48.97.229"),
								},
								Port: sl.Int(0),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(88654321),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671788),
									Hostname:                sl.String("bluemix-cll-analytics"),
									PrimaryBackendIpAddress: sl.String("10.186.44.219"),
								},
								Port: sl.Int(1),
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[1], "public")).To(BeTrue())
				Expect(strings.Contains(results[2], "13345678")).To(BeTrue())
				Expect(strings.Contains(results[3], "87654321")).To(BeTrue())
				Expect(strings.Contains(results[3], "private")).To(BeTrue())
				Expect(strings.Contains(results[4], "88654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[1], "public")).To(BeTrue())
				Expect(strings.Contains(results[2], "13345678")).To(BeTrue())
				Expect(strings.Contains(results[3], "87654321")).To(BeTrue())
				Expect(strings.Contains(results[3], "private")).To(BeTrue())
				Expect(strings.Contains(results[4], "88654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "virtualServerId")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[1], "private")).To(BeTrue())
				Expect(strings.Contains(results[2], "36671788")).To(BeTrue())
				Expect(strings.Contains(results[2], "public")).To(BeTrue())
				Expect(strings.Contains(results[3], "36671789")).To(BeTrue())
				Expect(strings.Contains(results[4], "36671790")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "bluemix-cli-analytics")).To(BeTrue())
				Expect(strings.Contains(results[4], "bluemix-cll-analytics")).To(BeTrue())
			})
		})

		Context("list non-zero result with PrimaryIpAddress or PrimaryBackendIpAddress = nil", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id: sl.Int(1234),
					NetworkComponentBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(87654321),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
									Hostname:                sl.String("bluemix-cll-analytics"),
									PrimaryBackendIpAddress: nil,
									PrimaryIpAddress:        nil,
								},
								Port: sl.Int(0),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(12345678),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671790),
									Hostname:                sl.String("bluemix-cli-analytics"),
									PrimaryBackendIpAddress: nil,
									PrimaryIpAddress:        nil,
								},
								Port: sl.Int(1),
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[2], "87654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[2], "87654321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "virtualServerId")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[2], "36671790")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "bluemix-cli-analytics")).To(BeTrue())
				Expect(strings.Contains(results[2], "bluemix-cll-analytics")).To(BeTrue())
			})
		})
		Context("list non-zero result with port = nil", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id: sl.Int(1234),
					NetworkComponentBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(87654321),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
									Hostname:                sl.String("bluemix-cll-analytics"),
									PrimaryBackendIpAddress: nil,
									PrimaryIpAddress:        nil,
								},
								Port: nil,
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Id: sl.Int(12345678),
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671790),
									Hostname:                sl.String("bluemix-cli-analytics"),
									PrimaryBackendIpAddress: nil,
									PrimaryIpAddress:        nil,
								},
								Port: sl.Int(1),
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "12345678")).To(BeTrue())
				Expect(strings.Contains(results[1], "public")).To(BeTrue())
				Expect(strings.Contains(results[2], "87654321")).To(BeTrue())
				Expect(strings.Contains(results[2], "private")).To(BeTrue())
			})
		})
	})
})
