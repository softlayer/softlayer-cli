package ipsec_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec order", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cmd              *ipsec.OrderCommand
		cliCommand       cli.Command
		context          plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = ipsec.NewOrderCommand(fakeUI, fakeIPSecManager, context)
		cliCommand = cli.Command{
			Name:        metadata.IpsecOrderMetaData().Name,
			Description: metadata.IpsecOrderMetaData().Description,
			Usage:       metadata.IpsecOrderMetaData().Usage,
			Flags:       metadata.IpsecOrderMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Context("order without -d", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-d|--datacenter' is required")).To(BeTrue())
		})
	})
	Context("order with server fails", func() {
		BeforeEach(func() {
			fakeIPSecManager.OrderTunnelContextReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "-d", "dal09")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to order IPSec.Please try again later.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("order", func() {
		BeforeEach(func() {
			fakeIPSecManager.OrderTunnelContextReturns(datatypes.Container_Product_Order_Receipt{
				OrderId: sl.Int(12345),
			}, nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "-d", "dal09")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345 was placed."}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("You may run '%s sl ipsec list --order 12345' to find this IPSec VPN after it is ready.", cmd.Context.CLIName())}))
		})
	})
})
