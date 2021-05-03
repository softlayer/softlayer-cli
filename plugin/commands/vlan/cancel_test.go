package vlan_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN Cancel", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.CancelCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = vlan.NewCancelCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.VlanCancelMetaData().Name,
			Description: metadata.VlanCancelMetaData().Description,
			Usage:       metadata.VlanCancelMetaData().Usage,
			Flags:       metadata.VlanCancelMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN cancel", func() {
		Context("VLAN cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VLAN cancel with wrong vlan id", func() {
			It("error resolving vlan ID", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'VLAN ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VLAN cancel with correct vlan id", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"VLAN 1234 was cancelled."}))
			})
		})

		Context("VLAN cancel with correct vlan id but vlan is not found", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Unable to find VLAN with ID 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound")).To(BeTrue())
			})
		})

		Context("VLAN cancel with correct vlan id but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.CancelVLANReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to cancel VLAN 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})
