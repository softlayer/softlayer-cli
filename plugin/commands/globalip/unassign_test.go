package globalip_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("GlobalIP unassign", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *globalip.UnassignCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = globalip.NewUnassignCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.GlobalIpUnassignMetaData().Name,
			Description: metadata.GlobalIpUnassignMetaData().Description,
			Usage:       metadata.GlobalIpUnassignMetaData().Usage,
			Flags:       metadata.GlobalIpUnassignMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("GlobalIP unassign", func() {
		Context("GlobalIP unassign without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("GlobalIP unassign with wrong IP ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Globalip ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("GlobalIP unassign with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.UnassignGlobalIPReturns(datatypes.Provisioning_Version1_Transaction{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to unassign global IP 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("GlobalIP assign with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.UnassignGlobalIPReturns(datatypes.Provisioning_Version1_Transaction{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The transaction to unroute a global IP address is created, routes will be updated in one or two minutes."}))
			})
		})
	})
})
