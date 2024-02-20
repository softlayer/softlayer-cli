package securitygroup_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
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
		cliCommand         *securitygroup.DetailCommand
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeNetworkManager = managers.NewNetworkManager(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	It("return no error", func() {
		err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
		Expect(err).NotTo(HaveOccurred())
		results := strings.Split(fakeUI.Outputs(), "\n")
		Expect(strings.Contains(results[9], "smsgwmongos-app-02-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[9], "private")).To(BeTrue())
		Expect(strings.Contains(results[9], "10.209.128.248")).To(BeTrue())
		Expect(strings.Contains(results[10], "smsgwmongos-app-03-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[10], "private")).To(BeTrue())
		Expect(strings.Contains(results[10], "10.209.128.85")).To(BeTrue())
		Expect(strings.Contains(results[11], "smsgwmongos-app-01-dal13.smsgwmongo.prd")).To(BeTrue())
		Expect(strings.Contains(results[11], "private")).To(BeTrue())
		Expect(strings.Contains(results[11], "10.209.128.228")).To(BeTrue())
	})

})

var _ = Describe("Securitygroup detail", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cliCommand         *securitygroup.DetailCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Securitygroup detail", func() {
		Context("detail without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("detail with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("detail with correct group id but server API call fails", func() {
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
		Context("detail with correct group id ", func() {
			BeforeEach(func() {
				fakeNetworkManager.GetSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id:          sl.Int(45507),
					Name:        sl.String("allow_ssh"),
					Description: sl.String("Allow all ingress TCP traffic on port 22"),
					Rules: []datatypes.Network_SecurityGroup_Rule{
						datatypes.Network_SecurityGroup_Rule{
							Id:           sl.Int(48805),
							Direction:    sl.String("ingress"),
							Ethertype:    sl.String("IPv4"),
							PortRangeMin: sl.Int(22),
							PortRangeMax: sl.Int(22),
							Protocol:     sl.String("TCP"),
						},
						datatypes.Network_SecurityGroup_Rule{
							Id:           sl.Int(48806),
							Direction:    sl.String("ingress"),
							Ethertype:    sl.String("IPv6"),
							PortRangeMin: sl.Int(22),
							PortRangeMax: sl.Int(22),
							Protocol:     sl.String("TCP"),
						},
					},
					NetworkComponentBindings: []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
									Hostname:                sl.String("bluemix-cli-analytics"),
									PrimaryIpAddress:        sl.String("169.48.97.229"),
									PrimaryBackendIpAddress: sl.String("10.186.44.219"),
								},
								Port: sl.Int(0),
							},
						},
						datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding{
							NetworkComponent: &datatypes.Virtual_Guest_Network_Component{
								Guest: &datatypes.Virtual_Guest{
									Id:                      sl.Int(36671787),
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
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "45507")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("45507"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("allow_ssh"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("   -           -                 ingress     IPv4         22               22               TCP"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("   -           -                 ingress     IPv6         22               22               TCP"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bluemix-cli-analytics   private     10.186.44.219"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("bluemix-cli-analytics   public      169.48.97.229"))
			})
		})
		Context("list non-zero result with PrimaryIpAddress or PrimaryBackendIpAddress", func() {
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
				Expect(strings.Contains(results[7], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[7], "private")).To(BeTrue())
				Expect(strings.Contains(results[8], "36671790")).To(BeTrue())
				Expect(strings.Contains(results[8], "public")).To(BeTrue())
				Expect(strings.Contains(results[9], "36671789")).To(BeTrue())
				Expect(strings.Contains(results[9], "private")).To(BeTrue())
				Expect(strings.Contains(results[10], "36671788")).To(BeTrue())
				Expect(strings.Contains(results[10], "public")).To(BeTrue())
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
				Expect(strings.Contains(results[7], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[8], "36671790")).To(BeTrue())
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
				Expect(strings.Contains(results[8], "36671790")).To(BeTrue())
				Expect(strings.Contains(results[8], "public")).To(BeTrue())
				Expect(strings.Contains(results[7], "36671787")).To(BeTrue())
				Expect(strings.Contains(results[7], "private")).To(BeTrue())
			})
		})
	})
})
