package globalip_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("GlobalIP assign", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *globalip.AssignCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = globalip.NewAssignCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        globalip.GlobalIpAssignMetaData().Name,
			Description: globalip.GlobalIpAssignMetaData().Description,
			Usage:       globalip.GlobalIpAssignMetaData().Usage,
			Flags:       globalip.GlobalIpAssignMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("GlobalIP assign", func() {
		Context("GlobalIP assign without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})
		Context("GlobalIP assign with not enough parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires two arguments.")).To(BeTrue())
			})
		})

		Context("GlobalIP assign with wrong IP ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc", "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Globalip ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("GlobalIP assign with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AssignGlobalIPReturns(false, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "1.2.3.4")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to assign global IP 1234 to target 1.2.3.4.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("GlobalIP assign with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.AssignGlobalIPReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "1.2.3.4")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The transaction to modify a global IP route is created, routes will be updated in one or two minutes."}))
			})
		})
	})
})
