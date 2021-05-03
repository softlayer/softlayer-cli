package securitygroup_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup rule remove", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *securitygroup.RuleRemoveCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = securitygroup.NewRuleRemoveCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.SecurityGroupRuleRemoveMetaData().Name,
			Description: metadata.SecurityGroupRuleRemoveMetaData().Description,
			Usage:       metadata.SecurityGroupRuleRemoveMetaData().Usage,
			Flags:       metadata.SecurityGroupRuleRemoveMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Securitygroup rule remove", func() {
		Context("rule remove without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("rule remove without rule id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments."))
			})
		})
		Context("rule remove with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("rule remove with wrong rule id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group rule ID'. It must be a positive integer."))
			})
		})
		Context("rule remove with correct params but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234", "5678")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This will remove rule 5678 in security group 1234 and cannot be undone. Continue?"))
			})
		})
		Context("rule remove with correct params but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.RemoveSecurityGroupRuleReturns(errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to remove rule 5678 in security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("rule remove with correct params", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "5678", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Rule 5678 in security group 1234 is removed."))
			})
		})
	})
})
