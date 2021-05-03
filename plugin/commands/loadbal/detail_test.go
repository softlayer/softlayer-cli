package loadbal_test

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
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.DetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewDetailCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        metadata.LoadbalDetailMetadata().Name,
			Description: metadata.LoadbalDetailMetadata().Description,
			Usage:       metadata.LoadbalDetailMetadata().Usage,
			Flags:       metadata.LoadbalDetailMetadata().Flags,
			Action:      cmd.Run,
		}
	})
	Context("detail without loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "'--id' is required")).To(BeTrue())
		})
	})
	Context("detail with wrong loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "invalid value")).To(BeTrue())
		})
	})
	Context("detail with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(datatypes.Network_LBaaS_LoadBalancer{},
				errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "1234")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to get load balancer with ID 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("detail with loadbal ID", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerReturns(
				datatypes.Network_LBaaS_LoadBalancer{
					Id: sl.Int(13162),
				}, nil)
		})
		It("return loadbalancer", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "13162")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"13162"}))
		})
	})
})
