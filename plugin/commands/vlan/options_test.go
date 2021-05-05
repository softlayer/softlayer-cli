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

var _ = Describe("VLAN options", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *vlan.OptionsCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = vlan.NewOptionsCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.VlanOptionsMetaData().Name,
			Description: metadata.VlanOptionsMetaData().Description,
			Usage:       metadata.VlanOptionsMetaData().Usage,
			Flags:       metadata.VlanOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VLAN options", func() {
		Context("VLAN options but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListDatacentersReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list datacenters.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VLAN options", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListDatacentersReturns(map[int]string{1: "dal07"}, nil)
				fakeNetworkManager.ListRoutersReturns([]string{"bcr01a.dal07", "fcr01a.dal07"}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"public,private"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"bcr01a.dal07,fcr01a.dal07"}))
			})
		})
	})
})
