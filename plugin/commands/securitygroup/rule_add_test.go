package securitygroup_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup ruleadd", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *securitygroup.RuleAddCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = securitygroup.NewRuleAddCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SecurityGroupRuleAddMetaData().Name,
			Description: metadata.SecurityGroupRuleAddMetaData().Description,
			Usage:       metadata.SecurityGroupRuleAddMetaData().Usage,
			Flags:       metadata.SecurityGroupRuleAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Securitygroup rule add", func() {
		Context("rule add without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("rule list with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("rule add without direction", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -d|--direction has to be either egress or ingress."))
			})
		})
		Context("rule add with wrong direction", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -d|--direction has to be either egress or ingress."))
			})
		})
		Context("rule add with wrong ethertype", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "egress", "-e", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -e|--ether-type has to be either IPv4 or IPv6."))
			})
		})
		Context("rule add without protocol", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "egress", "-M", "8888")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -p|--protocal must be set when -M or -m is specified."))
			})
		})
		Context("rule add with wrong protocol", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "egress", "-M", "8888", "-p", "http")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Options for -p|--protocol are: icmp,tcp,udp"))
			})
		})
		Context("rule add with correct params but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSecurityGroupRuleReturns(datatypes.Network_SecurityGroup_Rule{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "egress", "-e", "IPv4", "-M", "8888", "-m", "20", "-p", "tcp")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add rule to security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("rule add with correct params", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddSecurityGroupRuleReturns(datatypes.Network_SecurityGroup_Rule{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "egress", "-e", "IPv4", "-M", "8888", "-m", "20", "-p", "tcp")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Rule is added to security group 1234."))
			})
		})
	})
})
