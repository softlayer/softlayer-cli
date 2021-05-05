package vlan_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN Edit", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.EditCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = vlan.NewEditCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.VlanEditMetaData().Name,
			Description: metadata.VlanEditMetaData().Description,
			Usage:       metadata.VlanEditMetaData().Usage,
			Flags:       metadata.VlanEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN edit", func() {
		Context("VLAN edit without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VLAN edit with wrong vlan id", func() {
			It("error resolving vlan ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'VLAN ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VLAN edit without -n", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '-n|--name' is required")).To(BeTrue())
			})
		})

		Context("VLAN edit with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditVlanReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "myvlan")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to edit VLAN: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VLAN edit with correct vlan id", func() {
			BeforeEach(func() {
				fakeNetworkManager.EditVlanReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "myvlan")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"VLAN 1234 was updated."}))
			})
		})
	})
})
