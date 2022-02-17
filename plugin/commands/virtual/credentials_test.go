package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS credentials", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CredentialsCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCredentialsCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSCredentialsMetaData().Name,
			Description: virtual.VSCredentialsMetaData().Description,
			Usage:       virtual.VSCredentialsMetaData().Usage,
			Flags:       virtual.VSCredentialsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS credentials", func() {
		Context("VS credentials without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS credentials with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS credentials with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS credentials with correct VS ID ", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(
					datatypes.Virtual_Guest{
						Id: sl.Int(1234),
						OperatingSystem: &datatypes.Software_Component_OperatingSystem{
							Software_Component: datatypes.Software_Component{
								Passwords: []datatypes.Software_Component_Password{
									datatypes.Software_Component_Password{
										Username: sl.String("root"),
										Password: sl.String("password4root"),
									},
									datatypes.Software_Component_Password{
										Username: sl.String("db2admin"),
										Password: sl.String("password4db2admin"),
									},
								},
							},
						},
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"root"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"password4root"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"db2admin"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"password4db2admin"}))
			})
		})
	})
})
