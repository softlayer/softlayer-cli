package securitygroup_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup edit", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *securitygroup.EditCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = securitygroup.NewEditCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        securitygroup.SecurityGroupEditMetaData().Name,
			Description: securitygroup.SecurityGroupEditMetaData().Description,
			Usage:       securitygroup.SecurityGroupEditMetaData().Usage,
			Flags:       securitygroup.SecurityGroupEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Securitygroup edit", func() {
		Context("edit without groupid", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("edit with wrong group id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Security group ID'. It must be a positive integer."))
			})
		})
		Context("edit with correct group id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditSecurityGroupReturns(errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "updated")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to edit security group 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("edit with correct group id ", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "updated")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Security group 1234 is updated."))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-d", "updated-desc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Security group 1234 is updated."))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "updated", "-d", "updated-desc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Security group 1234 is updated."))
			})
		})
	})
})
