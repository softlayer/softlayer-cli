package securitygroup_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup rule list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cliCommand         *securitygroup.RuleListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewRuleListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Securitygroup rule list", func() {
		Context("rulelist without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("rule list with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("rule list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})
		Context("rule list with correct group id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSecurityGroupRulesReturns(nil, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get rules of security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("rule list zero result", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSecurityGroupRulesReturns(nil, nil)
			})
			It("return not found", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No rules are found for security group 1234."))
			})
		})
		Context("list non-zero result", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSecurityGroupRulesReturns([]datatypes.Network_SecurityGroup_Rule{
					datatypes.Network_SecurityGroup_Rule{
						Id:            sl.Int(48815),
						RemoteIp:      sl.String("169.0.0.1"),
						RemoteGroupId: sl.Int(45678),
						Direction:     sl.String("egress"),
						Ethertype:     sl.String("IPv6"),
						PortRangeMin:  sl.Int(80),
						PortRangeMax:  sl.Int(1000),
						Protocol:      sl.String("HTTP"),
					},
					datatypes.Network_SecurityGroup_Rule{
						Id:            sl.Int(48816),
						RemoteIp:      sl.String("168.0.0.1"),
						RemoteGroupId: sl.Int(45478),
						Direction:     sl.String("ingress"),
						Ethertype:     sl.String("IPv4"),
						PortRangeMin:  sl.Int(22),
						PortRangeMax:  sl.Int(400),
						Protocol:      sl.String("TCP"),
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "48815")).To(BeTrue())
				Expect(strings.Contains(results[2], "48816")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "48815")).To(BeTrue())
				Expect(strings.Contains(results[2], "48816")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "remoteIp")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "168.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[2], "169.0.0.1")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "remoteGroupId")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "45478")).To(BeTrue())
				Expect(strings.Contains(results[2], "45678")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "direction")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "egress")).To(BeTrue())
				Expect(strings.Contains(results[2], "ingress")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "portRangeMin")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "22")).To(BeTrue())
				Expect(strings.Contains(results[2], "80")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "portRangeMax")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "400")).To(BeTrue())
				Expect(strings.Contains(results[2], "1000")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--sortby", "protocol")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "HTTP")).To(BeTrue())
				Expect(strings.Contains(results[2], "TCP")).To(BeTrue())
			})
		})
	})
})
