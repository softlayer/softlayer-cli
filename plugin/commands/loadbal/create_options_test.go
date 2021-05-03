package loadbal_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer create options", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeLBManager      *testhelpers.FakeLoadBalancerManager
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *loadbal.OptionsCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = loadbal.NewOptionsCommand(fakeUI, fakeLBManager, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.LoadbalOrderOptionsMetadata().Name,
			Description: metadata.LoadbalOrderOptionsMetadata().Description,
			Usage:       metadata.LoadbalOrderOptionsMetadata().Usage,
			Flags:       metadata.LoadbalOrderOptionsMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("create options returns error", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerOptionsReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get load balancer product packages.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})

})
