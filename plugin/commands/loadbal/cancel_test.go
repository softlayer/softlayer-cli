package loadbal_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer cancel", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.CancelCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewCancelCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalCancelMetadata().Name,
			Description: loadbal.LoadbalCancelMetadata().Description,
			Usage:       loadbal.LoadbalCancelMetadata().Usage,
			Flags:       loadbal.LoadbalCancelMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("cancel without loadbalID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '--id' is required")).To(BeTrue())
		})
	})
	Context("cancel without confirmation", func() {
		It("return aborted", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCommand(cliCommand, "--id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will cancel the load balancer: 1234 and cannot be undone. Continue?"}))
		})
	})
	Context("cancel with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCommand(cliCommand, "--id", "1234")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This will cancel the load balancer: 1234 and cannot be undone. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("cancel with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.CancelLoadBalancerReturns(true, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "1234", "-f")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to cancel load balancer 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("cancel with correct load balancer ID", func() {
		BeforeEach(func() {
			fakeLBManager.CancelLoadBalancerReturns(true, nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "1234", "-f")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Load balancer 1234 is cancelled."}))
		})
	})
	Context("cancel with server error, load balancer with UUID fail", func() {
		BeforeEach(func() {
			fakeLBManager.GetLoadBalancerUUIDReturns("", errors.New("Internal server error"))
		})
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "--id", "1234", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get load balancer: Internal server error."))
		})
	})
})
